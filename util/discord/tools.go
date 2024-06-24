package discord

import (
	"fmt"
	"ninjin/util/cls"
	"ninjin/util/slack"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	DISCORD_API_TOKEN	string
	TEST_CHANNEL_ID		string
	Bot  				*discordgo.Session
	SERVER_ID 			string
	webhooks			[]Webhook
}

func (r *Router) Setup() error {
	discordbot, err := discordgo.New("Bot " + r.DISCORD_API_TOKEN)
	if err != nil {
		fmt.Println("Error creating discord bot : ", err)
		return err
	}
	err = discordbot.Open()
	r.Bot = discordbot

	if err != nil {
		fmt.Println("Error opening Discord session:", err)
		return err
	}

	r.SERVER_ID = "1202873405221773343"
	r.webhooks, err = r.GetWebhookList()
	if err != nil {
		fmt.Println("Error getting webhook list : ", err)
	}

	return nil
}

func (r *Router) EventMassage(user *slack.User, msg *cls.Message) {
	webhook, err := r.SelectWebhook(msg)
	if err != nil {
		fmt.Println("Error select Webhook : ", err)
		return
	}
	msg.DiscordChannelID = webhook.ChannelID
	r.MessageSend(user, msg, webhook)
}

func (r *Router) SelectWebhook(msg *cls.Message) (*Webhook, error) {
	for _, webhook := range r.webhooks {
		if webhook.ChannelName == msg.ChannelName {
			return &webhook, nil
		}
	}
	webhook, err := r.CreateWebhook(msg)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

func (r *Router) CreateWebhook(msg *cls.Message) (*Webhook, error) {
	created_channel, err := r.Bot.GuildChannelCreate(r.SERVER_ID, msg.ChannelName, discordgo.ChannelTypeGuildText)		
	if err != nil {
		return nil, err
	}
	webhookName := fmt.Sprintf("ninjin_webhook_of_%s", msg.ChannelName)
	created_webhook, err := r.Bot.WebhookCreate(created_channel.ID, webhookName, "")
	if err != nil {
		return nil, err
	}

	webhook := Webhook {
		ID 		: created_webhook.ID,
		TOKEN 	: created_webhook.Token,
		Name	: webhookName,
		ChannelID: created_webhook.ChannelID,	
		ChannelName: msg.ChannelName,	
	}

	r.webhooks = append(r.webhooks, webhook)

	return &webhook, nil
}