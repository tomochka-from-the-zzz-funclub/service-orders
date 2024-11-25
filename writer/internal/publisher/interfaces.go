package publisher

import "writer/internal/models"

type InterfaceKafkaClient interface {
	SendOrderToKafka(topic string, message models.Order) error
}
