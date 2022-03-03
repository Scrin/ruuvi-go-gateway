package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type MQTT struct {
	Enabled       *bool  `yaml:"enabled,omitempty"`
	BrokerUrl     string `yaml:"broker_url"`
	BrokerAddress string `yaml:"broker_address"`
	BrokerPort    int    `yaml:"broker_port"`
	ClientID      string `yaml:"client_id"`
	Username      string `yaml:"username"`
	Password      string `yaml:"password"`
	TopicPrefix   string `yaml:"topic_prefix"`
}

type HTTP struct {
	Enabled  *bool         `yaml:"enabled,omitempty"`
	URL      string        `yaml:"url"`
	Interval time.Duration `yaml:"interval"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
}

type Logging struct {
	Type       string `yaml:"type"`
	Level      string `yaml:"level"`
	Timestamps *bool  `yaml:"timestamps,omitempty"`
	WithCaller bool   `yaml:"with_caller,omitempty"`
}

type Config struct {
	GwMac             string  `yaml:"gw_mac"`
	AllAdvertisements bool    `yaml:"all_advertisements"`
	HciIndex          int     `yaml:"hci_index"`
	MQTT              *MQTT   `yaml:"mqtt,omitempty"`
	HTTP              *HTTP   `yaml:"http,omitempty"`
	Logging           Logging `yaml:"logging"`
	Debug             bool    `yaml:"debug"`
}

func ReadConfig(configFile string, strict bool) (Config, error) {
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return Config{}, fmt.Errorf("no config found! Tried to open \"%s\"", configFile)
	}

	f, err := os.Open(configFile)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var conf Config
	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(strict)
	err = decoder.Decode(&conf)

	if err != nil {
		return Config{}, err
	}
	return conf, nil
}
