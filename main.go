package main

import (
	"fmt"
	"time"

	"github.com/getlantern/systray"
)

const defaultTitle = "Slatusify"

var pollTicker time.Ticker
var currentStatus PlayStatus

// Set up systray menu
func onReady() {
	// systray.SetIcon(icon.Data)
	systray.SetTitle(defaultTitle)
	systray.SetTooltip("Set Slack status to what's playing on Spotify")
	mQuit := systray.AddMenuItem("Quit", "")
	go func() { // mQuit handler
		_ = <-mQuit.ClickedCh
		exit()
	}()
	mQuit.Enable()
}

func setNewStatus(status PlayStatus) {
	var infoStr string
	if status.playing {
		infoStr = fmt.Sprintf("%s - %s", status.artist, status.track)
		systray.SetTitle(infoStr)
	} else {
		systray.SetTitle(defaultTitle)
	}
	currentStatus = status
}

func runStatusPolling() {
	pollTicker = *time.NewTicker(1000000000)
	for tick := range pollTicker.C {
		status := getPlayStatus()
		if currentStatus != status {
			fmt.Printf("%s Player is %s: %s - %s\n", tick.Format(time.RFC3339), status.state, status.artist, status.track)
			setNewStatus(status)
		}
	}
}

func main() {
	go runStatusPolling()
	systray.Run(onReady, nil)
}

func exit() {
	fmt.Println(time.Now().Format(time.RFC3339), "Exiting...")
	pollTicker.Stop()
	systray.Quit()
}
