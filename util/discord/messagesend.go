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
}

type WebhookResponse struct {
	ID 				string `json:"id"`
}

func (r *Router) MessageSend(user *slack.User, msg *cls.Message) (string) {
	webhook := r.webhooks[0]
	WebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s?wait=true", webhook.ID, webhook.TOKEN)
	whm := WebhookMessage {
		Content: msg.Content,
		Username: user.RealName,
		Avatar: user.Usericon,
	}
	msgByte, err := json.Marshal(whm)
	if err != nil {
		fmt.Println("Error marshaling message : ", err)
		return ""
	}

	resp, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(msgByte))
	if err != nil {
		fmt.Println("Error sending message : ", err)
		return ""
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Body)
	var respbody WebhookResponse
    if err := json.NewDecoder(resp.Body).Decode(&respbody); err != nil {
		fmt.Println("Error Decode Discord Response : ", err)
        return ""
    }

	return respbody.ID
}