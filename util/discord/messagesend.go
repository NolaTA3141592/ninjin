package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"ninjin/util/slack"
)

type WebhookMessage struct {
	Content 		string `json:"content"`
	Username		string `json:"username"`
	avatar	 		string `json:"avatar_url"`
}

func (r *Router) MessageSend(user *slack.User, msg string) {
	webhook := r.webhooks[0]
	WebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.TOKEN)
	usericon := "https://github.com/qiita.png"
	fmt.Println(usericon)
	whm := WebhookMessage {
		Content: msg,
		Username: user.RealName,
		avatar: usericon,
	}
	msgByte, err := json.Marshal(whm)
	if err != nil {
		fmt.Println("Error marshaling message : ", err)
		return
	}

	resp, err := http.Post(WebhookURL, "application/json", bytes.NewBuffer(msgByte))
	if err != nil {
		fmt.Println("Error sending message : ", err)
	}
	defer resp.Body.Close()
}