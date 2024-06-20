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
	if files, ok := data["files"].([]interface{}); ok {
		for _, file := range files {
			if fileMap, ok := file.(map[string]interface{}); ok {
				if thumbs, ok := fileMap["thumb_360"]; ok {
					msg.FileURLs = append(msg.FileURLs, thumbs.(string))
				}
				if other, ok := fileMap["url_private"]; ok {
					msg.FileURLs = append(msg.FileURLs, other.(string))
				}
			}
		}
	}
	if files, ok := data["files"].([]interface{}); ok {
		for _, file := range files {
			if fileMap, ok := file.(map[string]interface{}); ok {
				if name, ok := fileMap["name"]; ok {
					msg.FileNames = append(msg.FileNames, name.(string))
				}
			}
		}
	}

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