package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
)

//types should be in separate type.go file
type urlVerification struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

type slackEvent struct {
	Token    string `json:"token"`
	TeamID   string `json:"team_id"`
	APIAppID string `json:"api_app_id"`
	Event    struct {
		ClientMsgID string `json:"client_msg_id"`
		Type        string `json:"type"`
		Text        string `json:"text"`
		User        string `json:"user"`
		Ts          string `json:"ts"`
		Team        string `json:"team"`
		Blocks      []struct {
			Type     string `json:"type"`
			BlockID  string `json:"block_id"`
			Elements []struct {
				Type     string `json:"type"`
				Elements []struct {
					Type   string `json:"type"`
					UserID string `json:"user_id"`
				} `json:"elements"`
			} `json:"elements"`
		} `json:"blocks"`
		Channel string `json:"channel"`
		EventTs string `json:"event_ts"`
	} `json:"event"`
	Type           string `json:"type"`
	EventID        string `json:"event_id"`
	EventTime      int    `json:"event_time"`
	Authorizations []struct {
		EnterpriseID        interface{} `json:"enterprise_id"`
		TeamID              string      `json:"team_id"`
		UserID              string      `json:"user_id"`
		IsBot               bool        `json:"is_bot"`
		IsEnterpriseInstall bool        `json:"is_enterprise_install"`
	} `json:"authorizations"`
	IsExtSharedChannel bool   `json:"is_ext_shared_channel"`
	EventContext       string `json:"event_context"`
}

func send(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		fmt.Println("GET request")
		return
	case "POST":

		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			fmt.Println("ERROR\n")
			return
		}
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(body))

		var veriTest urlVerification
		err = json.Unmarshal(body, &veriTest)

		if err != nil {
			//panic(err)
			fmt.Println("error")
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

	http.HandleFunc("/", send)

	fmt.Printf("Starting sever for testing on port 8080...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
