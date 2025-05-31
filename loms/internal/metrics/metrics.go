package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "loms",
			Name:      "grpc_requests_total",
			Help:      "Total number of gRPC requests",
		},
		[]string{"method"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "grpc_request_duration_seconds",
			Help:      "gRPC request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "status"},
	)

	ExternalRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "loms",
			Name:      "external_requests_total",
			Help:      "Total number of requests to external services",
		},
		[]string{"service", "method"},
	)

	ExternalRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "external_request_duration_seconds",
			Help:      "External service request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"service", "method", "status"},
	)

	DatabaseOperationsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "loms",
			Name:      "database_operations_total",
			Help:      "Total number of database operations",
		},
		[]string{"operation"},
	)

	DatabaseOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "loms",
			Name:      "database_operation_duration_seconds",
			Help:      "Database operation duration in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"operation", "status"},
	)
)

func NewMetricsHandler() http.Handler {
	return promhttp.Handler()
}

func ObserveGRPCRequest(method string, statusCode int, duration time.Duration) {
	RequestsTotal.WithLabelValues(method).Inc()
	RequestDuration.WithLabelValues(
		method,
		strconv.Itoa(statusCode),
	).Observe(duration.Seconds())
}

func ObserveExternalRequest(service, method string, statusCode int, duration time.Duration) {
	ExternalRequestsTotal.WithLabelValues(service, method).Inc()
	ExternalRequestDuration.WithLabelValues(
		service,
		method,
		strconv.Itoa(statusCode),
	).Observe(duration.Seconds())
}

func ObserveDatabaseOperation(operation string, err error, duration time.Duration) {
	DatabaseOperationsTotal.WithLabelValues(operation).Inc()

	status := "success"
	if err != nil {
		status = "error"
	}

	DatabaseOperationDuration.WithLabelValues(
		operation,
		status,
	).Observe(duration.Seconds())
}
