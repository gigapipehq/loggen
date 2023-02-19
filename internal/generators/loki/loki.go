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

//easyjson:json
type logLine map[string]string

func (l logLine) toLogFMT() (string, error) {
	buf := bytes.NewBuffer([]byte{})
	e := logfmt.NewEncoder(buf)
	for k, v := range l {
		if err := e.EncodeKeyval(k, v); err != nil {
			return "", err
		}
	}
	return buf.String(), nil
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

	marshalLine := func(l logLine) string {
		s, _ := l.toLogFMT()
		return s
	}
	if logConfig.Format == "json" {
		marshalLine = func(l logLine) string {
			b, _ := l.MarshalJSON()
			return string(b)
		}
	}

	rand := gofakeit.New(0)
	for i := 0; i < count; i++ {
		l.Streams[0].Values[i] = []string{
			fmt.Sprintf("%d", time.Now().UnixNano()),
		}
		line := generateLine(rand, logConfig.Structure)
		l.Streams[0].Values[i] = append(l.Streams[0].Values[i], marshalLine(line))
	}
	return l.MarshalJSON()
}

func GenerateLokiExampleLog(logConfig config.LogConfig) []byte {
	marshalLine := func(l logLine) string {
		s, _ := l.toLogFMT()
		return s
	}
	if logConfig.Format == "json" {
		marshalLine = func(l logLine) string {
			b, _ := l.MarshalJSON()
			return string(b)
		}
	}
	return []byte(marshalLine(generateLine(gofakeit.New(0), logConfig.Structure)))
}

func generateLine(rand *gofakeit.Faker, structure map[string]string) logLine {
	line := logLine{}
	for k, v := range structure {
		line[k] = rand.Generate(fmt.Sprintf("{%s}", v))
	}
	line["traceId"] = trace.TraceID([16]byte{1}).String()
	return line
}
