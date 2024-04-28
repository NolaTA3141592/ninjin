package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type SlackEvent struct {
	Type	string `json:"type"`
	Channel	string `json:"channel"`
	User 	string `json:"user"`
	Text	string `json:"text"`
}

func main() {
	SLACK_API_TOKEN := "xoxb-6567041518326-6586379363121-ENqFZK3Uj1plShdRPOyffVIm"
	SLACK_VERIFY_TOKEN := SLACK_API_TOKEN
	DISCORD_API_TOKEN := "MTIwMjg3ODc0MDA5OTc2ODM1MQ.GsHYhv.6LAmr-TwUDmDCQNGiDgNuaxqUUA3qVVlaSL0x4"
	test_channel_id := "1202886521842049056"

	discord, err := discordgo.New("Bot " + DISCORD_API_TOKEN)
	if err != nil {
		fmt.Println("Error creating discord bot : ", err)
		return
	}

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

		if event.Type == "message" {
			discord.ChannelMessageSend(test_channel_id, fmt.Sprintf("[%s]: %s", event.User, event.Text))
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	SlackEventVerifyEndpoint := "/slack/events/verify"
	http.HandleFunc(SlackEventVerifyEndpoint, func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")

		if token != SLACK_VERIFY_TOKEN {
			http.Error(w, "Invalid verification token", http.StatusUnauthorized)
			return
		}

		w.WriteHeader((http.StatusOK))
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