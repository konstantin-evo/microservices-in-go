package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"broker/event/data"
)

type Logger struct {
	logServiceURL string
	client        *http.Client
}

func NewLogger(logServiceURL string) *Logger {
	return &Logger{
		logServiceURL: logServiceURL,
		client:        &http.Client{},
	}
}

func (logger *Logger) logEvent(entry data.Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest(http.MethodPost, logger.logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := logger.client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

	return nil
}
