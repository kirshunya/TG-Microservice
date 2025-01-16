package kafka_broker

import (
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}

func ConsumeMessage(topic string) {

	msgCnt := 0

	worker, err := ConnectConsumer([]string{"localhost:9093"})
	if err != nil {
		panic(err)
	}

	consumer, err := worker.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	fmt.Println("Consumer connected")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{})
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCnt++
				fmt.Printf("Message %d: | Topic(%s) | Message(%s) \n", msgCnt, string(msg.Topic), string(msg.Value))
				user := string(msg.Value)
				fmt.Println(user)
			case sig := <-sigchan:
				fmt.Println("Received terminate signal", sig)
				doneChan <- struct{}{}
			}
		}
	}()

	<-doneChan
	fmt.Println("Processed " + strconv.Itoa(msgCnt) + " messages")

	if err := worker.Close(); err != nil {
		panic(err)
	}
}
