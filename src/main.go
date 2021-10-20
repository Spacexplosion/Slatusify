package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/getlantern/systray"
)

const defaultTooltip = "Slatusify"

var pollTicker time.Ticker
var isSpotifyRunning bool
var currentSpotifyStatus PlayStatus

var logger *log.Logger = log.Default()
var iconData []byte

func loadAssets() {
	var ioErr error
	iconData, ioErr = os.ReadFile("../Resources/icon.png")
	if ioErr != nil {
		panic(ioErr)
	}
}

// Set up systray menu
func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip(defaultTooltip)
	mQuit := systray.AddMenuItem("Quit", "")
	go func() { // mQuit handler
		_ = <-mQuit.ClickedCh
		exit()
	}()
	mQuit.Enable()
}

// Set state on changed player status
func setNewStatus(status PlayStatus) {
	infoStr := fmt.Sprintf("%s - %s", status.artist, status.track)
	logger.Printf("Player is %s: %s\n", status.state, infoStr)
	if status.playing {
		systray.SetTooltip(infoStr)
	} else {
		systray.SetTooltip(status.state)
	}
	currentSpotifyStatus = status
}

// Poll for state changes
func runStatusPolling() {
	pollTicker = *time.NewTicker(1000000000)

	defer func() {
		err := recover()
		if err != nil {
			pollTicker.Stop()
			isSpotifyRunning = false
			logger.Println(err)
			logger.Println("Restarting polling...")
			go runStatusPolling()
		}
	}()

	for range pollTicker.C {
		if isSpotifyRunning {
			status := getPlayStatus()
			if currentSpotifyStatus != status {
				setNewStatus(status)
			}
		} else {
			isSpotifyRunning = isPlayerRunning()
		}
	}
}

func main() {
	loadAssets()
	go runStatusPolling()
	systray.Run(onReady, nil)
}

func exit() {
	logger.Println("Exiting...")
	pollTicker.Stop()
	systray.Quit()
}
