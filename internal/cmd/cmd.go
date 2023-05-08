package cmd

import (
	"context"
	"fmt"
	"net/url"
	"time"

	otelsdk "go.opentelemetry.io/otel"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
	"github.com/gigapipehq/loggen/internal/otel"
	"github.com/gigapipehq/loggen/internal/prom"
	"github.com/gigapipehq/loggen/internal/senders"
	"github.com/gigapipehq/loggen/internal/senders/file"
	surl "github.com/gigapipehq/loggen/internal/senders/url"
)

type progressTracker interface {
	Add(int)
}

func Do(ctx context.Context, cfg *config.Config, progress progressTracker) error {
	s, err := getSender(cfg)
	if err != nil {
		return fmt.Errorf("unable to create sender: %v", err)
	}

	shutdownMT := configureMetricsAndTraces(ctx, cfg, s)
	defer shutdownMT()

	gen := generators.New(cfg)
	go func() {
		for i := range s.Progress() {
			progress.Add(i)
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(cfg.Timeout))
	defer cancel()
	senders.Start(ctx, s, gen)
	return nil
}

func configureMetricsAndTraces(ctx context.Context, cfg *config.Config, s senders.Sender) func() {
	pctx, cancel := context.WithCancel(ctx)
	var qch chan struct{}

	withMetrics := s.SupportsMetrics() && cfg.EnableMetrics
	if withMetrics {
		qch = prom.Initialize(pctx, cfg)
	}

	tp := otel.NewProvider(s.TracesExporter(), cfg)
	otelsdk.SetTracerProvider(tp)

	return func() {
		cancel()
		if tp != nil {
			_ = tp.Shutdown(context.Background())
			fmt.Printf("Total spans sent: %d\n", otel.GetTotalSpansSent())
		}

		if withMetrics {
			<-qch
			close(qch)
		}
	}
}

func getSender(cfg *config.Config) (senders.Sender, error) {
	u, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to parse sender output location '%s': %v", cfg.URL, err)
	}
	if u.Scheme == "file" {
		s, err := file.New(fmt.Sprintf("%s%s", u.Host, u.Path))
		if err != nil {
			return nil, fmt.Errorf("unable to create file: %v", err)
		}
		return s, nil
	}
	u.Path = "/loki/api/v1/push"
	return surl.New(u).WithHeaders(cfg.GetHeaders()), nil
}
