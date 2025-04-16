// internal/api/metrics.go
package api

import (
    "sync"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "net/http"
)

var (
    httpRequestsTotal       *prometheus.CounterVec
    httpRequestDuration     *prometheus.HistogramVec
    kafkaMessagesProduced   prometheus.Counter
    kafkaMessagesConsumed   prometheus.Counter

    registerMetricsOnce sync.Once
)

func initMetrics() {
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests processed, labeled by status and method",
        },
        []string{"status", "method"},
    )

    httpRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Histogram of latencies for HTTP requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method"},
    )

    kafkaMessagesProduced = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "kafka_messages_produced_total",
            Help: "Total number of messages produced to Kafka",
        },
    )

    kafkaMessagesConsumed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "kafka_messages_consumed_total",
            Help: "Total number of messages consumed from Kafka",
        },
    )

    prometheus.MustRegister(httpRequestsTotal)
    prometheus.MustRegister(httpRequestDuration)
    prometheus.MustRegister(kafkaMessagesProduced)
    prometheus.MustRegister(kafkaMessagesConsumed)
}

// RegisterMetrics initializes and registers Prometheus metrics only once
func RegisterMetrics() {
    registerMetricsOnce.Do(initMetrics)
}

// MetricsHandler exposes the /metrics endpoint
func MetricsHandler() http.Handler {
    return promhttp.Handler()
}
