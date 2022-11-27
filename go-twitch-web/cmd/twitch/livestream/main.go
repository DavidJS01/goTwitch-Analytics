package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/url"
	"os"
	"regexp"
	s "strings"
	"test.com/m/internal/database"
)

func parseUserName(twitchMessage string) string {
	rx := regexp.MustCompile(`:(.*?)!`)
	username := rx.FindAllStringSubmatch(twitchMessage, -1)[0][0]
	username = s.Trim(username, ":")
	username = s.Trim(username, "!")
	return username
}

func parseMessage(twitchMessage string) string {
	rx := regexp.MustCompile(`#(.*?):(.*)`)
	message := rx.FindAllStringSubmatch(twitchMessage, -1)[0][2]
	message = s.Trim(message, "\n")
	return message
}

func createWebSocketClient(host string, scheme string) *websocket.Conn {
	log.Print("Creating websocket client")
	u := url.URL{Scheme: scheme, Host: host}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func authenticateClient(connection *websocket.Conn, twitchChannel string) {
	log.Print("authenticating websocket client")
	oauth := fmt.Sprintf("PASS %s", os.Getenv("twitchAuth"))
	username := fmt.Sprintf("NICK %s", os.Getenv("twitchUsername"))

	err := connection.WriteMessage(websocket.TextMessage, []byte(oauth))
	if err != nil {
		log.Println(err)
	}
	err = connection.WriteMessage(websocket.TextMessage, []byte(username))
	if err != nil {
		log.Println(err)
	}
	connection.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("JOIN #%s", twitchChannel)))
	if err != nil {
		log.Println(err)
	}
}


func parseTwitchMessage(message []byte, channel string, connection *websocket.Conn) {
	messageString := string(message)
	if s.Contains(messageString, "PRIVMSG") {
		message := parseMessage(messageString)
		username := parseUserName(messageString)
		fmt.Printf("%s: %s \n", username, messageString)
		database.InsertTwitchMesasge(username, message, channel)
	}
	if s.Contains(messageString, "PING") {
		connection.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv"))
	}
}


func receiveHandler(connection *websocket.Conn, channel string) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error while recieving a twitch message:", err)
			return
		}
		parseTwitchMessage(msg, channel, connection)
	}
}

func StartStream(twitch_channel string) {
	connection := createWebSocketClient("irc-ws.chat.twitch.tv:443", "wss")
	authenticateClient(connection, twitch_channel)
	receiveHandler(connection, twitch_channel)
	defer connection.Close()
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var twitch_channel string = os.Args[1]
	connection := createWebSocketClient("irc-ws.chat.twitch.tv:443", "wss")
	authenticateClient(connection, twitch_channel)
	fmt.Print("got to the handler")
	receiveHandler(connection, twitch_channel)
	defer connection.Close()
}
