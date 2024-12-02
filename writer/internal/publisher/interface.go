package publisher

import "writer/internal/models"

//go:generate  -source=interface.go -destination=mocks/mock_kafka_client.go -package=mocks

type InterfaceKafkaClient interface {
	SendOrderToKafka(topic string, message models.Order) error
}
