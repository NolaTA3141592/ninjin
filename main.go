package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ninjin/util/cls"
	"ninjin/util/discord"
	"ninjin/util/mdb"
	slacktool "ninjin/util/slack"
	"os"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
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
	if err != nil {
		fmt.Println(err)
		return
	}
	var (
		SLACK_VERIFY_TOKEN = config_data["SLACK_VERIFY_TOKEN"]
	)

	dr := discord.Router {
		DISCORD_API_TOKEN: config_data["DISCORD_API_TOKEN"],
		TEST_CHANNEL_ID: config_data["test_channel_id"],
	}
	err = dr.Setup()
	if err != nil {
		return
	}

	slack := slacktool.SlackUtil {
		SLACK_API_TOKEN: 	config_data["SLACK_API_TOKEN"],
	}

	SlackEventsEndPoint := "/slack/events"
	
	db, err := mdb.Setup()
	if err != nil {
		return
	}
	defer db.Data.Close()

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
				slack.Verify(w, r, body, SLACK_VERIFY_TOKEN)
			case "event_callback":
				jsonbody2, ok := jsonbody["event"].(map[string]interface{})
				if !ok {
					http.Error(w, "Failed to parse json", http.StatusBadRequest)
					return
				}
				if jsonbody2["type"] == "message" {
					user := slacktool.User {
						UserID: jsonbody2["user"].(string),
					}
					msg := cls.Message{}
					slack.AttachUserInfo(&user)
					slack.AttachMessageInfo(&msg, jsonbody2)
					if err != nil {
						http.Error(w, "Failed to get user info", http.StatusInternalServerError)
						return
					}
					dr.EventMassage(&user, &msg, &db)
					db.Insert(&msg)
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

	defer dr.Bot.Close()
	select {}
}