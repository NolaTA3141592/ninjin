package slack

func GetUserName(userinfo *User) (string) {
	if userinfo.Profile.DisplayName != ""{
		return userinfo.Profile.DisplayName
	} else {
		return userinfo.RealName
	}
	return ""
}
