package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ninjin/util/cls"
	"ninjin/util/slack"
)

type WebhookMessage struct {
	Content 		string `json:"content"`
	Username		string `json:"username"`
	Avatar	 		string `json:"avatar_url"`
	Embeds			[]Embed `json:"embeds"`
}

type Embed struct {
	Title			string `json:"title"`
	URL				string `json:"url"`
}

type WebhookResponse struct {
	ID 				string `json:"id"`
}

func (r *Router) MessageSend(user *slack.User, msg *cls.Message, webhook *Webhook) {
	WebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s?wait=true", webhook.ID, webhook.TOKEN)
	if msg.ThreadMode {
		WebhookURL += fmt.Sprintf("&thread_id=%s", msg.Discord_thread_ID)
	}
	whm := WebhookMessage {
		Content: msg.Content,
		Username: user.RealName,
		Avatar: user.Usericon,
	}

	for i := 0; i < len(msg.FileNames) && i < len(msg.FileURLs); i++ {
		if msg.FileURLs[i] != "" && msg.FileNames[i] != "" {
			embed := Embed{
				Title: msg.FileNames[i],
				URL:   msg.FileURLs[i],
			}
			whm.Embeds = append(whm.Embeds, embed)
		}
	}

	msgByte, err := json.Marshal(whm)
	if err != nil {
		fmt.Println("Error marshaling message : ", err)
		return
	}

	resp, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(msgByte))
	if err != nil {
		fmt.Println("Error sending message : ", err)
		return
	}
	defer resp.Body.Close()

	var respbody WebhookResponse
    if err := json.NewDecoder(resp.Body).Decode(&respbody); err != nil {
		fmt.Println("Error Decode Discord Response : ", err)
        return
    }

	msg.Discord_ID = respbody.ID
}