package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func (app *Config) routes() http.Handler {
	r := gin.Default()

	// CORS configuration
	config := cors.Config{
		AllowOrigins:     []string{"https://*", "http://*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Link"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// Heartbeat endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// Prometheus metrics endpoint
	// get global Monitor object
	m := ginmetrics.GetMonitor()

	// +optional set metric path, default /debug/metrics
	m.SetMetricPath("/metrics")
	// +optional set slow time, default 5s
	m.SetSlowTime(10)
	// +optional set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// set middleware for gin
	m.Use(r)

	// Custom logger middleware
	r.Use(LoggerMiddleware)

	// HandleSubmission endpoint
	r.POST("/handle", func(c *gin.Context) {
		app.HandleSubmission(c.Writer, c.Request)
	})

	// Prometheus' metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
