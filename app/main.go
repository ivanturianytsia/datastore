package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	mailgun "github.com/mailgun/mailgun-go"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to read .env")
	}

	// It is here because it have to load after .env
	mg = mailgun.NewMailgun(
		os.Getenv("MAILGUN_DOMAIN"),
		os.Getenv("MAILGUN_KEY"),
		os.Getenv("MAILGUN_PUBKEY"))

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
