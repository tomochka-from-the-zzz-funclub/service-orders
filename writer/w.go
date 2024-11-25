package main

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	// Создаем продюсера
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
	})
	if err != nil {
		fmt.Println("Ошибка создания продюсера:", err.Error())
		return
	}
	defer producer.Close()

	deliveryChan := make(chan kafka.Event)

	topic := "records"

	// Отправляем сообщение
	err = producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte("tttttteesr"),
	}, deliveryChan)

	fmt.Println(producer.Flush(10000))

	fmt.Println("Producer: сообщение отправлено, ждем подтверждения...", err)

	// Ждем подтверждения доставки
	event := <-deliveryChan
	msgDelivery := event.(*kafka.Message)
	if msgDelivery.TopicPartition.Error != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", msgDelivery.TopicPartition.Error)
	} else {
		fmt.Println("Сообщение успешно отправлено")
	}

	// Создаем консюмера
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "test",
		"auto.offset.reset": "earliest", // Ключевой момент
	})
	if err != nil {
		fmt.Println("Ошибка создания консюмера:", err.Error())
		return
	}
	defer consumer.Close()

	fmt.Println("Подписываемся на топик:", topic)
	if err := consumer.Subscribe(topic, nil); err != nil {
		fmt.Println("Ошибка подписки:", err.Error())
		return
	}

	// Читаем сообщения из темы
	fmt.Println("Консумер: Ждем сообщения...")
	for {
		msg, err := consumer.ReadMessage(5 * time.Second) // Ждем 5 секунд на сообщение
		if err == nil {
			fmt.Printf("Полученное сообщение: %s\n", msg.Value)
			break // Если вы хотите читать продолжительно, удалите этот break
		} else if err.(kafka.Error).Code() == kafka.ErrTimedOut {
			// Если время ожидания истекло, продолжаем ожидать
			continue
		} else {
			fmt.Println("Ошибка чтения сообщения:", err.Error())
			break
		}
	}

	// Ожидание перед завершением
	time.Sleep(1 * time.Second)
}
