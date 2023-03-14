package main

import (
	"broker/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	authenticationServiceURL = "http://authentication-service/authenticate"
	logServiceURL            = "http://logger-service/log"
)

// HandleSubmission is the main point of entry into the broker. It accepts a JSON
// payload and performs an action based on the value of "action" in that JSON.
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload data.RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case data.Ping:
		app.ping(w)
	case data.Auth:
		app.authenticate(w, requestPayload.Auth)
	case data.Log:
		app.logItem(w, requestPayload.Log)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) ping(w http.ResponseWriter) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) logItem(w http.ResponseWriter, entry data.LogPayload) {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	request.Header.Set(string(data.HeaderContentType), string(data.ContentTypeJSON))

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	var responsePayload data.ResponsePayload
	if err := json.NewDecoder(response.Body).Decode(&responsePayload); err != nil {
		app.errorJSON(w, err)
		return
	}

	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, fmt.Errorf("status code %d: %s", response.StatusCode, responsePayload.Message))
		return
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a data.AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(http.MethodPost, authenticationServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer response.Body.Close()

	// create a variable we'll read response.Body into
	var responsePayload data.ResponsePayload

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&responsePayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if responsePayload.Error {
		app.errorJSON(w, fmt.Errorf("status code %d: %s", response.StatusCode, responsePayload.Message), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}
