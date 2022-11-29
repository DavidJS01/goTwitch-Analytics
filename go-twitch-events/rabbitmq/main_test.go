package rabbitmq

import (
	"testing"
	"fmt"
	"os"

)

func TestConnectToRabbitMQ(t *testing.T) {
	if os.Getenv("INTEGRATION") == "" {
		t.Skip("Skipping integration test")
	}
	url :=  fmt.Sprintf("amqp://guest:guest@%s:%s/", os.Getenv("RABBITMQ_HOST"), os.Getenv("RABBITMQ_PORT"))         

	t.Setenv("AMQP_SERVER_URL", url)
	t.Log(os.Getenv("AMQP_SERVER_URL"))

	conn := ConnectToRabbitMQ()
	fmt.Print(conn.RemoteAddr())
	
}