package main

import (
	"encoding/json"
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
var slackRegistrationIndexContent string

var menuItemSlackRegister *systray.MenuItem

func loadAssets() {
	var ioErr error
	iconData, ioErr = os.ReadFile("../Resources/icon.png")
	if ioErr != nil {
		panic(ioErr)
	}

	var slackRegistrationIndexData []byte
	slackRegistrationIndexData, ioErr = os.ReadFile("../Resources/index.html")
	if ioErr != nil {
		panic(ioErr)
	}
	slackRegistrationIndexContent = string(slackRegistrationIndexData)

	var clientCredsData []byte
	clientCredsData, ioErr = os.ReadFile("../Resources/clientCreds.json")
	jsonErr := json.Unmarshal(clientCredsData, &slackClientCreds)
	if jsonErr != nil {
		logger.Fatalln(jsonErr)
	}
}

// Set up systray menu
func onReady() {
	systray.SetIcon(iconData)
	systray.SetTooltip(defaultTooltip)

	menuItemSlackRegister = systray.AddMenuItem("Authorize Slack API", "Get a token to allow editing status on Slack")
	go func() { // mRegister handler
		_ = <-menuItemSlackRegister.ClickedCh
		menuItemSlackRegister.Disable()
		doSlackRegistration()
	}()

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

	var nextStatus string
	if status.playing {
		nextStatus = infoStr
		systray.SetTooltip(infoStr)
	} else {
		systray.SetTooltip(status.state)
	}

	setOk, setErr := setSlackStatus(nextStatus)
	if setOk {
		currentSpotifyStatus = status
	} else {
		logger.Println("Set failed", setErr)
	}
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
