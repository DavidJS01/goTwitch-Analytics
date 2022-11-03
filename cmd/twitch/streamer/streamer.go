package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"io"
)


func GetLivestreamInfo(streamerName string){
	client := &http.Client{}
	clientID := os.Getenv("twitchClientId")
	oauth := strings.Split(os.Getenv("twitchAuth"), ":")[1]
	fmt.Print(fmt.Sprintf("https://api.twitch.tv/helix/streams?%s", streamerName))
	request, _ := http.NewRequest("GET", fmt.Sprintf("https://api.twitch.tv/helix/streams?%s", streamerName), nil)

	request.Header = http.Header{
		"Client-ID": {clientID},
		"Authorization": {fmt.Sprintf("Bearer %s", oauth)},
	}
	log.Print(request.Header)

	response, err := client.Do(request)
	if err != nil {
		log.Print("Error when calling request")
		log.Print(err)
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	fmt.Println(string(body))
}


func main() {
	GetLivestreamInfo("tyler1lol")
}