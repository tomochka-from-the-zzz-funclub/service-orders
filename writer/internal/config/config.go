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
	// err := godotenv.Load("configs/docker.env")

	// if err != nil {
	// 	log.Fatalf("Ошибка загрузки .env файла: %v", err)
	// }

	log.Println("=", os.Getenv("ENV"), "=", os.Getenv("ENV") != "docker")
	fmt.Printf("=%s=", os.Getenv("ENV"))
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
