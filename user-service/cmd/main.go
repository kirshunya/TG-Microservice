package main

import (
	"github.com/gin-gonic/gin"
	"microservice/initializers"
	"microservice/middleware"
	"microservice/user-service/internal"
)

func init() {
	initializers.LoadEnv("D:\\TG-Microservice\\initializers\\.env")
	initializers.ConnectToDB()
	initializers.SyncDB()
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/signup", internal.SignUp)
	router.POST("/login", internal.LogIn)
	router.GET("/validate", middleware.RequireAuth, internal.Validate)
	router.GET("/recommendation", middleware.RequireAuth, internal.Recommendation)

	router.Run()
}
