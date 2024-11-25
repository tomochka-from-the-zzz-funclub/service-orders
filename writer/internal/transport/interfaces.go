package transport

type PubInterface interface {
	// Close - метод для корректного закрытия продюсера и консьюмера
	Close()

	// SendMessage - метод для отправки сообщения в Kafka
	SendMessage(topic string, key string, value interface{}) error
}
