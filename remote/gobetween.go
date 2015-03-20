package main

import (
	"log"
	//"os"
	"flag"
	"os/exec"
	"strings"

	//"github.com/crowdmob/goamz/sqs"
	"github.com/joho/godotenv"
)

// start of script tag
const ScriptStart = "tell application \"Spotify\" to "

var commands = map[string]string{
	"state":      "player state",
	"play":       "play",
	"pause":      "pause",
	"duration":   "duration of current track",
	"name":       "name of current track",
	"album":      "album of current track",
	"id":         "id of current track",
	"artwork":    "artwork of current track",
	"vol_loud":   "set sound volume to 100",
	"vol_soft":   "set sound volume to 20",
	"vol_norm":   "set sound volume to 50",
	"set_volume": "set sound volume to ", //requires parameter
	"play_track": "play track ",          //requires parameter
	"position":   "player position",
}

// API Host Flag
var apiHost = flag.String("host", "localhost:8080",
	"Hostname of Groupify API")

func callSpotify(command string, param string) string {
	fullcmd := ScriptStart + commands[command] + param

	out, err := exec.Command("/usr/bin/osascript", "-e", fullcmd).Output()
	if err != nil {
		log.Fatal(err)
		log.Fatal(out)
	}
	return strings.TrimSpace(string(out))
}

func main() {
	log.Println("Starting Groupify Remote")

	flag.Parse()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Create host string
	url := strings.Join([]string{"ws://", *apiHost, "/api/v1/remote"}, "")

	// Create WS Connection
	c, err := NewWSConnection(url, "http://localhost/")

	// Spotify status queue
	spotifyState := make(chan interface{})
	go polling(spotifyState)

	// Create notification message channel
	notifications := make(chan interface{})
	go WSMessageReceiver(c, notifications)
	go processQueue(notifications)

	for s := range spotifyState {
		log.Println("Spotify State:", s)
		SendWSMessage(c, s)
	}

	done := make(chan bool)
	<-done
}
