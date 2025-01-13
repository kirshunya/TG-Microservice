package internal

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafka/producer/model"
	"log"
)

func CreateProducer(brokers string) (*kafka.Producer, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers:": brokers,
	})
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func HandleDelivery(producer *kafka.Producer) {
	deliveryChan := producer.Events()

	go func() {
		for e := range deliveryChan {
			switch msg := e.(type) {
			case *kafka.Message:
				if msg.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v", msg.TopicPartition.Error)
				} else {
					log.Printf("Message delivered to %v [%d] at offset %v",
						*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset)
				}
			}
		}
	}()
}

func SendMessage(producer *kafka.Producer, topic string, user model.User) error {
	value, err := json.Marshal(user)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key: nil,
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: value,
	}

	if err = producer.Produce(&msg, nil); err != nil {
		return err
	}
	return nil
}
