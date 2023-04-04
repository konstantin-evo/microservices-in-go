package main

import (
	"broker/data"
	"broker/event"
	eventData "broker/event/data"
	"broker/logs"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/rpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	case data.LogRPC:
		app.logItemViaRPC(w, requestPayload.Log)
	case data.LogGRPC:
		app.logItemViaGRPC(w, requestPayload.Log)
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
func (app *Config) authenticate(w http.ResponseWriter, requestPayload data.AuthPayload) {
	responsePayload, err := callExternalService(app.AuthenticationServiceURL, requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if responsePayload.Error {
		app.errorJSON(w, fmt.Errorf("status code %d: %s", responsePayload.StatusCode, responsePayload.Message), http.StatusUnauthorized)
		return
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

func (app *Config) sendMail(w http.ResponseWriter, requestPayload data.MailPayload) {
	responsePayload, err := callExternalService(app.MailServiceURL, requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if responsePayload.Error {
		app.errorJSON(w, fmt.Errorf("status code %d: %s", responsePayload.StatusCode, responsePayload.Message), http.StatusUnauthorized)
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

func (app *Config) logItemViaGRPC(w http.ResponseWriter, requestPayload data.LogPayload) {
	gRPCaddress := fmt.Sprintf("%s:%s", app.LogServiceAddress, app.LogServiceGRPCPort)
	client, err := grpc.Dial(gRPCaddress, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer client.Close()

	gRPCPayload := data.RPCPayload(requestPayload)

	conn := logs.NewLogServiceClient(client)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	logResponse, err := conn.WriteLog(ctx, &logs.LogRequest{
		LogEntry: &logs.Log{
			Name: gRPCPayload.Name,
			Data: gRPCPayload.Data,
		},
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := data.ResponsePayload{
		Error:   false,
		Message: logResponse.Result,
		Data:    gRPCPayload,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}

func (app *Config) logItemViaRPC(w http.ResponseWriter, requestPayload data.LogPayload) {
	RPCaddress := fmt.Sprintf("%s:%s", app.LogServiceAddress, app.LogServiceRPCPort)
	client, err := rpc.Dial("tcp", RPCaddress)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rpcPayload := data.RPCPayload(requestPayload)

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

// util func to send post request with json payload
func callExternalService(url string, requestPayload interface{}) (data.ResponsePayload, error) {
	jsonData, err := json.MarshalIndent(requestPayload, "", "\t")
	if err != nil {
		return data.ResponsePayload{}, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return data.ResponsePayload{}, err
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return data.ResponsePayload{}, err
	}
	defer response.Body.Close()

	var responsePayload data.ResponsePayload
	err = json.NewDecoder(response.Body).Decode(&responsePayload)
	if err != nil {
		return data.ResponsePayload{}, err
	}

	responsePayload.StatusCode = response.StatusCode

	return responsePayload, nil
}
