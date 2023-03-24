package main

import (
	"broker/data"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// readJSON tries to read the body of a request and converts it into JSON
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, requestPayload any) error {
	maxBytes := 1048576 // one megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	// Use a temporary struct to decode the request payload
	var temp struct {
		Action string           `json:"action"`
		Auth   data.AuthPayload `json:"auth,omitempty"`
		Log    data.LogPayload  `json:"log,omitempty"`
		Mail   data.MailPayload `json:"mail,omitempty"`
	}

	err := dec.Decode(&temp)
	if err != nil {
		return err
	}

	// Convert the string action to an enum value
	var action data.ActionType
	switch temp.Action {
	case "ping":
		action = data.Ping
	case "auth":
		action = data.Auth
	case "log":
		action = data.Log
	case "logGrpc":
		action = data.LogGRPC
	case "mail":
		action = data.Mail
	default:
		return errors.New("unknown action")
	}

	// Copy the decoded values to the output parameter
	switch payload := requestPayload.(type) {
	case *data.RequestPayload:
		payload.Action = action
		payload.Auth = temp.Auth
		payload.Log = temp.Log
		payload.Mail = temp.Mail
	default:
		return errors.New("unsupported payload type")
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON takes a response status code and arbitrary data and writes a json response to the client
func (app *Config) writeJSON(w http.ResponseWriter, status int, payload any, headers ...http.Header) error {
	out, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set(string(data.HeaderContentType), string(data.ContentTypeJSON))
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON takes an error, and optionally a response status code, and generates and sends
// a json error response
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload data.ResponsePayload
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJSON(w, statusCode, payload)
}
