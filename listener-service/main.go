package main

import (
	"fmt"
	"listener/event"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to rabbitmq
	connection, err := connect()
	if err != nil {
		log.Panic(err)
	}
	defer connection.Close()

	// start listening for messages
	log.Println("Listening for and consuming RabbitMQ messages...")

	// create consumer
	consumer, err := event.NewConsumer(connection)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	// don't continue until rabbit is ready
	for {
		connectToRabbitMQ, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready... ", err)
			counts++
		} else {
			log.Println("Connected to RabbitMQ!")
			connection = connectToRabbitMQ
			break
		}

		if counts > 5 {
			fmt.Println(err)
			return nil, err
		}

		// calculates the backoff time to wait before attempting to reconnect to RabbitMQ using an exponential strategy
		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Println("backing off...")
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
