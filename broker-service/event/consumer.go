package event

import (
	"broker/event/data"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	logger    *Logger
	queueName string
}

func NewConsumer(conn *amqp.Connection, logger *Logger) *Consumer {
	return &Consumer{
		conn:   conn,
		logger: logger,
	}
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

func (consumer *Consumer) Listen(topics []string) error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	queue, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		channel.QueueBind(
			queue.Name,
			topic,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload data.Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload, consumer.logger)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", queue.Name)
	<-forever

	return nil
}

func handlePayload(payload data.Payload, logger *Logger) {
	switch payload.Name {
	case "log", "event":
		err := logger.logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	case "auth":
	// you can have as many cases as you want, as long as you write the logic
	default:
		err := logger.logEvent(payload)
		if err != nil {
			fmt.Println(err)
		}
	}
}
