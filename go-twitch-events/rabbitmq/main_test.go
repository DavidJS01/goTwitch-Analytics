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
	t.Setenv("AMQP_SERVER_URL", "amqp://guest:guest@message-broker:1337/")
	conn := ConnectToRabbitMQ()
	fmt.Print(conn.RemoteAddr())
	
}