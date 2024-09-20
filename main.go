package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/supporttools/hello-world/pkg/config"
	"github.com/supporttools/hello-world/pkg/logging"
	"github.com/supporttools/hello-world/pkg/metrics"
	"github.com/supporttools/hello-world/pkg/templates"
	"github.com/supporttools/hello-world/pkg/version"
)

var (
	// Global logger variable
	logger *logrus.Logger
)

func main() {
	// Load configuration from environment variables
	config.LoadConfiguration()

	// Set up logging based on the configuration
	logger = logging.SetupLogging(&config.CFG)

	// Log the current configuration if debug mode is enabled
	if config.CFG.Debug {
		logger.Infoln("Debug mode enabled")
		logger.Infof("Port: %d", config.CFG.Port)
		logger.Infof("Metrics Port: %d", config.CFG.MetricsPort)
	}

	// Start the metrics server in a separate goroutine
	go metrics.StartMetricsServer(config.CFG.MetricsPort)

	// Start the web server
	webserver(config.CFG.Port)
}

func webserver(port int) {
	logger.Println("Starting web server...")

	// Set up HTTP handlers
	http.HandleFunc("/", helloWorldHandler)

	// Serve static files from the /img directory
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./img"))))

	// Serve Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())

	// Build the server address
	serverAddress := fmt.Sprintf(":%d", port)
	logger.Printf("Serving HelloWorld on HTTP port: %s\n", serverAddress)

	// Define server with timeouts
	server := &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  10 * 60 * 60, // Max time to read request body
		WriteTimeout: 10 * 60 * 60, // Max time to write response
		IdleTimeout:  60 * 60 * 60, // Max time to keep idle connections open
	}

	// Start the HTTP server with timeout settings
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}

// helloWorldHandler handles requests to the root path and serves the HTML template
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	logger.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
	}).Info("Received request")

	// Prepare template data
	data := map[string]interface{}{
		"Hostname":  getHostname(),
		"GitCommit": version.GitCommit,
		"Host":      r.Host,
		"Headers":   r.Header,
	}

	// Render the HTML template
	output, err := templates.CompileTemplateFromMap(templates.HelloWorldTemplate, data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		logger.Error("Error rendering template: ", err)
		return
	}

	// Write the rendered template to the response
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprint(w, output)
}

// getHostname gets the hostname of the server
func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Error getting hostname: ", err)
		return "unknown"
	}
	return hostname
}
