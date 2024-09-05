package main

import (
	"TaskMaster/configs/config"
	"TaskMaster/pkg/apiserver"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	api := apiserver.New(cfg)
	
	if err := api.Start(); err != nil {
		log.Fatalf("Failed to start the sevrer: %v", err)
	}
}
