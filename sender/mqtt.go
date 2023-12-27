package sender

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Scrin/ruuvi-go-gateway/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/rigado/ble"
	log "github.com/sirupsen/logrus"
)

type mqttMessage struct {
	GwMac  string        `json:"gw_mac"`
	Rssi   int           `json:"rssi"`
	Aoa    []interface{} `json:"aoa"`
	Gwts   string        `json:"gwts"`
	Ts     string        `json:"ts"`
	Data   string        `json:"data"`
	Coords string        `json:"coords"`
}

var client mqtt.Client

func SetupMQTT(conf config.MQTT) {
	address := conf.BrokerAddress
	if address == "" {
		address = "localhost"
	}
	port := conf.BrokerPort
	if port == 0 {
		port = 1883
	}
	server := conf.BrokerUrl
	if server == "" {
		server = fmt.Sprintf("tcp://%s:%d", address, port)
	}
	log.WithFields(log.Fields{
		"target":       server,
		"topic_prefix": conf.TopicPrefix,
	}).Info("Starting MQTT")

	clientID := conf.ClientID
	if clientID == "" {
		clientID = "ruuvi-go-gateway"
	}
	opts := mqtt.NewClientOptions()
	opts.SetCleanSession(false)
	opts.AddBroker(server)
	opts.SetClientID(clientID)
	opts.SetUsername(conf.Username)
	opts.SetPassword(conf.Password)
	opts.SetKeepAlive(10 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)
	if conf.LWTTopic == nil || *conf.LWTTopic != "" {
		topic := conf.TopicPrefix + "/gw_status"
		payload := conf.LWTOfflinePayload
		if conf.LWTTopic != nil {
			topic = *conf.LWTTopic
		}
		if payload == "" {
			payload = "{\"state\":\"offline\"}"
		}
		opts.SetWill(topic, payload, 0, true)
	}
	client = mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if conf.LWTTopic == nil || *conf.LWTTopic != "" {
		topic := conf.TopicPrefix + "/gw_status"
		payload := conf.LWTOnlinePayload
		if conf.LWTTopic != nil {
			topic = *conf.LWTTopic
		}
		if payload == "" {
			payload = "{\"state\":\"online\"}"
		}
		client.Publish(topic, 0, true, payload)
	}
}

func SendMQTT(conf config.MQTT, adv ble.Advertisement, gwMac string) {
	mac := strings.ToUpper(adv.Addr().String())
	data := adv.ManufacturerData()
	flags := []byte{0x00} // the actual advertisement flags don't seem to be available, so just use zero
	message := mqttMessage{
		GwMac:  gwMac,
		Rssi:   adv.RSSI(),
		Aoa:    []interface{}{},
		Gwts:   fmt.Sprint(time.Now().Unix()),
		Ts:     fmt.Sprint(time.Now().Unix()),
		Data:   fmt.Sprintf("0201%X%XFF%X", flags, len(data)+1, data),
		Coords: "",
	}
	data, err := json.Marshal(message)
	if err != nil {
		log.WithError(err).Error("Failed to serialize data")
	} else {
		client.Publish(conf.TopicPrefix+"/"+mac, 0, false, string(data))
	}
}
