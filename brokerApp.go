package main

import (
	"log"
	"net/http"

	"flag"

	"github.com/sklevenz/cf-api-broker/server"
)

const (
	staticDir string = "./static"
)

var (
	// Version set by go build via -ldflags "'-X main.Version=1.0'"
	Version string = "n/a"
	// Commit set by go build via -ldflags "'-X main.Commit=123'"
	Commit string = "n/a"

	// ConfigPath keeps path to configuration file
	configPath string
)

func init() {
	flag.StringVar(&configPath, "f", "./data/config.yaml", "path to config file")
}

func main() {
	flag.Parse()
	log.Println("start application")

	log.Printf("version %v", Version)
	log.Printf("commit %v", Commit)

	server.SetBuildVersion(Version, Commit)
	brokerServer := server.NewRouter(staticDir, configPath)

	if err := http.ListenAndServe(":5000", brokerServer); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
