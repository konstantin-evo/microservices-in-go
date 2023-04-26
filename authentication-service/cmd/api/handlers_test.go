package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (fn RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req), nil
}

func NewTestClient(roundTripFunc RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: roundTripFunc,
	}
}

func Test_Authenticate(t *testing.T) {
	prepareTestApp()

	body := prepareRequestBody()

	req, _ := http.NewRequest(http.MethodPost, "/authenticate", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(testApp.Authenticate)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Errorf("expected http.StatusAccepted but got %d", rr.Code)
	}
}

func prepareTestApp() {
	jsonToReturn := prepareMockResponse()

	client := NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusAccepted,
			Body:       io.NopCloser(bytes.NewBufferString(jsonToReturn)),
			Header:     make(http.Header),
		}
	})

	testApp.Client = client
}

func prepareMockResponse() string {
	return `
	{
		"error": false,
		"message": "test message"
	}
	`
}

func prepareRequestBody() []byte {
	postBody := map[string]interface{}{
		"email":    "me@here.com",
		"password": "verysecret",
	}

	body, _ := json.Marshal(postBody)
	return body
}
