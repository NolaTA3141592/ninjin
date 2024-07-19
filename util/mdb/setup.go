package mdb

import (
	"database/sql"
	"fmt"
	"os"
	"time"
	_ "github.com/lib/pq"
)

type Mdb struct {
	Data *sql.DB
}

func Setup() (Mdb, error) {
	dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)
	
	time.Sleep(10 * time.Second)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        fmt.Println("error Open DataBase", err)
		return Mdb{db}, err
    }
	
	err = db.Ping()
	if err != nil {
        fmt.Println("error", err)
		return Mdb{db}, err
	}

	fmt.Println("Successfully connecting DataBase!")

	createTable := `
	CREATE TABLE IF NOT EXISTS MessageDatabase (
		slackID TEXT PRIMARY KEY,
		discordID TEXT,
		ChannelName TEXT,
		slackChannelID TEXT,
		discordChannelID TEXT
	);

	CREATE TABLE IF NOT EXISTS ThreadDatabase (
		slackThreadID TEXT PRIMARY KEY,
		discordThreadID TEXT,
		ChannelName TEXT,
		slackChannelID TEXT,
		discordChannelID TEXT
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		fmt.Println("error create table", err)
	}
	return Mdb{db}, nil
}