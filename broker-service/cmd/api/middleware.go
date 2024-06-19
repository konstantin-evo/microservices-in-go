package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Printf("Incoming request: %s %s %s\nBody: %s", r.Method, r.RequestURI, r.RemoteAddr, string(bodyBytes))

		ww := &responseWriter{ResponseWriter: w, body: new(bytes.Buffer)}
		next.ServeHTTP(ww, r)

		log.Printf("Response: %d %s %s %s %s\nBody: %s", ww.status, r.Method, r.RequestURI, r.RemoteAddr, time.Since(start), ww.body.String())
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
	body   *bytes.Buffer
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}
