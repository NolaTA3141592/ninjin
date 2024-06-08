package slack

import (
	"fmt"
	"ninjin/util/cls"

	slackgo "github.com/nlopes/slack"
)

func (sl SlackUtil) AttachMessageInfo(msg *cls.Message, data map[string]interface{}) error {
	msg.Content = data["text"].(string)
	msg.Slack_ID = data["ts"].(string)
	msg.ChannelName = sl.GetChannelNameByID(data["channel"].(string))
	return nil
}

func (sl SlackUtil) GetChannelNameByID(channelID string) string {
	api := slackgo.New(sl.SLACK_API_TOKEN)
	channelinfo, err := api.GetConversationInfo(channelID, false)
	if err != nil {
		fmt.Println("error GetChannelInfo : ", err)
		return ""
	}
	return channelinfo.Name
}