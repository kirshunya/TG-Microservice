package main

import (
	"kafka/producer/internal"
	"kafka/producer/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	brokers := "localhost:9092"

	producer, err := internal.CreateProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create producer %s", err)
	}

	defer producer.Close()

	internal.HandleDelivery(producer)

	router.POST("/registration", func(c *gin.Context) {
		var user model.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		topic := "user-registration"

		if err := internal.SendMessage(producer, topic, user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "user sent"})
	})

	router.Run(":8081")
}
