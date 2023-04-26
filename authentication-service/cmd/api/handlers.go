package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type authenticationRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authenticationResponsePayload struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload authenticationRequestPayload
	if err := app.readJSON(w, r, &requestPayload); err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate the user against the database
	user, err := app.Repo.GetByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("user does not exist"), http.StatusUnauthorized)
		return
	}

	if valid, err := app.Repo.PasswordMatches(requestPayload.Password, *user); err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	// Log authentication
	if err := app.logRequest("authentication", fmt.Sprintf("%s logged in", user.Email)); err != nil {
		app.errorJSON(w, err)
		return
	}

	responsePayload := authenticationResponsePayload{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	app.writeJSON(w, http.StatusAccepted, responsePayload)
}

type logEntry struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) logRequest(name, data string) error {
	entry := logEntry{
		Name: name,
		Data: data,
	}

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		return err
	}

	logServiceURL := "http://logger-service/log"
	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	if _, err := app.Client.Do(request); err != nil {
		return err
	}

	return nil
}
