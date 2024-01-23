package main

import (
	"context"
	"log"
	"os"

	"github.com/denys-kovalchuk/kbot/cmd"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	teleToken := os.Getenv("TELE_TOKEN")
	if teleToken == "" {
		log.Fatal("Please set the TELE_TOKEN environment variable.")
	}

	meterProvider := metric.NewMeterProvider()

	traceProvider := trace.NewTracerProvider()

	cmd.TeleToken = teleToken
	cmd.MeterProvider = meterProvider
	cmd.TraceProvider = traceProvider
	cmd.Execute()

	meterProvider.Stop()
	traceProvider.Shutdown(context.Background())
}
