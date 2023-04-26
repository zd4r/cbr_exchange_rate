package main

import (
	"log"

	"github.com/zd4r/cbr_exchange_rate/internal/app"
	"github.com/zd4r/cbr_exchange_rate/pkg/config"
)

func main() {
	// Configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Config error: %v", err)
	}

	// Run app
	app.Run(cfg)
}
