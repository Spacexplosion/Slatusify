package main

import (
	"encoding/json"
	"fmt"
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

func setSlackStatus(status string) {
	// todo
}
