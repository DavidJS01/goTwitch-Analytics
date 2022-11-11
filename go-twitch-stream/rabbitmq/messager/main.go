package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	mq "test.com/m/go-twitch-stream/rabbitmq"
)

type MessageResponse struct {
	State string `json:"state"`
	Channel string `json:"channel"`
	Status_Code int `json:"status_code"`
}


func WriteMessage(body []byte) error {
	rabbitMQConnection := mq.ConnectToRabbitMQ()
	defer rabbitMQConnection.Close()
	rabbitMQChannel := mq.ConnectToRabbitMQChannel(rabbitMQConnection)
	defer rabbitMQChannel.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
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

	message := mq.CreateMessage(body)

	// Attempt to publish a message to the queue.
	if err := rabbitMQChannel.Publish(
		"",        // exchange
		"Streams", // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	); err != nil {
		return err
	}

	return nil
}

func messageResponse(state string, channel string, status_code int) MessageResponse {
	var response MessageResponse
	response.State = state
	response.Channel = channel
	response.Status_Code = status_code

	return response
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	state := mux.Vars(r)["state"]
	channel := mux.Vars(r)["channel"]
	response := messageResponse(state, channel, 200)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
	WriteMessage([]byte(fmt.Sprintf("%s %s", state, channel)))

}

func main() {

	mux := mux.NewRouter()
	mux.HandleFunc("/message", sendMessageHandler).Queries("state", "{state}", "channel", "{channel}").Methods("POST")
	http.ListenAndServe(":9090", mux)
}
