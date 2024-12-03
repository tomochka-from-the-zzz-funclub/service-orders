package publisher

import (
	"encoding/json"

	myLog "writer/internal/logger"
	"writer/internal/models"

	"github.com/IBM/sarama"
)

type KafkaClient struct {
	Producer sarama.SyncProducer
}

func NewKafkaClient() *KafkaClient {
	producer, err := sarama.NewSyncProducer([]string{"kafka:9092"}, nil)

	if err != nil {
		myLog.Log.Debugf("Error create producer", err.Error())
		return nil
	}
	k := KafkaClient{
		Producer: producer,
	}
	myLog.Log.Debugf("Succes create producer kafka")
	return &k
}

func (kc *KafkaClient) SendOrderToKafka(topic string, message models.Order) error {

	msg_json, err := json.Marshal(message)
	if err != nil {
		myLog.Log.Errorf("Error marshal msg-order")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: "order_set",
		Key:   sarama.StringEncoder(message.OrderUID),
		Value: sarama.ByteEncoder(msg_json),
	}

	_, _, err = kc.Producer.SendMessage(msg)
	if err != nil {
		myLog.Log.Errorf("Failed to send message to Kafka: %v", err)
		return err
	}

	return err
}
