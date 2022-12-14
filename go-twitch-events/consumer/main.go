package main

import (
	"fmt"
	"log"
	"os/exec"
	s "strings"
	mq "test.com/m/go-twitch-events/rabbitmq"
	"test.com/m/internal/database"
)

func listenStream(stream string) {
	log.Printf("Connecting to the chat for stream %s", stream)
	// https://serverfault.com/a/903631
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo $$; exec ./stream %s", stream))
	err := cmd.Start()
	if err != nil {
		log.Print(err)
	}
	database.InsertStreamEvent(stream, true, cmd.Process.Pid)
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

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages")

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf(" > Received message: %s\n", message.Body)
			command := string(message.Body)
			if s.Contains(command, "start") {
				stream := s.Split(command, " ")[1]
				err := database.InsertStreamer(stream, true)
				if err != nil {
					log.Print(err)
				}
				listenStream(stream)
			}
			if s.Contains(command, "stop") {
				stream := s.Split(command, " ")[1]
				pid := database.GetLatestPID(stream)
				cmd := exec.Command("kill", fmt.Sprint(pid))
				err := cmd.Run()
				if err != nil {
					log.Print(err)
				}
				err = cmd.Wait()
				if err != nil {
					log.Print(err)
				}
				database.InsertStreamEvent(stream, false, pid)
				err = database.InsertStreamer(stream, false)
				if err != nil {
					log.Print(err)
				}
			}

		}
	}()

	<-forever
}
