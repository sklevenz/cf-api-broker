package main

import (
	"log"
	"net/http"

	"flag"
	"os"

	"github.com/sklevenz/cf-api-broker/config"
	"github.com/sklevenz/cf-api-broker/server"
)

const (
	staticDir string = "./static"
	defaultPort             = "5000"
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
	flag.StringVar(&configPath, "f", "./config/config.yaml", "path to config file")
}

func main() {

	flag.Parse()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}


	log.Printf("start application on port %v", port)
	log.Printf("version %v", Version)
	log.Printf("commit %v", Commit)

	if err := config.Read(configPath); err != nil {
		log.Fatalf("could not read configuration %v", err)
	}

	server.SetBuildVersion(Version, Commit)
	brokerServer := server.NewRouter(staticDir)

	log.Printf("call server: http://localhost:%v", port)

	if err := http.ListenAndServe(":"+port, brokerServer); err != nil {
		log.Fatalf("could not listen on port %v: %v", port, err)
	}
}
