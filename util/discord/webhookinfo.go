package discord

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
    DiscordAPIURL = "https://discord.com/api/v9"
)

type Webhook struct {
    ID   string `json:"id"`
    Name string `json:"name"`
	ChannelID string `json:"channel_id"`
	ChannelName string
}

func (r *Router) GetWebhookList() ([]Webhook, error) {
    var webhooks []Webhook
    req, err := http.NewRequest("GET", DiscordAPIURL + "/guilds/" + r.SERVER_ID + "/webhooks", nil)
    if err != nil {
        fmt.Println("Error creating HTTP request:", err)
        return webhooks, err
    }
    req.Header.Set("Authorization", "Bot " + r.DISCORD_API_TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error sending HTTP request:", err)
        return webhooks, err
    }
    defer resp.Body.Close()

    if err := json.NewDecoder(resp.Body).Decode(&webhooks); err != nil {
        fmt.Println("Error decoding JSON:", err)
        return webhooks, err
    }

    for i, webhook := range webhooks {
        channelname, err := r.getChannelName(webhook.ChannelID)
        if err != nil {
            fmt.Println("Error getting channel Name:", err)
            continue
        }
        webhooks[i].ChannelName = channelname
		fmt.Printf("%s", channelname)
    }

	return webhooks, nil
}

func (r *Router) getChannelName(channelID string) (string, error) {
    var channelname string

    req, err := http.NewRequest("GET", DiscordAPIURL+"/channels/"+channelID, nil)
    if err != nil {
        return channelname, err
    }
    req.Header.Set("Authorization", "Bot "+r.DISCORD_API_TOKEN)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return channelname, err
    }
    defer resp.Body.Close()

	var jsonbody map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&jsonbody); err != nil {
        return channelname, err
    }

	return jsonbody["name"].(string), nil
}