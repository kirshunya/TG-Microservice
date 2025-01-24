package initializers

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"microservice/model"
	"os"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error

	dsn := "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("DB_USER") + " password=" + os.Getenv("DB_PASSWORD") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	fmt.Println(dsn)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to DB")
	}
}

func SyncDB() {
	DB.AutoMigrate(&model.User{})
}
