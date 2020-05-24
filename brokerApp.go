package main

import (
	"log"
	"net/http"

	"github.com/sklevenz/cf-api-broker/server"
)

func main() {
	log.Println("start application")

	server := server.NewBrokerServer()

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
