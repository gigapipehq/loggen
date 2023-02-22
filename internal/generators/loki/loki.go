package loki

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-logfmt/logfmt"
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

func GenerateLokiLogs(ctx context.Context, logConfig config.LogConfig, count int, labels map[string]string) ([]byte, error) {
	l := log{
		Streams: []stream{
			{
				Stream: labels,
				Values: make([][]string, count),
			},
		},
	}
	ctx, span := otel.Tracer.Start(ctx, "generate loki logs batch")
	defer span.End()

	marshalLine := config.GetLogLineMarshaller[logLine](logConfig)
	rand := gofakeit.New(0)
	for i := 0; i < count; i++ {
		l.Streams[0].Values[i] = []string{
			fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		line := generateLine(rand, logConfig.Structure, span.SpanContext().TraceID())
		line["spanId"] = span.SpanContext().SpanID().String()
		l.Streams[0].Values[i] = append(l.Streams[0].Values[i], marshalLine(line))
	}
	return l.MarshalJSON()
}

func GenerateLokiExampleLog(logConfig config.LogConfig) []byte {
	marshalLine := config.GetLogLineMarshaller[logLine](logConfig)
	return []byte(marshalLine(generateLine(gofakeit.New(0), logConfig.Structure, [16]byte{1})))
}

func generateLine(rand *gofakeit.Faker, structure map[string]string, traceID trace.TraceID) logLine {
	line := logLine{}
	for k, v := range structure {
		line[k] = rand.Generate(fmt.Sprintf("{%s}", v))
	}
	line["traceId"] = traceID.String()
	return line
}
