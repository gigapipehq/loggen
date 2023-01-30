package loki

import (
	"context"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"

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

func GenerateLokiLogs(ctx context.Context, count int, labels map[string]string) ([]byte, error) {
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

	for i := 0; i < count; i++ {
		_, lspan := otel.Tracer.Start(ctx, "generate batch")
		l.Streams[0].Values[i] = []string{
			fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		ll := logLine{
			TraceID:    lspan.SpanContext().TraceID().String(),
			Level:      gofakeit.LogLevel("general"),
			Host:       gofakeit.DomainName(),
			Method:     gofakeit.HTTPMethod(),
			StatusCode: gofakeit.HTTPStatusCodeSimple(),
			Bytes:      gofakeit.Number(0, 3000),
		}
		b, _ := ll.MarshalJSON()
		l.Streams[0].Values[i] = append(l.Streams[0].Values[i], string(b))
		lspan.End()
	}
	_, mspan := otel.Tracer.Start(ctx, "marshal batch")
	defer mspan.End()
	return l.MarshalJSON()
}
