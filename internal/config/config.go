package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

type AppConfig struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string

	Kafka_Brokers  string
	Kafka_Topic    []string
	Kafka_Group_ID string

	Elastic_URL string
}

func Load() AppConfig {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg AppConfig
	err = viper.Unmarshal(&cfg)
	if err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	cfg.Kafka_Topic = strings.Split(viper.GetString("KAFKA_TOPICS"), ",")

	return cfg
}
