package discord

import (
	"fmt"
	"ninjin/util/slack"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	DISCORD_API_TOKEN	string
	TEST_CHANNEL_ID		string
	Bot  				*discordgo.Session
	SERVER_ID 			string
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

	r.SERVER_ID = ""
	webhooks, err := r.GetWebhookList()
	fmt.Println(webhooks)
	if err != nil {
		fmt.Println("Error getting webhook list : ", err)
	}

	return nil
}

func (r *Router) EventMassage(user *slack.User, msg string) {
	r.Bot.ChannelMessageSend(r.TEST_CHANNEL_ID, fmt.Sprintf("[%s]: %s", user.RealName, msg))
}