package otel

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"sync/atomic"

	"github.com/denisbrodbeck/machineid"
	"github.com/matishsiao/goInfo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/gigapipehq/loggen/internal/config"
)

var (
	Tracer trace.Tracer

	knownResourceWithAttributes = map[attribute.Key]func() resource.Option{
		semconv.HostNameKey:                  resource.WithHost,
		semconv.OSDescriptionKey:             resource.WithOSDescription,
		semconv.OSTypeKey:                    resource.WithOSType,
		semconv.ProcessExecutableNameKey:     resource.WithProcessExecutableName,
		semconv.ProcessExecutablePathKey:     resource.WithProcessExecutablePath,
		semconv.ProcessCommandArgsKey:        resource.WithProcessCommandArgs,
		semconv.ProcessOwnerKey:              resource.WithProcessOwner,
		semconv.ProcessPIDKey:                resource.WithProcessPID,
		semconv.ProcessRuntimeDescriptionKey: resource.WithProcessRuntimeDescription,
		semconv.ProcessRuntimeNameKey:        resource.WithProcessRuntimeName,
		semconv.ProcessRuntimeVersionKey:     resource.WithProcessRuntimeVersion,
	}
	knownBaseAttributes = []attribute.KeyValue{
		semconv.HostArchKey.String(runtime.GOARCH),
		semconv.TelemetrySDKNameKey.String("opentelemetry"),
		semconv.TelemetrySDKLanguageKey.String("go"),
		semconv.TelemetrySDKVersionKey.String(otel.Version()),
	}
	totalSpansSent int64 = 0
)

func init() {
	mid, err := machineid.ID()
	if err != nil {
		mid = "00000000-0000-0000-0000-000000000000"
	}
	info, _ := goInfo.GetInfo()
	knownBaseAttributes = append(knownBaseAttributes, semconv.DeviceIDKey.String(mid))
	knownBaseAttributes = append(knownBaseAttributes, attribute.Int("host.cpus", info.CPUs))
}

type zipkinExporterWithLogsWrapper struct {
	*zipkin.Exporter
}

func (z *zipkinExporterWithLogsWrapper) ExportSpans(ctx context.Context, spans []sdktrace.ReadOnlySpan) error {
	atomic.AddInt64(&totalSpansSent, int64(len(spans)))
	return z.Exporter.ExportSpans(ctx, spans)
}

func NewExporter(collectorURL string, httpClient *http.Client) sdktrace.SpanExporter {
	curl := fmt.Sprintf("%s/tempo/spans", collectorURL)
	e, _ := zipkin.New(curl, zipkin.WithClient(httpClient))
	return &zipkinExporterWithLogsWrapper{e}
}

func GetTotalSpansSent() int64 {
	return atomic.LoadInt64(&totalSpansSent)
}

func NewNoopExporter() sdktrace.SpanExporter {
	e, _ := stdouttrace.New(
		stdouttrace.WithWriter(io.Discard),
		stdouttrace.WithoutTimestamps(),
	)
	return e
}

func NewProvider(exporter sdktrace.SpanExporter, cfg *config.Config) *sdktrace.TracerProvider {
	opts := []resource.Option{}
	attrs := []attribute.KeyValue{}
	for _, k := range cfg.Traces.Defaults {
		if f, ok := knownResourceWithAttributes[k]; ok {
			opts = append(opts, f())
			continue
		}
		for _, kv := range knownBaseAttributes {
			if kv.Key == k {
				attrs = append(attrs, kv)
				break
			}
		}
	}

	var serviceName string
	for _, kv := range cfg.Traces.Custom {
		if kv.Key == semconv.ServiceNameKey {
			serviceName = kv.Value.AsString()
			continue
		}
		attrs = append(attrs, kv)
	}
	if serviceName == "" {
		serviceName = "loggen"
	}
	Tracer = otel.Tracer(serviceName)

	attrs = append(attrs, semconv.ServiceNameKey.String(serviceName))
	baseResource, _ := resource.New(context.Background(), opts...)
	r, _ := resource.Merge(
		baseResource,
		resource.NewWithAttributes(
			semconv.SchemaURL,
			attrs...,
		),
	)
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithMaxExportBatchSize(sdktrace.DefaultMaxExportBatchSize),
			sdktrace.WithMaxQueueSize(100000)),
		sdktrace.WithResource(r),
	)
}
