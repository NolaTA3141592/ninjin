package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ninjin/util/SlackUtil"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/go-yaml/yaml"
	"github.com/slack-go/slack"
)

type SlackEvent struct {
	Type	string `json:"type"`
	Channel	string `json:"channel"`
	User 	string `json:"user"`
	Text	string `json:"text"`
}

func main() {
	config_buf, err := os.ReadFile("./config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var config_data map[string]string
	err = yaml.Unmarshal(config_buf, &config_data)
	var (
		SLACK_API_TOKEN = config_data["SLACK_API_TOKEN"]
		SLACK_VERIFY_TOKEN = config_data["SLACK_VERIFY_TOKEN"]
		DISCORD_API_TOKEN = config_data["DISCORD_API_TOKEN"]
		test_channel_id = config_data["test_channel_id"]
	)
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

		
		switch event.Type {
			case "url_verification":
				SlackUtil.Verify(w, r, body, SLACK_VERIFY_TOKEN)
			case "event_callback":
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