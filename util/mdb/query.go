package mdb

import (
	"fmt"
	"ninjin/util/cls"

	_ "github.com/lib/pq"
)

func (db *Mdb) Insert(msg *cls.Message) (error) {
	sqlStatement := `
    INSERT INTO MessageDatabase (slackID, discordID, ChannelName)
    VALUES ($1, $2, $3)`	
	_, err := db.Data.Exec(sqlStatement, msg.Slack_ID, msg.Discord_ID, msg.ChannelName)
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