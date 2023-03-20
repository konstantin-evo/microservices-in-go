package main

import (
	"broker/event"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	WebPort                  string
	AuthenticationServiceURL string
	MailServiceURL           string
	LogServiceURL            string
	RabbitURL                string
	Rabbit                   *amqp.Connection
}

func main() {
	// Load configuration from environment variables or command-line arguments
	config, err := loadConfig()
	if err != nil {
		log.Panic(err)
	}

	// Connect to RabbitMQ
	rabbitMqConnection, err := connect(config.RabbitURL)
	if err != nil {
		log.Panic(err)
	}
	defer rabbitMqConnection.Close()

	// Initialize the app with the configuration
	app := Config{
		WebPort:                  config.WebPort,
		AuthenticationServiceURL: config.AuthenticationServiceURL,
		MailServiceURL:           config.MailServiceURL,
		LogServiceURL:            config.LogServiceURL,
		RabbitURL:                config.RabbitURL,
		Rabbit:                   rabbitMqConnection,
	}

	// Initialize the Consumer with the RabbitMQ connection and Logger instance
	logger := event.NewLogger(app.LogServiceURL)
	consumer := event.NewConsumer(app.Rabbit, logger)

	// Set up topics you want to listen to
	topics := []string{"log", "auth", "event"}

	// Listen to the events in a separate goroutine
	go func() {
		if err := consumer.Listen(topics); err != nil {
			log.Panic(err)
		}
	}()

	// Start the HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	log.Printf("Starting broker service on port %s\n", app.WebPort)
	if err := server.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func connect(url string) (*amqp.Connection, error) {
	// Use an exponential backoff to connect to RabbitMQ
	const maxRetries = 5
	var counts int
	var backOff time.Duration
	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(url)
		if err == nil {
			log.Println("Connected to RabbitMQ!")
			break
		}

		if counts >= maxRetries {
			return nil, err
		}

		counts++

		// calculates the backoff time to wait before attempting to reconnect to RabbitMQ using an exponential strategy
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Failed to connect to RabbitMQ. Retrying in %v...", backOff)
		time.Sleep(backOff)
	}

	return conn, nil
}

func loadConfig() (*Config, error) {
	// Use a default value if the environment variable is not set
	rabbitURL, ok := os.LookupEnv("RABBITMQ_URL")
	if !ok {
		rabbitURL = "amqp://guest:guest@rabbitmq"
	}

	config := &Config{
		RabbitURL:                rabbitURL,
		AuthenticationServiceURL: "http://authentication-service/authenticate",
		MailServiceURL:           "http://mailer-service/send",
		LogServiceURL:            "http://logger-service/log",
		WebPort:                  "80",
	}

	return config, nil
}
