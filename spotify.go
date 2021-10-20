package main

import (
	"strings"

	"github.com/andybrewer/mack"
	"github.com/ghetzel/go-stockutil/sliceutil"
)

const appName = "Spotify"

type PlayStatus struct {
	playing bool
	state   string
	artist  string
	track   string
}

func isPlayerRunning() bool {
	procList, procListErr := mack.Tell("System Events", "name of processes")
	if procListErr != nil {
		panic(procListErr)
	}
	procNames := strings.Split(procList, ",")
	for i, n := range procNames {
		procNames[i] = strings.TrimLeft(n, " ")
	}
	return sliceutil.ContainsString(procNames, appName)
}

func getPlayStatus() PlayStatus {
	var status PlayStatus
	state, stateErr := mack.Tell(appName, "player state")
	if stateErr != nil {
		panic(stateErr)
	}
	status.state = state
	status.playing = state == "playing"

	if status.playing {
		var trackErr error
		status.artist, trackErr = mack.Tell(appName, "artist of current track")
		if trackErr != nil {
			panic(trackErr)
		}

		status.track, trackErr = mack.Tell(appName, "name of current track")
		if trackErr != nil {
			panic(trackErr)
		}
	}
	return status
}
