// package service

// import (
// 	store "consumer/internal/cache"
// 	"consumer/internal/config"
// 	"consumer/internal/database"
// 	myLog "consumer/internal/logger"
// 	"consumer/internal/models"
// 	"encoding/json"
// 	"fmt"

// 	"github.com/IBM/sarama"
// )

// type Srv struct {
// 	cache store.Store
// 	db    *database.Postgres
// 	//consumer *subscriber.KafkaClient
// }

// func NewSrv(cfg config.Config) *Srv {
// 	// d, err := database.NewPostgres(cfg)
// 	// if err != nil {
// 	// 	myLog.Log.Fatalf("Don't create service consumer", err.Error())
// 	// }
// 	//c := subscriber.NewKafkaClient(cfg)
// 	s := &Srv{
// 		cache: store.NewStore(),
// 		db:    database.NewPostgres(cfg),
// 		//consumer: c,
// 	}
// 	// go s.Read(cfg)
// 	return s
// }

// func (s *Srv) GetOrderSrv(orderUUID string) (models.Order, error) {
// 	myLog.Log.Debugf("GetOrderSrv")
// 	order, err := s.cache.Get(orderUUID)
// 	if err != nil {
// 		orderdb, err := s.db.GetOrder(orderUUID)
// 		if err != nil {
// 			return models.Order{}, err
// 		}
// 		s.cache.Add(order)
// 		return orderdb, nil
// 	}
// 	return order, nil
// }

// func (s *Srv) SetOrder(order models.Order) (string, error) {
// 	myLog.Log.Debugf("SetOrder")
// 	id, err := s.db.AddOrder(order)
// 	if err != nil {
// 		return "", err
// 	}
// 	s.cache.Add(order)
// 	return id, nil
// }

// func (s *Srv) Read(cfg config.Config) {
// 	myLog.Log.Debugf("Start func Read in consumer")
// 	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}, nil)
// 	if err != nil {
// 		myLog.Log.Errorf("Failed to create consumer: %v", err)
// 	}
// 	defer consumer.Close()
// 	myLog.Log.Errorf("Suc consumer")
// 	partConsumer, err := consumer.ConsumePartition(cfg.KafkaTopic, 0, sarama.OffsetNewest)
// 	if err != nil {
// 		myLog.Log.Errorf("Failed to consume partition: %v", err)
// 	}
// 	defer partConsumer.Close()
// 	myLog.Log.Errorf("Suc patr")
// 	for {
// 		select {
// 		case msg, ok := <-partConsumer.Messages():
// 			if !ok {
// 				myLog.Log.Debugf("Channel closed, exiting")
// 				return
// 			}
// 			myLog.Log.Debugf("in select")
// 			var receivedMessage models.Order
// 			err := json.Unmarshal(msg.Value, &receivedMessage)

// 			if err != nil {
// 				myLog.Log.Debugf("Error unmarshaling JSON: %v\n", err)
// 				continue
// 			}

// 			myLog.Log.Debugf("Received message: %+v\n", receivedMessage)

// 			s.SetOrder(receivedMessage)
// 			s.cache.Add(receivedMessage)
// 			myLog.Log.Debugf("Success read order in kafka with id: ", receivedMessage.OrderUID)
// 		}
// 		myLog.Log.Debugf("END")
// 	}
// }

package service

import (
	store "consumer/internal/cache"
	"consumer/internal/database"
	myLog "consumer/internal/logger"
	"consumer/internal/models"
	"encoding/json"
	"fmt"

	"consumer/internal/config"

	"github.com/IBM/sarama"
)

type Srv struct {
	cache store.Store
	db    *database.Postgres
}

func NewSrv(cfg config.Config) *Srv {
	return &Srv{
		cache: store.NewStore(),
		db:    database.NewPostgres(cfg),
	}
}

func (s *Srv) GetOrderSrv(orderUUID string) (models.Order, error) {
	myLog.Log.Debugf("GetOrderSrv")
	order, err := s.cache.Get(orderUUID)
	if err != nil {
		orderdb, err := s.db.GetOrder(orderUUID)
		if err != nil {
			return models.Order{}, err
		}
		s.cache.Add(order)
		return orderdb, nil
	}
	return order, nil
}

func (s *Srv) SetOrder(order models.Order) (string, error) {
	myLog.Log.Debugf("SetOrder")
	id, err := s.db.AddOrder(order)
	if err != nil {
		return "", err
	}
	s.cache.Add(order)
	return id, nil
}

func (s *Srv) Read(cfg config.Config) {
	myLog.Log.Debugf("Start Read")
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", cfg.KafkaHost, cfg.KafkaPort)}, nil)
	if err != nil {
		myLog.Log.Errorf("Failed to create consumer: %v", err)
	}
	defer consumer.Close()

	partConsumer, err := consumer.ConsumePartition(cfg.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		myLog.Log.Errorf("Failed to consume partition: %v", err)
	}
	defer partConsumer.Close()

	for {
		select {
		case msg, ok := <-partConsumer.Messages():
			if !ok {
				myLog.Log.Debugf("Channel closed, exiting")
				return
			}

			var receivedMessage models.Order
			err := json.Unmarshal(msg.Value, &receivedMessage)

			if err != nil {
				myLog.Log.Debugf("Error unmarshaling JSON: %v\n", err)
				continue
			}

			myLog.Log.Debugf("Received message: %+v\n", receivedMessage)

			// responseText := receivedMessage.Name + " " + receivedMessage.Value + " ( " + receivedMessage.ID + " ) "

			// resp := &sarama.ProducerMessage{
			// 	Topic: "pong",
			// 	Key:   sarama.StringEncoder(receivedMessage.ID),
			// 	Value: sarama.StringEncoder(responseText),
			// }

			// _, _, err = producer.SendMessage(resp)
			// if err != nil {
			// 	myLog.Log.Printf("Failed to send message to Kafka: %v", err)
			// }
			s.SetOrder(receivedMessage)
			s.cache.Add(receivedMessage)
			myLog.Log.Debugf("success")
		}
	}
}
