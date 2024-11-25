package publisher

import (
	"encoding/json"
	//"sync"

	//"github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/config"
	//"github.com/IBM/sarama"
	//"github.com/IBM/sarama"
	"github.com/IBM/sarama"
	myLog "github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/logger"
	"github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/models"
	//"github.com/IBM/sarama"
	//"github.com/confluentinc/confluent-kafka-go/kafka"
	//"github.com/tomochka-from-the-zzz-funclub/go-L0-Kafka/internal/config"
)

// KafkaClient - структура для хранения продюсера и консьюмера Kafka
type KafkaClient struct {
	Producer sarama.SyncProducer
	//Consumer sarama.Consumer
}

// var responseChannels map[string]chan *sarama.ConsumerMessage
// var mu sync.Mutex

// NewKafkaClient - конструктор для создания новой структуры KafkaClient
func NewKafkaClient() *KafkaClient {
	// producer, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": "kafka:9092",
	// })
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

// // CreateTopic создает новый топик в Kafka
// func (kc *KafkaClient) CreateTopic(topic string, partitions int32, replicationFactor int16) error {
// 	admin, err := sarama.NewClusterAdmin([]string{"kafka:9092"}, nil)
// 	if err != nil {
// 		myLog.Log.Errorf("Create admin", err)
// 		return err
// 	}
// 	defer admin.Close()

// 	// Создание конфигурации топика
// 	topicDetail := &sarama.TopicDetail{
// 		NumPartitions:     partitions,
// 		ReplicationFactor: replicationFactor,
// 	}

// 	// Создание топика
// 	err = admin.CreateTopic(topic, topicDetail, false)
// 	if err != nil {
// 		myLog.Log.Errorf("Create topic", err)
// 		return err
// 	}

// 	log.Printf("Топик %s успешно создан", topic)
// 	return nil
// }

// // // Close - метод для корректного закрытия продюсера и консьюмера
// // func (kc *KafkaClient) Close() {
// // 	if err := (kc.Producer).Close(); err != nil {
// // 		log.Printf("Failed to close producer: %v", err)
// // 	}
// // }

// func (kc *KafkaClient) SendOrderToKafka(topic string, order models.Order) error {
// 	//fmt.Println("SendOrderToKafka")
// 	//fmt.Println(kc)
// 	// if kc.Producer == nil {
// 	// 	fmt.Println("nil")
// 	// }
// 	myLog.Log.Debugf("Start SendOrderToKafka")

// 	//kc.CreateTopic(topic, 1, 1)
// 	//fmt.Println("SendOrderToKafka")
// 	// Сериализуем структуру Order в JSON
// 	orderData, err := json.Marshal(order)
// 	if err != nil {
// 		myLog.Log.Debugf("Err parse")
// 		return err
// 	}
// 	//fmt.Println(orderData)
// 	// Создаем сообщение с сериализованными данными
// 	// msg := &sarama.ProducerMessage{
// 	// 	Topic: topic,
// 	// 	Key:   sarama.StringEncoder(topic),
// 	// 	Value: sarama.ByteEncoder(orderData),
// 	// }
// 	fmt.Println("сообщение создано")
// 	myLog.Log.Debugf("Отправлено сообщение в топик %s", topic)
// 	// Отправляем сообщение в Kafka
// 	err = kc.Producer.Produce(&kafka.Message{
// 		TopicPartition: kafka.TopicPartition{
// 			Topic:     &topic,
// 			Partition: kafka.PartitionAny,
// 		},
// 		Value: []byte(orderData),
// 	}, nil)
// 	fmt.Println("сообщение отправлено")
// 	if err != nil {
// 		myLog.Log.Errorf("Send////")
// 		return err
// 	}
// 	return nil
// }

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
