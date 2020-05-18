package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/sklevenz/cf-api-broker/server"
)

func main() {
	fmt.Println("-- start application")

	server := http.HandlerFunc(server.BrokerServer)
	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
