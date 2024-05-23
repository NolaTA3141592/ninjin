package slack

import (
	slackgo "github.com/slack-go/slack"
)

func GetUserName(userinfo *slackgo.User) (string) {
	if userinfo.Profile.DisplayName != ""{
		return userinfo.Profile.DisplayName
	} else {
		return userinfo.RealName
	}
	return ""
}
