package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	SourceType string
	SinkType   string

	MySQL struct {
		DSN string
	}

	Kafka struct {
		Brokers []string
		Topic   string
		GroupID string
	}

	Elastic struct {
		URL   string
		Index string
	}
}

func Load() AppConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg AppConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	return cfg
}
