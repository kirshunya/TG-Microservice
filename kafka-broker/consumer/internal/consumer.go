package internal

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

func CreateConsumer(brokers, topic string) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          "group",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create consumer %s", err)
		return nil, err
	}

	if err := consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		log.Fatalf("Failed to subscribe to topic: %s", err)
		return nil, err
	}

	return consumer, nil
}

func ConsumeMessages(consumer *kafka.Consumer) {
	defer consumer.Close()

	// Основной цикл для получения сообщений
	for {
		msg, err := consumer.ReadMessage(-1) // -1 - блокировка до получения сообщения
		if err == nil {
			fmt.Printf("Received message: %s\n", string(msg.Value))
		} else {
			log.Printf("Error while receiving message: %v\n", err)
		}
	}
}
