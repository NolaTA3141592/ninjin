package slack

func (sl SlackUtil) GetUserName(userinfo *User) error {
	if userinfo.Profile.DisplayName != nil{
		return userinfo.Profile.DisplayName
	} else {
		return userinfo.RealName
	}
	return nil
}
