package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Streamer struct {
	Name      string
	Is_Active bool
}

const (
	host     = "172.17.0.1"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "postgres"
)

type InsertMessage func(username string, message string, channel string)

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


func newNullString(pid int) sql.NullInt64 {
	if pid == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: int64(pid),
		Valid: true,
	}
}

func connectToDB() *sql.DB {
	psqlInfo := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		user,
		password,
		host,
		port,
		dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	return db
}

func CreateMessageTable() {
	db := connectToDB()
	sqlStatement := `
	create schema if not exists twitch;
	create table if not exists twitch.streamers  (
		twitch_channel varchar,
		time_added timestamp without time zone default (now() at time zone 'utc'));
	`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	db.Close()

}



func CreateStreamerTable() {
	db := connectToDB()
	sqlStatement := `
	create schema if not exists twitch;
	create table if not exists twitch.streamers  (
		twitch_channel varchar,
		is_active boolean,
		time_added timestamp without time zone default (now() at time zone 'utc')
		);
		
	`
	_, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	db.Close()
}

func InsertStreamer(channel string, is_active bool) error {
	db := connectToDB()
	sqlStatement := `
				INSERT INTO "twitch"."streamers" (twitch_channel, is_active)
				VALUES ($1, $2);`
	_, err := db.Exec(sqlStatement, channel, is_active)
	db.Close()
	if err != nil {
		panic(err)
	}
	return err
}

func GetStreamerData() ([]Streamer, error) {
	streamers := []Streamer{}
	db := connectToDB()
	sqlStatement := `
		SELECT  "twitch_channel", "is_active" FROM (
		SELECT *, rank() OVER (PARTITION BY "twitch_channel" ORDER BY "time_added" DESC) rank_number from twitch.streamers
		) AS t
		WHERE rank_number = 1
		`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	var streamer Streamer
	for rows.Next() {
		err = rows.Scan(&streamer.Name, &streamer.Is_Active)
		streamers = append(streamers, streamer)

	}
	db.Close()
	if err != nil {
		panic(err)
	}
	return streamers, err

}

func CreateStreamEventsTable() {
	db := connectToDB()
	sqlStatement := `
	create table if not exists twitch.stream_events  (
		twitch_channel varchar,
		listening bool, 
		pid varchar,
		time_added timestamp without time zone default (now() at time zone 'utc')
		);`
	_, err := db.Exec(sqlStatement)
	db.Close()
	if err != nil {
		panic(err)
	}
}

func GetLatestPID(streamer string) int {
	db := connectToDB()
	sqlStatement := `
	SELECT  "pid" FROM (
		SELECT *, rank() OVER (PARTITION BY "twitch_channel" ORDER BY "time_added" desc ) rank_number from twitch.stream_events 
		WHERE "twitch_channel" = $1
		) AS t
		WHERE rank_number = 1
	`
	var pid int
	rows := db.QueryRow(sqlStatement, streamer)
	err := rows.Scan(&pid)
	if err != nil {
		log.Print("L")
	}

	db.Close()
	if err != nil {
		panic(err)
	}
	return pid
}

func InsertStreamEvent(twitchChannel string, listening bool, pid int) {
	db := connectToDB()
	sqlStatement := `
	INSERT INTO twitch.stream_events (twitch_channel, listening, pid)
	VALUES ($1, $2, $3);`
	_, err := db.Exec(sqlStatement, twitchChannel, listening, newNullString(pid))
	db.Close()
	if err != nil {
		panic(err)
	}
}
