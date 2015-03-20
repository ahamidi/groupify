package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendWSMessage(t *testing.T) {

	// Create connection
	c, err := NewWSConnection("ws://localhost:8080/api/v1/remote", "http://localhost/")

	// Send message
	err = SendWSMessage(c, "hello there")

	assert.Nil(t, err)
}

func TestSendTrackEnd(t *testing.T) {
	// Create notification
	type NotificationMessage struct {
		Event string `json:"event"`
		Value string `json:"values"`
		Track string `json:"track,omitempty"`
	}

	trackEnd := NotificationMessage{
		Event: "track_end",
		Value: "some random track",
	}

	// Create connection
	c, err := NewWSConnection("ws://localhost:8080/api/v1/remote", "http://localhost/")

	// Send message
	err = SendWSMessage(c, trackEnd)

	assert.Nil(t, err)

}
