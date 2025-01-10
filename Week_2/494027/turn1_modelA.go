package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-telemetry/otel"
	"github.com/google/go-telemetry/otel/exporters/otlp/otlptrace/otlptracehttp"
	"github.com/google/go-telemetry/otel/sdk/trace"
	"github.com/google/go-telemetry/otel/trace/export/traceotlp"
)

// Example handler function
func projectHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fmt.Fprintf(w, "Project: %s\n", r.URL.Query().Get("project"))

	// Simulate a time-consuming operation
	time.Sleep(2 * time.Second)
}

// Middleware to add tracing to handlers
func traceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := otel.StartSpan(r.Context(), "Project Request")
		defer span.End()

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Initialize the OpenTelemetry SDK
	exp, err := otlptracehttp.NewOTLPTraceExporter(
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("http://localhost:4317"),
	)
	if err != nil {
		panic(err)
	}

	defer exp.Shutdown(context.Background())

	// Create a trace provider with the exporter
	tp := trace.NewSDKTraceProvider(
		trace.WithTraceIDHandler(trace.AutoIDHandler()),
		trace.WithSpanProcessor(traceotlp.NewSpanProcessor(exp)),
	)

	// Set the global trace provider
	otel.SetTracerProvider(tp)
	defer tp.Shutdown(context.Background())

	// Define the routing
	http.HandleFunc("/project", traceMiddleware(http.HandlerFunc(projectHandler)))

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
