package main

import (
	"github.com/andybrewer/mack"
)

type PlayStatus struct {
	playing bool
	state   string
	artist  string
	track   string
}

func getPlayStatus() PlayStatus {
	var status PlayStatus
	stateResponse, stateErr := mack.Tell("Spotify", "player state")
	if stateErr != nil {
		panic(stateErr)
	}
	status.state = stateResponse
	status.playing = stateResponse == "playing"

	if status.playing {
		trackResponse, trackErr := mack.Tell("Spotify", "artist of current track")
		if trackErr != nil {
			panic(trackErr)
		}
		status.artist = trackResponse

		trackResponse, trackErr = mack.Tell("Spotify", "name of current track")
		if trackErr != nil {
			panic(trackErr)
		}
		status.track = trackResponse
	}
	return status
}
