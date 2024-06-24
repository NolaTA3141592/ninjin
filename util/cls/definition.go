package cls

type Message struct {
	Content		string
	Slack_ID	string
	Discord_ID	string
	ChannelName	string
	FileURLs	[]string
	FileNames	[]string
	SlackChannelID	string
	DiscordChannelID	string
	ThreadMode	bool
	Slack_parent_ID	string
	Discord_thread_ID	string
}