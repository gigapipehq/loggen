package loki

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/go-logfmt/logfmt"

	"github.com/gigapipehq/loggen/internal/otel"
)

type log struct {
	Streams []stream `json:"streams"`
}

type stream struct {
	Stream map[string]string `json:"stream"`
	Values [][]string        `json:"values"`
}

type logLine struct {
	TraceID    string `json:"traceId"`
	Level      string `json:"level"`
	Host       string `json:"host"`
	Method     string `json:"method"`
	StatusCode int    `json:"status_code"`
	Bytes      int    `json:"bytes"`
}

func (l *logLine) toLogFMT() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	e := logfmt.NewEncoder(buf)
	if err := e.EncodeKeyvals(
		"traceId", l.TraceID,
		"level", l.Level,
		"host", l.Host,
		"method", l.Method,
		"status_code", l.StatusCode,
		"bytes", l.Bytes,
	); err != nil {
		return "", err
	}
	if err := e.EndRecord(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func GenerateLokiLogs(ctx context.Context, lineFmt string, count int, labels map[string]string) ([]byte, error) {
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

	marshalLine := func(l logLine) string {
		s, _ := l.toLogFMT()
		return s
	}
	if lineFmt == "json" {
		marshalLine = func(l logLine) string {
			b, _ := l.MarshalJSON()
			return string(b)
		}
	}

	for i := 0; i < count; i++ {
		l.Streams[0].Values[i] = []string{
			fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		ll := logLine{
			TraceID:    span.SpanContext().TraceID().String(),
			Level:      gofakeit.LogLevel("general"),
			Host:       gofakeit.DomainName(),
			Method:     gofakeit.HTTPMethod(),
			StatusCode: gofakeit.HTTPStatusCodeSimple(),
			Bytes:      gofakeit.Number(0, 3000),
		}
		l.Streams[0].Values[i] = append(l.Streams[0].Values[i], marshalLine(ll))
	}
	return l.MarshalJSON()
}
