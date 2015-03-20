package main

import (
	"encoding/json"
	"log"
)

//NotificationMessage simple struct for sending data
type NotificationMessage struct {
	Event string `json:"event"`
	Value string `json:"values"`
	Track string `json:"track"`
}

func processQueue(ch chan interface{}) {
	for m := range ch {
		messagebody := map[string]interface{}{}
		err := json.Unmarshal(m.([]byte), &messagebody)
		if err != nil {
			log.Panic("the unmarshall plan")
		}
		switch messagebody["command"] {
		case "play_track":
			log.Println("Play Track:", messagebody)
			if str, ok := messagebody["param"].(string); ok {
				setNextTrack("spotify:track:" + str)
			} else {
				log.Panic("was unable to set current track")
			}
			// case "skip_track"
			// case "other command"
		} //end of switch

	}
}
