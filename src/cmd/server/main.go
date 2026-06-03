package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "status_code"},
	)
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "status_code"},
	)
)

func init() {
	prometheus.MustRegister(httpDuration, httpRequests)
}

func instrument(path string, handler func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Simulate slow response if BAD_VERSION env var is set
		if os.Getenv("BAD_VERSION") == "true" {
			time.Sleep(600 * time.Millisecond)
		}

		handler(w, r)
		duration := time.Since(start).Seconds()
		status := "200"
		httpDuration.WithLabelValues(path, status).Observe(duration)
		httpRequests.WithLabelValues(path, status).Add(1)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		version = "v1.0.0"
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Hello from GitOps Platform",
		"version": version,
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", instrument("/healthz", healthz))
	mux.HandleFunc("/api/v1/data", instrument("/api/v1/data", dataHandler))
	mux.Handle("/metrics", promhttp.Handler())

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
