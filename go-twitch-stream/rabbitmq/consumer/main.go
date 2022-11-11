package main

import (
	"fmt"
	"log"
	"os/exec"
	s "strings"
	mq "test.com/m/go-twitch-stream/rabbitmq"
	"test.com/m/internal/database"
)

func listenStream(stream string) {
	log.Printf("Connecting to the chat for stream %s", stream)
	// https://serverfault.com/a/903631
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo $$; exec ./stream %s", stream))
	err := cmd.Start()
	if err != nil {
		log.Print("{fuck")
	}
	database.InsertStreamEvent(stream, true, cmd.Process.Pid)
	// err := cmd.Wait()
	// if err != nil {
	//     log.Printf("Error on stream command")
	//     log.Print(err)
	// }
}

func main() {

	rabbitMQConnection := mq.ConnectToRabbitMQ()
	defer rabbitMQConnection.Close()
	rabbitMQChannel := mq.ConnectToRabbitMQChannel(rabbitMQConnection)
	_, err := rabbitMQChannel.QueueDeclare(
		"Streams", // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}
	defer rabbitMQChannel.Close()

	// Subscribing to Streams for getting messages.
	messages, err := rabbitMQChannel.Consume(
		"Streams", // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		log.Println(err)
	}

	// Build a welcome message.
	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			log.Printf(" > Received message: %s\n", message.Body)
			command := string(message.Body)
			if s.Contains(command, "start") {
				stream := s.Split(command, " ")[1]
				listenStream(stream)
			}
			if s.Contains(command, "stop") {
				stream := s.Split(command, " ")[1]
				pid := database.GetLatestPID(stream)
				cmd := exec.Command("kill", fmt.Sprint(pid))
				cmd.Run()
				cmd.Wait()
				database.InsertStreamEvent(stream, false, pid)
			}

		}
	}()

	<-forever
}
