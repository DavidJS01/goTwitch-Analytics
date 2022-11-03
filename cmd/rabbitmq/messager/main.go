package main 

import (
    mq "test.com/m/cmd/rabbitmq"
    "github.com/gorilla/mux"
    "net/http"
    "fmt"

)


func WriteMessage(body []byte) error {
	rabbitMQConnection := mq.ConnectToRabbitMQ()
    defer rabbitMQConnection.Close()
    rabbitMQChannel := mq.ConnectToRabbitMQChannel(rabbitMQConnection)
    defer rabbitMQChannel.Close()


    // With the instance and declare Queues that we can
    // publish and subscribe to.
    _, err := rabbitMQChannel.QueueDeclare(
        "Streams", // queue name
        true,            // durable
        false,           // auto delete
        false,           // exclusive
        false,           // no wait
        nil,             // arguments
    )
    if err != nil {
        panic(err)
    }

    message := mq.CreateMessage(body)

	// Attempt to publish a message to the queue.
	if err := rabbitMQChannel.Publish(
		"",              // exchange
		"Streams", // queue name
		false,           // mandatory
		false,           // immediate
		message,         // message to publish
	); err != nil {
		return err
	}

	return nil
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	state := mux.Vars(r)["state"]
    channel := mux.Vars(r)["channel"]
    WriteMessage([]byte(fmt.Sprintf("%s %s", state, channel)))
    w.Write([]byte("msg sent"))
	
}

func main() {

    mux := mux.NewRouter()
    mux.HandleFunc("/message", sendMessageHandler).Queries("state", "{state}", "channel", "{channel}").Methods("POST")
    http.ListenAndServe(":9090", mux)
}

