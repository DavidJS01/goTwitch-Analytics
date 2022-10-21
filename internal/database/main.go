package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateMessageTable() {
	db := connectToDB()
	sqlStatement := `
	create schema "postgres"."twitch" if not exists;
	create table "postgres"."twitch"."messages" if not exists (
		username varchar,
		twitch_channel varchar,
		message varchar,
		message_timestamp timestamp without time zone default (now() at time zone 'utc')
	`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	db.Close()

}

func InsertTwitchMesasge(username string, message string, channel string) {
	db := connectToDB()
	sqlStatement := `
				INSERT INTO "postgres"."twitch"."messages" (username, twitch_channel, message)
				VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, username, channel, message)
	db.Close()
	if err != nil {
		panic(err)
	}
}
