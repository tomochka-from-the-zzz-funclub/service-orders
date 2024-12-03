package transport

type PubInterface interface {
	Close()
	SendMessage(topic string, key string, value interface{}) error
}
