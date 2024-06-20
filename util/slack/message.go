package slack

import (
	"fmt"
	"ninjin/util/cls"
	"regexp"

	slackgo "github.com/nlopes/slack"
	slackgo2 "github.com/slack-go/slack"
)

type Mentions struct {
	UserID   		[]string
	UserName 		[]string
}

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

	re := regexp.MustCompile(`<@[A-Z0-9]{11}>`)

	matches := re.FindAllString(msg.Content, -1)

	mentions := Mentions{
		UserID:   []string{},
		UserName: []string{},
	}

	mentions.UserID = append(mentions.UserID, matches...)

	for _, userID := range mentions.UserID {
		api := slackgo2.New(sl.SLACK_API_TOKEN)
		userinfo, err := api.GetUserInfo(userID)
		if err != nil {
			fmt.Printf("Error fetching user info for %s: %v\n", userID, err)
			continue
		}
		usernames := GetUserName(userinfo)
		mentions.UserName = append(mentions.UserName, usernames)
	}

	fmt.Println(mentions)


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