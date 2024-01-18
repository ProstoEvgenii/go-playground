package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresConnection *gorm.DB

func InitPostgres() {
	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASS"),
		os.Getenv("PG_DB"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"))
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal(err)
	}
	postgresConnection = conn

	if err := conn.AutoMigrate(
		&Person{},
	); err != nil {
		log.Fatal(err)
	}
	log.Println("Postgres is running.")

}

func GetPostgres() *gorm.DB {
	return postgresConnection
}

func Close() {
	db, err := postgresConnection.DB()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
}
