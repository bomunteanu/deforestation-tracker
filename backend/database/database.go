package database

import (
	"deforestation/models"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func Init() {
	var err error
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"))

	for {
		DB, err = gorm.Open("postgres", dsn)
		if err == nil {
			break
		}
		fmt.Println("Failed to connect to the database. Retrying in 5 seconds...")
		time.Sleep(5 * time.Second)
	}

	DB.AutoMigrate(&models.Area{}, &models.History{}, &models.User{})
}

func GetDB() *gorm.DB {
	return DB
}
