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

	return nil
}

func (r *Router) Execute(user *slack.User, msg string) {
	r.Bot.ChannelMessageSend(r.TEST_CHANNEL_ID, fmt.Sprintf("[%s]: %s", user.RealName, msg))
}