package configs

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectMysql() *gorm.DB {
	dsn := os.Getenv("GCP_PRIVATE_DATABASE_INFO") // private IP
	//dsn := os.Getenv("GCP_PUBLIC_DATABASE_INFO") // public IP
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Failed to connect Database")
	}
	log.Println("Database connection success")

	return db
}
