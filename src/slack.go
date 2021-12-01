package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/skratchdot/open-golang/open"
)

type ApiCredentials struct {
	ClientId     string
	ClientSecret string
}

var slackClientCreds ApiCredentials = ApiCredentials{}

func delayKill(server *http.Server) {
	time.Sleep(time.Second)
	server.Close()
}

func doSlackRegistration() {
	http.HandleFunc("/oauth", func(resp http.ResponseWriter, req *http.Request) {
		fmt.Fprint(resp, slackRegistrationIndexContent)
	})
	http.HandleFunc("/getCreds", func(resp http.ResponseWriter, req *http.Request) {
		bodyData, err := json.Marshal(slackClientCreds)
		if err != nil {
			logger.Fatalln("Failed to marshal creds", err.Error())
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(resp, err)
		} else {
			fmt.Fprint(resp, string(bodyData))
		}
	})
	http.HandleFunc("/finish", func(resp http.ResponseWriter, req *http.Request) {
		// store token
		config.OauthToken = req.FormValue("token")
		config.store()

		resp.WriteHeader(http.StatusAccepted)
		resp.(http.Flusher).Flush()

		// quit web host
		context := req.Context()
		server := context.Value(http.ServerContextKey).(*http.Server)
		go delayKill(server)
	})
	go func() {
		logger.Println("HTTP serve terminated", http.ListenAndServe(":51401", nil))
		menuItemSlackRegister.Enable()
	}()
	open.Start(fmt.Sprintf("https://slack.com/oauth/authorize?client_id=%s&scope=client&redirect_uri=http%%3A%%2F%%2Flocalhost%%3A51401%%2Foauth", slackClientCreds.ClientId))
}

const statusEndpoint = "https://slack.com/api/users.profile.set"

type ProfileEntry struct {
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}

type profileSetResult struct {
	Ok      bool         `json:"ok"`
	Profile ProfileEntry `json:"profile"`
	Error   string       `json:"error"`
}

func setSlackStatus(status string) (bool, error) {
	var emoji string
	if status != "" { // empty means unsetting status
		emoji = config.StatusEmoji
	}

	body := make(map[string]ProfileEntry)
	body["profile"] = ProfileEntry{
		StatusText:  status,
		StatusEmoji: emoji,
	}
	bodyData, _ := json.Marshal(body)

	request, reqError := http.NewRequest(http.MethodPost, statusEndpoint, bytes.NewBuffer(bodyData))
	if reqError == nil {
		request.Header.Add("Authorization", "Bearer "+config.OauthToken)
		request.Header.Add("Content-Type", "application/json; charset=UTF-8")
		response, httpErr := http.DefaultClient.Do(request)
		if httpErr == nil {
			var result profileSetResult
			logger.Println("Status set -", response.Status)
			if response.StatusCode == http.StatusOK {
				resultData, buffErr := ioutil.ReadAll(response.Body)
				if buffErr != nil {
					return false, buffErr
				}
				if parseErr := json.Unmarshal(resultData, &result); parseErr != nil {
					return false, parseErr
				}
				if !result.Ok {
					return false, fmt.Errorf("Slack API Error: %s", result.Error)
				}
				return true, nil
			}
			// else
			return false, fmt.Errorf("HTTP status %s", response.Status)
		}
		// else
		return false, httpErr
	}
	// else
	logger.Fatalln("Failed to build Slack request", reqError)
	return false, reqError
}
