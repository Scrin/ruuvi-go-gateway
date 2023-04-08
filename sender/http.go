package sender

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Scrin/ruuvi-go-gateway/config"
	"github.com/go-ble/ble"
	log "github.com/sirupsen/logrus"
)

type httpMessage struct {
	Data httpMessageData `json:"data"`
}

type httpMessageData struct {
	Coordinates string                    `json:"coordinates"`
	Timestamp   string                    `json:"timestamp"`
	Nonce       string                    `json:"nonce"`
	GwMac       string                    `json:"gw_mac"`
	Tags        map[string]httpMessageTag `json:"tags"`
}

type httpMessageTag struct {
	Rssi      int    `json:"rssi"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
}

var lock sync.Mutex
var tags map[string]httpMessageTag = make(map[string]httpMessageTag)
var httpClient http.Client

func handleBodyClose(req *http.Request) {
	if err := req.Body.Close(); err != nil {
		log.WithError(err).Error(fmt.Sprintf("Failed to close request body: %s", err))
	}
}

func SetupHTTP(conf config.HTTP, gwMac string) {
	log.WithFields(log.Fields{
		"url":      conf.URL,
		"interval": conf.Interval,
	}).Info("Starting HTTP")

	httpClient = http.Client{
		Timeout: conf.Interval,
	}

	go func() {
		ticker := time.NewTicker(conf.Interval)
		for {
			<-ticker.C
			msg := httpMessage{Data: httpMessageData{
				Coordinates: "",
				Timestamp:   fmt.Sprint(time.Now().Unix()),
				Nonce:       "",
				GwMac:       gwMac,
				Tags:        tags,
			}}
			lock.Lock()
			data, err := json.Marshal(msg)
			lock.Unlock()
			if err != nil {
				log.WithError(err).Error("Failed to serialize data")
			}
			req, err := http.NewRequest(http.MethodPost, conf.URL, strings.NewReader(string(data)))
			if err != nil {
				handleBodyClose(req)
				log.WithError(err).Error("Failed create a POST request")
			}
			if conf.Username != "" {
				req.SetBasicAuth(conf.Username, conf.Password)
			}

			_, err = httpClient.Do(req)
			if err != nil {
				handleBodyClose(req)
				log.WithError(err).Error("Failed POST data")
			}

			handleBodyClose(req)
		}
	}()
}

func SendHTTP(conf config.HTTP, adv ble.Advertisement) {
	mac := strings.ToUpper(adv.Addr().String())
	data := adv.ManufacturerData()
	flags := []byte{0x00} // the actual advertisement flags don't seem to be available, so just use zero
	tag := httpMessageTag{
		Rssi:      adv.RSSI(),
		Timestamp: fmt.Sprint(time.Now().Unix()),
		Data:      fmt.Sprintf("0201%X%XFF%X", flags, len(data)+1, data),
	}
	lock.Lock()
	tags[mac] = tag
	lock.Unlock()
}
