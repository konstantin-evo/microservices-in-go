package event

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func logEvent(entry Payload, logServiceURL string) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	request, err := http.NewRequest(http.MethodPost, logServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
