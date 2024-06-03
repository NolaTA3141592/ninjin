package slack

import (
	slackgo "github.com/slack-go/slack"
)

func GetIcon(userinfo *slackgo.User) string {
	if userinfo.Profile.Image512 != ""{
		return userinfo.Profile.Image512
	}
	return ""
}
