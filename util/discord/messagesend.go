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
	Avatar	 		string `json:"avatar_url"`
	Embeds   		[]Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Image Image `json:"image"`
}

type Image struct {
	URL string `json:"url"`
}

func (r *Router) MessageSend(user *slack.User, msg string) {
	webhook := r.webhooks[0]
	WebhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", webhook.ID, webhook.TOKEN)
	whm := WebhookMessage {
		Content: msg,
		Username: user.RealName,
		Avatar: user.Usericon,
	}
	
	imageURL := "https://github.com/qiita.png"
	if imageURL != "" {
		whm.Embeds = []Embed{
			{
				Image: Image{
					URL: imageURL,
				},
			},
		}
	}

	msgByte, err := json.Marshal(whm)

	fmt.Println(string(msgByte))
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