package main

import (
	"log"
)

var trackQueueSchema = `
CREATE TABLE IF NOT EXISTS track_queue (
	id SERIAL PRIMARY KEY,
	track_id	TEXT UNIQUE,
	name	TEXT,
	artist	TEXT,
	album	TEXT,
	album_art	TEXT,
	time	TEXT,
	queued_by	TEXT,
	queued_by_avatar	TEXT
);`

var userTableSchema = `
CREATE TABLE IF NOT EXISTS user_table (
	id	SERIAL PRIMARY KEY,
	user_id	TEXT UNIQUE,
	access_token	TEXT,
	avatar_url	TEXT,
	created_at	TEXT,
	updated_at	TEXT
);`

func initDB() error {
	log.Println("Initializing Database")

	context.db.MustExec(trackQueueSchema)
	context.db.MustExec(userTableSchema)

	return nil	
}
