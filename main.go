package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/slack-go/slack"
)

type SlackEvent struct {
	Type	string `json:"type"`
	Channel	string `json:"channel"`
	User 	string `json:"user"`
	Text	string `json:"text"`
}

const (
	SLACK_API_TOKEN = "xoxb-6567041518326-7037473636820-KUTMSWnPrV4X3ngn5oyOqr41"
	SLACK_VERIFY_TOKEN = "VeZpxISb8b0yqZwxQr2OX5ol"
	DISCORD_API_TOKEN = "MTIwMjg3ODc0MDA5OTc2ODM1MQ.GsHYhv.6LAmr-TwUDmDCQNGiDgNuaxqUUA3qVVlaSL0x4"
	test_channel_id = "1202886521842049056"
)


func varify(w http.ResponseWriter, r *http.Request, body []byte) {
	var jsonbody map[string]interface{}
	err := json.Unmarshal(body, &jsonbody)
	if err != nil {
		http.Error(w, "Failed to parse json", http.StatusBadRequest)
		return
	}

	token, ok := jsonbody["token"].(string)
	if !ok {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	eventType, ok := jsonbody["type"].(string)
	if !ok {
		http.Error(w, "Missing event type", http.StatusBadRequest)
		return
	}

	if token != SLACK_VERIFY_TOKEN {
		http.Error(w, "Invalid verification token", http.StatusUnauthorized)
		return
	}
	
	if eventType == "url_verification" {
		challenge, ok := jsonbody["challenge"].(string)
		if !ok {
			http.Error(w, "Missing challenge", http.StatusBadRequest)
			return
		}
		fmt.Fprint(w, challenge)
	}
	w.WriteHeader((http.StatusOK))
}

func main() {
	discord, err := discordgo.New("Bot " + DISCORD_API_TOKEN)
	if err != nil {
		fmt.Println("Error creating discord bot : ", err)
		return
	}

	api := slack.New(SLACK_API_TOKEN)

	SlackEventsEndPoint := "/slack/events"
	
	http.HandleFunc(SlackEventsEndPoint, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		var event SlackEvent
		err = json.Unmarshal(body, &event)
		if err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}

		var jsonbody map[string]interface{}
		err = json.Unmarshal(body, &jsonbody)
		if err != nil {
			http.Error(w, "Failed to parse json", http.StatusBadRequest)
			return
		}

		if event.Type == "url_verification" {
			varify(w, r, body)
			return
		}

		if event.Type == "event_callback" {
			jsonbody2, ok := jsonbody["event"].(map[string]interface{})
			if !ok {
				http.Error(w, "Failed to parse json", http.StatusBadRequest)
				return
			}
			if jsonbody2["type"] == "message" {
				userid := jsonbody2["user"]
				userinfo, err := api.GetUserInfo(userid.(string))	
				if err != nil {
					http.Error(w, "Failed to get user info", http.StatusInternalServerError)
					return
				}
				discord.ChannelMessageSend(test_channel_id, fmt.Sprintf("[%s]: %s", userinfo.RealName, jsonbody2["text"]))
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	})


	go func() {
		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			fmt.Println("Error stargting HTTP server:", err)
		}
	}()

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return
	}
	defer discord.Close()

	select {}
}