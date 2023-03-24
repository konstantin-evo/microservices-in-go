package main

import (
	"broker/data"
	"broker/event"
	eventData "broker/event/data"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
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
		app.logEvent(w, requestPayload.Log)
	case data.LogGRPC:
		app.logItemViaRPC(w, requestPayload.Log)
	case data.Mail:
		app.sendMail(w, requestPayload.Mail)
	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) ping(w http.ResponseWriter) {
	payload := data.ResponsePayload{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

// authenticate calls the authentication microservice and sends back the appropriate response
func (app *Config) authenticate(w http.ResponseWriter, a data.AuthPayload) {
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(http.MethodPost, app.AuthenticationServiceURL, bytes.NewBuffer(jsonData))
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

func (app *Config) sendMail(w http.ResponseWriter, msg data.MailPayload) {
	jsonData, err := json.MarshalIndent(msg, "", "\t")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// post to mail service
	request, err := http.NewRequest(http.MethodPost, app.MailServiceURL, bytes.NewBuffer(jsonData))
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

	// make sure we get back the right status code
	if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling mail service"))
		return
	}

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

// logEvent logs an event using the logger-service. It makes the call by pushing the data to RabbitMQ.
func (app *Config) logEvent(w http.ResponseWriter, logPayload data.LogPayload) {
	err := app.pushToQueue(logPayload.Name, logPayload.Data)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var responsePayload data.ResponsePayload
	responsePayload.Error = false
	responsePayload.Message = "The event info is sent to the queue."
	responsePayload.Data = logPayload

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

// pushToQueue pushes a message into RabbitMQ
func (app *Config) pushToQueue(name, msg string) error {
	emitter, err := event.NewEventEmitter(app.Rabbit)
	if err != nil {
		return err
	}

	payload := data.LogPayload{
		Name: name,
		Data: msg,
	}

	eventPayload, _ := json.MarshalIndent(&payload, "", "\t")
	err = emitter.Push(string(eventPayload), string(eventData.SeverityLog))
	if err != nil {
		return err
	}
	return nil
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, l data.LogPayload) {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := data.RPCPayload{
		Name: l.Name,
		Data: l.Data,
	}

	var result string
	err = client.Call("RPCServer.LogInfo", rpcPayload, &result)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := data.ResponsePayload{
		Error:   false,
		Message: result,
		Data:    rpcPayload,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
