package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusMethodNotAllowed) //there should never be a get request
		return
	case "POST":
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			panic(err)
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		var veriTest urlVerification
		err = json.Unmarshal(body, &veriTest)
		if err != nil {
			panic(err)
		}

		if veriTest.Challenge != "" {
			fmt.Println("Bot is being verified by Slack\n")
			challengeResp := veriTest.Challenge
			w.Write([]byte(challengeResp))
			return

		}
		var msg slackEvent
		err = json.Unmarshal(body, &msg)

		if err != nil {
			panic(err)
		}

		eventType := msg.Event.Type
		fmt.Println("\n" + eventType + " is the event type")

		//probably shouldn't hardcode the bot's key
		api := slack.New("xoxb-3796174298896-3769530677237-FFoh89wMBAcp79s09YDq9Bjm")

		if eventType == "app_mention" {
			messageText := msg.Event.Text
			//same with hardcoding the channel ID
			channelID, timestamp, err := api.PostMessage("C03N9Q7P5JB", slack.MsgOptionText(messageText, false))

			if err != nil {
				fmt.Printf("%s\n", err)
				return
			}
			fmt.Printf("\nMessage sent successfully to channel %s at %s \n", channelID, timestamp)
		}
	}

}

func main() {

	http.HandleFunc("/", handler)

	fmt.Printf("Starting sever for testing on port 8080...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
