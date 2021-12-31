package logging

import (
	"github.com/Scrin/ruuvi-go-gateway/config"
	log "github.com/sirupsen/logrus"
)

func Setup(conf config.Logging) {
	timestamps := true
	if conf.Timestamps != nil {
		timestamps = *conf.Timestamps
	}

	log.SetReportCaller(conf.WithCaller)

	if conf.WithCaller {
		if timestamps {
			log.SetFormatter(new(PlainFormatterWithTsWithCaller))
		} else {
			log.SetFormatter(new(PlainFormatterWithoutTsWithCaller))
		}
	} else {
		if timestamps {
			log.SetFormatter(new(PlainFormatterWithTsWithoutCaller))
		} else {
			log.SetFormatter(new(PlainFormatterWithoutTsWithoutCaller))
		}
	}

	switch conf.Type {
	case "structured":
		log.SetFormatter(&log.TextFormatter{
			DisableTimestamp: !timestamps,
		})
	case "json":
		log.SetFormatter(&log.JSONFormatter{
			DisableTimestamp: !timestamps,
		})
	case "simple":
	case "":
	default:
		log.Fatal("Invalid logging type: ", conf.Type)
	}

	switch conf.Level {
	case "trace":
		log.SetLevel(log.TraceLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "":
		log.SetLevel(log.InfoLevel)
	default:
		log.Fatal("Invalid logging level: ", conf.Level)
	}
}
