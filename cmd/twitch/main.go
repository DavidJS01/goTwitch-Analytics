package main

import (
	"fmt"
	"github.com/gorilla/websocket"
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

func createWebSocketClient() *websocket.Conn {
	u := url.URL{Scheme: "wss", Host: "irc-ws.chat.twitch.tv:443"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func authenticateClient(connection *websocket.Conn, twitchChannel string) {
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

func receiveHandler(connection *websocket.Conn, channel string) {
	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error in receive:", err)
			return
		}
		if msg != nil {
			// when a message is recieved in the channel, ping and then parse the message
			if s.Contains(string(msg), "PRIVMSG") {
				message := parseMessage(string(msg))
				username := parseUserName(string(msg))
				fmt.Printf("%s: %s \n", username, message)
				database.InsertTwitchMesasge(username, message, channel)
			}
			if s.Contains(string(msg), "PING") {
				connection.WriteMessage(websocket.TextMessage, []byte("PONG :tmi.twitch.tv"))
			}
		}
	}
}

func main() {
	// Channel to indicate that the receiverHandler is done
	var twitch_channel string = os.Args[1]
	connection := createWebSocketClient()
	authenticateClient(connection, twitch_channel)
	receiveHandler(connection, twitch_channel)
	defer connection.Close()
}
