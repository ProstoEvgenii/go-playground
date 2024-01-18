package main

import (
	"log"
	"rest-service/db"
	"rest-service/server"

	"github.com/joho/godotenv"
)

func init() {

	if err := godotenv.Load(); err != nil {
		log.Println("Error during loading .env")
	}
}

func main() {

	db.InitPostgres()
	defer db.Close()

	server.Start()
}
