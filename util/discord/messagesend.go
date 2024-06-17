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

func (r *Router) MessageSend(user *slack.User, msg *cls.Message) {
	webhook := r.webhooks[0]
	WebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s?wait=true", webhook.ID, webhook.TOKEN)
	whm := WebhookMessage {
		Content: msg.Content,
		Username: user.RealName,
		Avatar: user.Usericon,
	}

	if msg.FileURL != "" && msg.FileName != ""{
		whm.Embeds = []Embed{
			{
				Title: msg.FileName,
				URL: msg.FileURL,
			},
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

	// fmt.Println(resp.Body)
	var respbody WebhookResponse
    if err := json.NewDecoder(resp.Body).Decode(&respbody); err != nil {
		fmt.Println("Error Decode Discord Response : ", err)
        return
    }

	msg.Discord_ID = respbody.ID
}