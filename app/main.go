package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to read .env")
	}

	server, err := NewServer()
	if err != nil {
		log.Fatalln(err)
		return
	}

	router := mux.NewRouter()
	server.Route(router)

	log.Printf("Listening on port %s", getPort())
	http.ListenAndServe(getPort(), router)
}
