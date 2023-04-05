package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	WebPort  string
	BrokerURL string
}

func main() {
	// Load configuration from environment variables or command-line arguments
	app, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Starting front end service on port: %s", app.WebPort)
	// Start the HTTP server
    server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
        Handler: routes(app.BrokerURL),
    }

    fmt.Println("Starting front end service on port 80")
    if err := server.ListenAndServe(); err != nil {
        log.Panic(err)
    }
}

func loadConfig() (*Config, error) {
		// Use a default value if the environment variable is not set
		brokerURL, ok := os.LookupEnv("BROKER_URL")
		if !ok {
			brokerURL = "http://localhost:8080"
		}

		config := &Config{
			BrokerURL:                brokerURL,
			WebPort:                  "80",
		}
	
		return config, nil
}
