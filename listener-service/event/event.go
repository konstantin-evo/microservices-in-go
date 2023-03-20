package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type (routes messages based on their topic routing key)
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // name (given a random name by the server)
		false, // durable (will not survive a broker restart)
		false, // not be automatically deleted when there are no more consumers
		true,  // exclusive (it can only be accessed by the connection that declares it)
		false, // no-wait for a response from the server
		nil,   // there are no additional arguments to be passed to the queue declaration
	)
}
