// Groupify API
//
// TeamOFP - GopherGala 2015
//

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/oauth2"

)

// App context
type Context struct {
	db *sqlx.DB
	//airbrake *gobrake.Notifier
	tq    *TrackQueue
	np    *nowPlaying
	oauth *oauth2.Config
	messages *chan []byte
	notifications *chan []byte
}

var context = &Context{}
var store = sessions.NewCookieStore([]byte("Groupify.go FTW!"))

// GetInfo - Info Endpoint. Returns versioning info.
func GetInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Spotify Groupify API v0.2.0")
}

func main() {
	log.Println("Starting Groupify API...")

	// Load .env
	err := godotenv.Load()
	if err != nil {
		// Can't load .env, so setenv defaults
		os.Setenv("SQL_HOST", "localhost")
		os.Setenv("SQL_USER", "root")
		os.Setenv("SQL_PASSWORD", "")
		os.Setenv("SQL_DB", "spotify_remote")
	}

	// Setup App Context
	// Setup DB
	db, err := sqlx.Open("sqlite3", "./spotify-remote.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	context.db = db

	// Sending Channel
	messages := make(chan []byte)
	context.messages = &messages

	// Receiving channel
	notifications := make(chan []byte)
	context.notifications = &notifications

	// Track Queue
	tq := &TrackQueue{}
	context.tq = tq

	// Now Playing
	context.np = &nowPlaying{}

	// Start Notification processor
	go notificationProcessor(notifications)

	router := mux.NewRouter()
	r := router.PathPrefix("/api/v1").Subrouter() // Prepend API Version

	// Setup Negroni
	n := negroni.Classic()

	// Info
	r.HandleFunc("/", GetInfo).Methods("GET")

	// TrackQueue
	r.HandleFunc("/queue/add", PostAddTrack).Methods("POST")
	r.HandleFunc("/queue/list", GetListTracks).Methods("GET")
	r.HandleFunc("/queue/delete", PostDeleteTrack).Methods("POST")
	r.HandleFunc("/queue/next", PostSkipTrack).Methods("POST")
	//r.HandleFunc("/queue/upvote", AddTrack).Methods("POST")
	//r.HandleFunc("/queue/downvote", AddTrack).Methods("POST")

	r.HandleFunc("/auth", Auth).Methods("GET")
	r.HandleFunc("/callback", Callback).Methods("GET")

	// Web Socket Handler
	r.HandleFunc("/remote", WSRemoteHandler).Methods("GET")

	// Setup router
	n.UseHandler(r)

	// Start Serve
	if os.Getenv("PORT") != "" {
		n.Run(strings.Join([]string{":", os.Getenv("PORT")}, ""))
	} else {
		n.Run(":8080")
	}

}
