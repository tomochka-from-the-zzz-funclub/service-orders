package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Kafka_BROKER []string
	KafkaTopic   string
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "docker" {
		if err := godotenv.Load("configs/local.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	fmt.Println([]string{os.Getenv("KafkaPort")}, os.Getenv("KafkaTopic"))
	return Config{
		Kafka_BROKER: []string{os.Getenv("KafkaPort")},
		KafkaTopic:   os.Getenv("KafkaTopic"),
	}
}
