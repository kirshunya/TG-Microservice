package main

import (
	"github.com/gin-gonic/gin"
	"kafka/consumer/internal"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	brokers := "localhost:9092"
	topic := "user-registration"

	router := gin.Default()

	consumer, err := internal.CreateConsumer(brokers, topic)
	if err != nil {
		log.Fatal(err)
	}

	go internal.ConsumeMessages(consumer)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigint
		log.Println("Shutting down consumer...")
		consumer.Close()
		os.Exit(0)
	}()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "running"})
	})

	router.Run(":8081")
}
