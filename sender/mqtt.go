package sender

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Scrin/ruuvi-go-gateway/config"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ble/ble"
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

var mqttClient mqtt.Client

func SetupMQTT(conf config.MQTT) {
	address := conf.BrokerAddress
	if address == "" {
		address = "localhost"
	}
	port := conf.BrokerPort
	if port == 0 {
		port = 1883
	}
	server := fmt.Sprintf("tcp://%s:%d", address, port)
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
	mqttClient = mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
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
		mqttClient.Publish(conf.TopicPrefix+"/"+mac, 0, false, string(data))
	}
}
