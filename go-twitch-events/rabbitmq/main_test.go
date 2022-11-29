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
	t.Setenv("AMQP_SERVER_URL", os.Getenv("AMQP_SERVER_URL"))
	conn := ConnectToRabbitMQ()
	fmt.Print(conn.RemoteAddr())
	
}