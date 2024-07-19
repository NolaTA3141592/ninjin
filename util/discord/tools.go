package discord

import (
	"fmt"
	"ninjin/util/cls"
	"ninjin/util/mdb"
	"ninjin/util/slack"

	"github.com/bwmarrin/discordgo"
)

type Router struct {
	DISCORD_API_TOKEN	string
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

	r.webhooks, err = r.GetWebhookList()
	if err != nil {
		fmt.Println("Error getting webhook list : ", err)
	}

	return nil
}

func (r *Router) EventMassage(user *slack.User, msg *cls.Message, db *mdb.Mdb) {
	webhook, err := r.SelectWebhook(msg)
	if err != nil {
		fmt.Println("Error select Webhook : ", err)
		return
	}
	msg.DiscordChannelID = webhook.ChannelID
	r.Threading(msg, db)
	r.MessageSend(user, msg, webhook)
}

func (r *Router) Threading(msg *cls.Message, db *mdb.Mdb) {
	if !msg.ThreadMode {
		return
	}
	dpID, err := db.QueryThreadID(msg.Slack_parent_ID)	
	if err == nil || dpID != "" {
		// もうDB内にスレッドが作られている
		msg.Discord_thread_ID = dpID
	} else {
		dmID, err := db.QueryMessageID(msg.Slack_parent_ID)
		if(err == nil || dmID != "") {
			// メッセージがDB内に存在する
			// 今からスレッドを作る
			dpID, err = r.MakeThread(msg.DiscordChannelID, dmID)
			if err != nil {
				fmt.Println("Error Make Thread : ", err)
				msg.ThreadMode = false
				return
			}
			msg.Discord_thread_ID = dpID
			db.InsertThread(msg.Slack_parent_ID, msg.Discord_thread_ID)
		} else {
			// Threadを作成できない
			msg.ThreadMode = false
		}
	}
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