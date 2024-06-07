package mdb

import (
	_ "github.com/lib/pq"
)

func (db *Mdb) Insert(SlackMessageID string, DiscordMessageID string) (error) {
	sqlStatement := `
    INSERT INTO MessageDataBase (slackID, discordID)
    VALUES ($1, $2)`	
	_, err := db.Data.Exec(sqlStatement, SlackMessageID, DiscordMessageID)
	return err
}

func (db *Mdb) Query(slackID string) (string, error) {
    var discordID string
    sqlStatement := `SELECT discordID FROM slackID WHERE id=$1`
    row := db.Data.QueryRow(sqlStatement, slackID)
    err := row.Scan(&discordID)
    if err != nil {
        return "", err
    }
    return discordID, nil
}