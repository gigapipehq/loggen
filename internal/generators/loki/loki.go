package loki

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-logfmt/logfmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/otel"
)

type log struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type logLine map[string]string

const traceIDKey = "traceId"

func (l logLine) ToLogFMT() string {
	buf := bytes.NewBuffer([]byte{})
	e := logfmt.NewEncoder(buf)
	for k, v := range l {
		if err := e.EncodeKeyval(k, v); err != nil {
			return ""
		}
	}
	return buf.String()
}

func (l logLine) ToJSON() string {
	b, _ := l.MarshalJSON()
	return string(b)
}

func GenerateLokiLogs(cfg *config.Config) ([]byte, error) {
	l := log{
		Streams: []stream{
			{
				Stream: cfg.Labels,
				Values: make([][]string, cfg.Rate),
			},
		},
	}

	marshalLine := config.GetLogLineMarshaller[logLine](cfg.LogConfig)
	rand := gofakeit.New(0)
	for i := 0; i < cfg.Rate; i++ {
		l.Streams[0].Values[i] = []string{
			fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		line := generateLine(rand, cfg.LogConfig.Structure, cfg.Traces)
		l.Streams[0].Values[i] = append(l.Streams[0].Values[i], marshalLine(line))
	}
	return l.MarshalJSON()
}

func GenerateLokiExampleLog(logConfig config.LogConfig) []byte {
	marshalLine := config.GetLogLineMarshaller[logLine](logConfig)
	tc := config.TracesConfig{Enabled: false}
	return []byte(marshalLine(generateLine(gofakeit.New(0), logConfig.Structure, tc)))
}

func generateLine(rand *gofakeit.Faker, structure map[string]string, traceConfig config.TracesConfig) logLine {
	line := logLine{}
	for k, v := range structure {
		line[k] = rand.Generate(fmt.Sprintf("{%s}", v))
	}

	if traceConfig.Enabled {
		ctx := context.Background()
		t := time.Now()
		for _, step := range traceConfig.Spans {
			line[traceIDKey] = createSpan(ctx, step, line, t).String()
		}
		if _, ok := line[traceIDKey]; !ok {
			line[traceIDKey] = trace.TraceID([16]byte{1}).String()
		}
	}
	return line
}

func createSpan(ctx context.Context, step config.SpanStep, line logLine, start time.Time) trace.TraceID {
	lctx, span := otel.Tracer.Start(
		ctx, step.Name, trace.WithSpanKind(trace.ValidateSpanKind(step.Kind)), trace.WithTimestamp(start),
	)
	defer span.End(trace.WithTimestamp(start.Add(step.Duration)))

	attrs := make([]attribute.KeyValue, len(step.Attributes))
	for i, attr := range step.Attributes {
		var strValue string
		if attr.ResolveFromLogValue != "" {
			if _, ok := line[attr.ResolveFromLogValue]; ok {
				strValue = line[attr.ResolveFromLogValue]
			}
		}
		if attr.ResolveFake != "" {
			strValue = gofakeit.New(0).Generate(fmt.Sprintf("{%s}", attr.ResolveFake))
		}
		attrs[i] = convertAttributeValue(attr.Name, attr.ValueType, strValue)
	}
	span.SetAttributes(attrs...)

	d := step.Duration
	for _, child := range step.Children {
		createSpan(lctx, child, line, start.Add(d))
		d += child.Duration
	}
	return span.SpanContext().TraceID()
}

func convertAttributeValue(name, kind, value string) attribute.KeyValue {
	switch kind {
	case "bool":
		b, _ := strconv.ParseBool(value)
		return attribute.Bool(name, b)
	case "int", "int64":
		i, _ := strconv.ParseInt(value, 10, 64)
		return attribute.Int64(name, i)
	case "float64":
		f, _ := strconv.ParseFloat(value, 64)
		return attribute.Float64(name, f)
	default:
		return attribute.String(name, value)
	}
}
