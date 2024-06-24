package discord

import (
	"fmt"

	_ "github.com/lib/pq"
)

func (r *Router) MakeThread(DiscordChannelID string, DiscordMessageID string) (string, error) {
	ch, err := r.Bot.MessageThreadStart(DiscordChannelID, DiscordMessageID, "", 4320)
	if err != nil {
		fmt.Println("Error MakeThread ", err)
		return "", err
	}
	return ch.ID, nil
}