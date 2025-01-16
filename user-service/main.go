package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	kafka_broker "microservice/kafka-broker"
	"microservice/model"
	"net/http"
	"strconv"
)

func registerUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInBytes, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to convert into bytes user": err.Error()})
		return
	}

	err = kafka_broker.PushUserToQueue("user_register", userInBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Failed to push user to queue": err.Error()})
		return
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     strconv.FormatInt(user.ID, 10) + user.Username + "was successfully registered",
	}

	c.JSON(http.StatusOK, response)
	return
}

func main() {

	router := gin.Default()

	router.POST("/user", registerUser)

	// Запуск горутины для потребления сообщений
	go kafka_broker.ConsumeMessage("user_register")

	router.Run(":8082")

}
