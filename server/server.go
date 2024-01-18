package server

import (
	"log"
	"net/http"
	"os"
	"rest-service/handlers"

	"github.com/gorilla/mux"
)

func Start() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9090"
	}

	router := mux.NewRouter()
	router.HandleFunc("/people", handlers.GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", handlers.GetPerson).Methods("GET")
	router.HandleFunc("/people", handlers.CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", handlers.UpdatePerson).Methods("PUT")
	router.HandleFunc("/people/{id}", handlers.DeletePerson).Methods("DELETE")

	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
