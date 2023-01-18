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
	rx := regexp.MustCompile(`(\w*)!`)
	username := rx.FindString(twitchMessage)
	username = s.Split(username, "!")[0]
	return username
}

func parseMessage(twitchMessage string) string {
	rx := regexp.MustCompile(`#(.*?):(.*)`)
	message := rx.FindAllStringSubmatch(twitchMessage, -1)[0][2]
	message = s.Trim(message, "\n")
	return message
}

func createWebSocketClient(host string, scheme string) (*websocket.Conn, error) {
	log.Print("Creating websocket client")
	u := url.URL{Scheme: scheme, Host: host}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return c, nil
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

func parseTwitchMessage(message []byte, channel string, connection *websocket.Conn, insertMessage database.InsertMessage) {
	messageString := string(message)
	if s.Contains(messageString, "PRIVMSG") {
		message := parseMessage(messageString)
		username := parseUserName(messageString)
		fmt.Printf("%s: %s \n", username, messageString)
		database.InsertStreamer(channel)
		insertMessage(username, message, channel)
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
		} else {
			parseTwitchMessage(msg, channel, connection, database.InsertTwitchMessage)
		}
	}
}

func StartStream(twitch_channel string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	connection, err := createWebSocketClient("irc-ws.chat.twitch.tv:443", "wss")
	if err != nil {
		log.Fatalf("Error establishing ws client: %s", err)
	}
	authenticateClient(connection, twitch_channel)
	receiveHandler(connection, twitch_channel)
	defer connection.Close()
}

func main() {
	// // database.SetupPostgres()
	// // database.InsertStreamer("katevolved")
	// // database.UpsertStreamEvent("katevolved")
	// database.InsertStreamEventStatus(true, 124, "katevolved")
	// database.UpdateStreamEventStatus(1234, "katevolved")
	// StartStream(os.Args[1])
	database.InsertStreamEventStatus(true, 1, "test")
}
