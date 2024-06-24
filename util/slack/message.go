package slack

import (
	"fmt"
	"ninjin/util/cls"
	"regexp"

	slackgo "github.com/nlopes/slack"
	slackgo2 "github.com/slack-go/slack"
)

type Mentions struct {
	UserID_row 		[]string
	UserID_list   	[]string
	UserName 		[]string
}

func (sl SlackUtil) AttachMessageInfo(msg *cls.Message, data map[string]interface{}) error {
	msg.Content = data["text"].(string)
	msg.Slack_ID = data["ts"].(string)
	msg.ChannelName = sl.GetChannelNameByID(data["channel"].(string))
	msg.FileURLs = sl.GetFileURLs(msg, data)
	msg.FileNames = sl.GetFileNames(msg, data)
	msg.Content = sl.ReplaceMentions(msg)
	msg.SlackChannelID = data["channel"].(string)
	msg.ThreadMode = sl.ContainsThreadTS(data)

	if msg.ThreadMode {
		msg.Slack_parent_ID = data["thread_ts"].(string)
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

func (sl SlackUtil) GetFileURLs(msg *cls.Message, data map[string]interface{}) []string {
	FileURLs := []string{}
	if files, ok := data["files"].([]interface{}); ok {
		for _, file := range files {
			if fileMap, ok := file.(map[string]interface{}); ok {
				if thumbs, ok := fileMap["thumb_360"]; ok {
					FileURLs = append(FileURLs, thumbs.(string))
				}
				if other, ok := fileMap["url_private"]; ok {
					FileURLs = append(FileURLs, other.(string))
				}
			}
		}
	}
	return FileURLs
}

func (sl SlackUtil) GetFileNames(msg *cls.Message, data map[string]interface{}) []string {
	FileNames := []string{}
	if files, ok := data["files"].([]interface{}); ok {
		for _, file := range files {
			if fileMap, ok := file.(map[string]interface{}); ok {
				if name, ok := fileMap["name"]; ok {
					FileNames = append(FileNames, name.(string))
				}
			}
		}
	}
	return FileNames
}

func (sl SlackUtil) ReplaceMentions(msg *cls.Message) string {
	message := msg.Content
	re := regexp.MustCompile(`<@[A-Z0-9]{11}>`)
	matches := re.FindAllString(message, -1)

	mentions := Mentions{
		UserID_row:   	[]string{},
		UserID_list:   	[]string{},
		UserName:		[]string{},
	}

	mentions.UserID_row = append(mentions.UserID_row, matches...)
	re = regexp.MustCompile(`<@(.*?)>`)

	for _, match := range mentions.UserID_row {
		userID := re.FindStringSubmatch(match)[1]
		mentions.UserID_list = append(mentions.UserID_list, userID)
	}
	for _, userID := range mentions.UserID_list {
		api := slackgo2.New(sl.SLACK_API_TOKEN)
		userinfo, err := api.GetUserInfo(userID)
		if err != nil {
			return ""
		}
		usernames := GetUserName(userinfo)
		mentions.UserName = append(mentions.UserName, usernames)
	}

	for i := 0; i < len(mentions.UserID_list); i++ {
		re := regexp.MustCompile(fmt.Sprintf("@%s", mentions.UserID_list[i]))
		replacement := fmt.Sprintf("@%s", mentions.UserName[i])
		message = re.ReplaceAllString(message, replacement)
	}
	return message
}

func (sl SlackUtil) ContainsThreadTS(data map[string]interface{}) bool {
	_, exist := data["thread_ts"]
	return exist
}