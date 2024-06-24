package mdb

import (
	"database/sql"
	"fmt"
	"os"
	"time"
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
	return Mdb{db}, nil
}