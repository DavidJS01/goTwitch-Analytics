package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"os"
	z "strings"
	"testing"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		fmt.Print(message)
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestParseUserName(t *testing.T) {
	mock_message := ":username!username@username.tmi.twitch.tv PRIVMSG #katevolved :message here"
	username := parseUserName(mock_message)
	if username != "username" {
		t.Errorf("Parsing a username was incorrect, expected 'username' got '%s'", username)
	}
}

func TestParseMessage(t *testing.T) {
	mock_message := ":username!username@username.tmi.twitch.tv PRIVMSG #katevolved :message here"
	message := parseMessage(mock_message)
	if message != "message here" {
		t.Errorf("Parsing a message was incorrect, expected 'message here' got '%s'", message)
	}
}

func TestCreateWebsocketClient(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(echo))
	address := z.Split(s.URL, "http://")[1]
	defer s.Close()

	client := createWebSocketClient(address, "ws")
	fmt.Print(client.RemoteAddr().String())
	if client.RemoteAddr().String() != address {
		t.Errorf("Error connecting test ws client to 127.0.0, server url %s, client url %s", address, client.RemoteAddr().String())
	}
	s.Close()
}

func TestAuthenticateClient(t *testing.T) {
	// temporarily modify environment variables
	t.Setenv("twitchAuth", "password")
	t.Setenv("twitchUsername", "username")

	// create websocket server and client
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()
	address := z.Split(s.URL, "http://")[1]
	client := createWebSocketClient(address, "ws")
	authenticateClient(client, "twitch_channel_here")

	// test for `PASS <oauth_password> message`
	_, msg, _ := client.ReadMessage()
	fmt.Print(string(msg))
	if string(msg) != fmt.Sprintf("PASS %s", os.Getenv("twitchAuth")) {
		t.Errorf("expected to recieve `PASS %s`, got %s", os.Getenv("twitchAuth"), string(msg))
	}

	// test for `NICK <username> message`
	_, msg, _ = client.ReadMessage()
	fmt.Print(string(msg))
	if string(msg) != fmt.Sprintf("NICK %s", os.Getenv("twitchUsername")) {
		t.Errorf("expected to recieve NICK %s, got %s", os.Getenv("twitchUsername"), string(msg))
	}

	// test for `JOIN #` message
	_, msg, _ = client.ReadMessage()
	fmt.Print(string(msg))
	if string(msg) != "JOIN #twitch_channel_here" {
		t.Errorf("expected to recieve JOIN #twitch_channel_here, got %s", string(msg))
	}

}


func mockInsertTwitchMessage(username string, message string, channel string) {
	return
}

func TestParseTwitchMessage(t *testing.T){
	// create messages to test on
	mockTwitchMessage := []byte("username!username@username.tmi.twitch.tv PRIVMSG #katevolved :message here")
	mockPingMessage := []byte("PING server message here")
	// start testing websocket server
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()
	// create client websocket conn
	address := z.Split(s.URL, "http://")[1]
	client := createWebSocketClient(address, "ws")

	// run twitch message test
	parseTwitchMessage(mockTwitchMessage, "Katevolved", client, mockInsertTwitchMessage)
	// run ping message test
	parseTwitchMessage(mockPingMessage, "Katevolved", client, mockInsertTwitchMessage)
	
	_, msg, _ := client.ReadMessage()
	if string(msg) != "PONG :tmi.twitch.tv" {
		t.Errorf("Unexpected PING response, expected PONG got %s", string(msg))
	}



	
}
