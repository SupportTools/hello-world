package metrics

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/supporttools/hello-world/pkg/config"
	"github.com/supporttools/hello-world/pkg/health"
	"github.com/supporttools/hello-world/pkg/logging"
)

// logger is the global logger for the metrics package.
var logger = logging.SetupLogging(&config.CFG)

var (
	// Register a counter metric for counting the total number of requests
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of get requests.",
		},
		[]string{"path"},
	)

	// Register a histogram to observe the response times
	responseDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_duration_seconds",
			Help:    "Duration of HTTP responses.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)
)

// init registers the metrics with Prometheus.
func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseDuration)
}

// StartMetricsServer starts the metrics server on the specified port.
func StartMetricsServer(port int) {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/healthz", health.HealthzHandler())
	mux.Handle("/version", health.VersionHandler())

	serverPort := strconv.Itoa(port)
	logger.Printf("Metrics server starting on port %d\n", port)

	if err := http.ListenAndServe(":"+serverPort, mux); err != nil {
		logger.Fatalf("Metrics server failed to start: %v", err)
	}
}

// RecordMetrics records the metrics for the given path and duration.
func RecordMetrics(path string, duration float64) {
	totalRequests.WithLabelValues(path).Inc()
	responseDuration.WithLabelValues(path).Observe(duration)
}
