package otel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/denisbrodbeck/machineid"
	"github.com/matishsiao/goInfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	"github.com/gigapipehq/loggen/internal/config"
)

var Tracer = otel.Tracer("loggen")

func NewExporter(collectorURL string, httpClient *http.Client) sdktrace.SpanExporter {
	curl := fmt.Sprintf("%s/tempo/spans", collectorURL)
	e, _ := zipkin.New(curl, zipkin.WithClient(httpClient))
	return e
}

func NewNoopExporter() sdktrace.SpanExporter {
	e, _ := stdouttrace.New(
		stdouttrace.WithWriter(io.Discard),
		stdouttrace.WithoutTimestamps(),
	)
	return e
}

func NewProvider(exporter sdktrace.SpanExporter, cfg *config.Config) *sdktrace.TracerProvider {
	mid, err := machineid.ID()
	if err != nil {
		mid = "00000000-0000-0000-0000-000000000000"
	}

	info, _ := goInfo.GetInfo()
	labels := []string{}
	for k, v := range cfg.Labels {
		labels = append(labels, fmt.Sprintf("%s: %s", k, v))
	}
	baseResource, _ := resource.New(
		context.Background(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithHost(),
		resource.WithTelemetrySDK(),
	)
	r, _ := resource.Merge(
		baseResource,
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("Loggen"),
			semconv.DeviceIDKey.String(mid),
			semconv.HostArchKey.String(runtime.GOARCH),
			attribute.Key("host.cpus").Int(info.CPUs),
			semconv.TelemetrySDKLanguageGo,
			attribute.Key("config.url").String(cfg.URL),
			attribute.Key("config.api_key").String(cfg.APIKey),
			attribute.Key("config.api_secret").String(cfg.APISecret),
			attribute.Key("config.labels").StringSlice(labels),
			attribute.Key("config.rate").Int(cfg.Rate),
			attribute.Key("config.timeout").String(cfg.Timeout.String()),
		),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)
}
