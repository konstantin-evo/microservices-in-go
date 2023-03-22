package event

import (
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type Consumer struct {
	conn          *amqp.Connection
	queueName     string
	logServiceURL string
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func NewConsumer(conn *amqp.Connection, logServiceURL string) (Consumer, error) {
	consumer := Consumer{
		conn:          conn,
		logServiceURL: logServiceURL,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
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

	err = bindTopicsToQueue(channel, queue, topics)
	if err != nil {
		return err
	}

	return consumeMessages(channel, queue, consumer.logServiceURL)
}

func bindTopicsToQueue(channel *amqp.Channel, queue amqp.Queue, topics []string) error {
	for _, topic := range topics {
		err := channel.QueueBind(
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

	return nil
}

func consumeMessages(channel *amqp.Channel, queue amqp.Queue, logServiceURL string) error {
	messages, err := channel.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload, logServiceURL)
		}
	}()

	log.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", queue.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload, logServiceURL string) {
	switch payload.Name {
	case "log", "event":
		err := logEvent(payload, logServiceURL)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// you can have as many cases as you want, as long as you write the logic
	default:
		err := logEvent(payload, logServiceURL)
		if err != nil {
			log.Println(err)
		}
	}
}
