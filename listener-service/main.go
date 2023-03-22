package main

import (
	"fmt"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
)

const (
	maxAttempts = 5
)

type config struct {
	RabbitMqURL   string
	LogServiceURL string
	Topics        []string
}

func main() {
	// Load configuration
	app := config{
		RabbitMqURL:   "amqp://guest:guest@rabbitmq",
		LogServiceURL: "http://logger-service/log",
		Topics:        []string{"log.INFO", "log.WARNING", "log.ERROR"},
	}

	// Try to connect to RabbitMQ
	connection, err := connect(app.RabbitMqURL)
	if err != nil {
		log.Panic(err)
	}
	defer connection.Close()

	// Create consumer
	consumer, err := event.NewConsumer(connection, app.LogServiceURL)
	if err != nil {
		panic(err)
	}

	// Start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// Watch the queue and consume events
	err = consumer.Listen(app.Topics)
	if err != nil {
		log.Println(err)
	}
}

func connect(RabbitMqURL string) (*amqp.Connection, error) {
	var (
		counts  int64
		backOff = 1 * time.Second
	)

	for {
		connection, err := amqp.Dial(RabbitMqURL)
		if err == nil {
			log.Println("Connected to RabbitMQ!")
			return connection, nil
		}

		if counts >= maxAttempts {
			return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %v", maxAttempts, err)
		}

		counts++
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Failed to connect to RabbitMQ. Retrying in %v...", backOff)
		time.Sleep(backOff)
	}
}
