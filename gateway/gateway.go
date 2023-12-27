package gateway

import (
	"context"
	"fmt"
	"strings"

	"github.com/Scrin/ruuvi-go-gateway/config"
	"github.com/Scrin/ruuvi-go-gateway/sender"
	"github.com/rigado/ble"
	"github.com/rigado/ble/examples/lib/dev"
	"github.com/rigado/ble/linux/hci/cmd"
	log "github.com/sirupsen/logrus"
)

func Run(config config.Config) {
	gwMac := config.GwMac
	if gwMac == "" {
		gwMac = "00:00:00:00:00:00"
	}
	useMQTT := false
	useHTTP := false
	if config.MQTT != nil && (config.MQTT.Enabled == nil || *config.MQTT.Enabled) {
		sender.SetupMQTT(*config.MQTT)
		useMQTT = true
	}
	if config.HTTP != nil && (config.HTTP.Enabled == nil || *config.HTTP.Enabled) {
		sender.SetupHTTP(*config.HTTP, gwMac)
		useHTTP = true
	}
	if !useMQTT && !useHTTP {
		log.Fatal("Neither MQTT nor HTTP is configured, check the config")
	}

	device, err := dev.NewDevice("default",
		ble.OptDeviceID(config.HciIndex),
		ble.OptScanParams(cmd.LESetScanParameters{
			LEScanType: 0, // passive scan
			LEScanInterval: 0x10,
			LEScanWindow: 0x10,
		}))
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"hci_index": config.HciIndex,
		}).Fatal("Can't setup bluetooth device")
	}
	ble.SetDefaultDevice(device)
	advHandler := func(adv ble.Advertisement) {
		data := adv.ManufacturerData()
		if len(data) > 2 {
			isRuuvi := data[0] == 0x99 && data[1] == 0x04 // ruuvi company identifier
			log.WithFields(log.Fields{
				"mac":      strings.ToUpper(adv.Addr().String()),
				"rssi":     adv.RSSI(),
				"is_ruuvi": isRuuvi,
				"data":     fmt.Sprintf("%X", data),
			}).Trace("Received data from BLE adapter")
			if config.AllAdvertisements || isRuuvi {
				if useMQTT {
					sender.SendMQTT(*config.MQTT, adv, gwMac)
				}
				if useHTTP {
					sender.SendHTTP(*config.HTTP, adv)
				}
			}
		}
	}

	err = ble.Scan(context.Background(), true, advHandler, nil)
	if err != nil {
		log.WithError(err).Error("Failed to scan")
	}
}
