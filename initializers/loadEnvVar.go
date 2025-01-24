package initializers

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv(filename string) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}
