package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"time"
)

func LoggerMiddleware(c *gin.Context) {
	start := time.Now()

	// Read the body
	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Printf("Incoming request: %s %s %s\nBody: %s", c.Request.Method, c.Request.RequestURI, c.ClientIP(), string(bodyBytes))

	// Use a custom response writer to capture the response body
	ww := &responseWriter{ResponseWriter: c.Writer, body: new(bytes.Buffer)}
	c.Writer = ww

	// Process request
	c.Next()

	// Log the response
	log.Printf("Response: %d %s %s %s %s\nBody: %s", ww.status, c.Request.Method, c.Request.RequestURI, c.ClientIP(), time.Since(start), ww.body.String())
}

type responseWriter struct {
	gin.ResponseWriter
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

func (rw *responseWriter) WriteString(s string) (int, error) {
	rw.body.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}
