package mdb

import (
	"fmt"
	"ninjin/util/cls"

	_ "github.com/lib/pq"
)

func (db *Mdb) Insert(msg *cls.Message) (error) {
	sqlStatement := `
    INSERT INTO MessageDatabase (slackID, discordID, ChannelName, slackChannelID, discordChannelID)
    VALUES ($1, $2, $3, $4, $5)`	
	_, err := db.Data.Exec(sqlStatement, msg.Slack_ID, msg.Discord_ID, msg.ChannelName, msg.SlackChannelID, msg.DiscordChannelID)
	if err != nil {
		fmt.Println("DB Insert error : ", err)
	}
	return err
}

func (db *Mdb) InsertThread(slackThreadID string, discordthreadID string) (error) {
	sqlStatement := `
    INSERT INTO ThreadDatabase (slackThreadID, discordThreadID)
    VALUES ($1, $2)`	
	_, err := db.Data.Exec(sqlStatement, slackThreadID, discordthreadID)
	if err != nil {
		fmt.Println("DB Insert error : ", err)
	}
	return err
}

func (db *Mdb) QueryMessageID(slackID string) (string, error) {
    var discordID string
    sqlStatement := `SELECT discordID FROM MessageDatabase WHERE slackID = $1`
    row := db.Data.QueryRow(sqlStatement, slackID)
    err := row.Scan(&discordID)
    if err != nil {
        return "", err
    }
    return discordID, nil
}

func (db *Mdb) QueryChannelName(slackID string) (string, error) {
    var ChannelName string
    sqlStatement := `SELECT ChannelName FROM MessageDatabase WHERE slackID = $1`
    row := db.Data.QueryRow(sqlStatement, slackID)
    err := row.Scan(&ChannelName)
    if err != nil {
        return "", err
    }
    return ChannelName, nil
}

func (db *Mdb) QueryThreadID(slackThreadID string) (string, error) {
    var discordThreadID string
    sqlStatement := `SELECT discordThreadID FROM ThreadDatabase WHERE slackThreadID = $1`
    row := db.Data.QueryRow(sqlStatement, slackThreadID)
    err := row.Scan(&discordThreadID)
    if err != nil {
        return "", err
    }
    return discordThreadID, nil
}