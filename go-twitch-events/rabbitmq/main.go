package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
)

func ConnectToRabbitMQ() *amqp.Connection {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL") //@TODO add this

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	return connectRabbitMQ

}

func ConnectToRabbitMQChannel(connectRabbitMQ *amqp.Connection) *amqp.Channel {
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	return channelRabbitMQ
}

func CreateMessage(body []byte) amqp.Publishing {
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        body,
	}
	return message

}
