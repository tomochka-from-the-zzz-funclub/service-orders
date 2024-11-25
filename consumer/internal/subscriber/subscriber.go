package subscriber

// import (
// 	"fmt"

// 	"consumer/internal/config"
// 	myLog "consumer/internal/logger"

// 	//"github.com/IBM/sarama"
// 	"github.com/IBM/sarama"
// 	//"github.com/IBM/sarama"
// )

// // KafkaClient для работы с Kafka
// type KafkaClient struct {
// 	Consumer  *sarama.Consumer
// 	Partition *sarama.PartitionConsumer
// }

// //var S service.Srv

// // NewKafkaClient - конструктор для создания клиента Kafka и запуска бесконечного чтения
// func NewKafkaClient(cfg config.Config) *KafkaClient {

// 	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}, nil)
// 	if err != nil {
// 		myLog.Log.Errorf("Failed to create consumer: %v", err)
// 	}
// 	if err != nil {
// 		myLog.Log.Errorf("Error create Consumer in consumer: %v", err.Error())
// 		return nil
// 	}

// 	partConsumer, err := consumer.ConsumePartition(cfg.KafkaTopic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		myLog.Log.Errorf("Error create Partition in consumer %v", err.Error())
// 	}
// 	kc := &KafkaClient{
// 		Consumer:  &consumer,
// 		Partition: &partConsumer,
// 	}
// 	return kc
// }

// func (kc *KafkaClient) ReadOrders() models.Order {
// 	var receivedMessage models.Order
// 	for {
// 				//select {
// 				//case
// 				//myLog.Log.Debugf(s.consumer.Consumer, (s.consumer.Partition == nil))
// 				fmt.Println((s.consumer.Consumer == nil), (s.consumer.Partition == nil))
// 				receivedMessage := s.consumer.ReadOrders()
// 				// msg, ok := <-(*s.consumer.Partition).Messages():
// 				// if !ok {
// 				// 	myLog.Log.Debugf("Channel closed, exiting")
// 				// 	return
// 				// }

// 				// var receivedMessage models.Order
// 				// err := json.Unmarshal(msg.Value, &receivedMessage)

// 				// if err != nil {
// 				// 	myLog.Log.Debugf("Error unmarshaling JSON: %v\n", err)
// 				// 	continue
// 				// }

// 				myLog.Log.Debugf("Received message: %+v\n", receivedMessage)

// 				//}
// 			}
// 	select {
// 	case msg, ok := <-(*kc.Partition).Messages():
// 		myLog.Log.Debugf("ERROR: %v", ok)
// 		if !ok {
// 			myLog.Log.Debugf("Channel closed, exiting")
// 			return models.Order{}
// 		}

// 		err := json.Unmarshal(msg.Value, &receivedMessage)

// 		if err != nil {
// 			myLog.Log.Debugf("Error unmarshaling JSON: %v\n", err)
// 			return models.Order{}
// 		}

// 		myLog.Log.Debugf("Received message: %+v\n", receivedMessage)
// 	}
// 	return receivedMessage
// }
