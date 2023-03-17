package cmd

import (
	"context"
	"fmt"
	"time"

	otelsdk "go.opentelemetry.io/otel"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
	"github.com/gigapipehq/loggen/internal/otel"
	"github.com/gigapipehq/loggen/internal/prom"
	"github.com/gigapipehq/loggen/internal/senders"
	_default "github.com/gigapipehq/loggen/internal/senders/default"
)

type progressTracker interface {
	Add(int)
}

func Do(ctx context.Context, cfg *config.Config, progress progressTracker) error {
	shutdownMT := configureMetricsAndTraces(ctx, cfg)
	defer shutdownMT()

	gen := generators.New(cfg)
	s, err := _default.New().WithHeaders(cfg.GetHeaders()).WithURL(fmt.Sprintf("%s/loki/api/v1/push", cfg.URL))
	if err != nil {
		return fmt.Errorf("unable to create sender: %v", err)
	}

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

func configureMetricsAndTraces(ctx context.Context, cfg *config.Config) func() {
	pctx, cancel := context.WithCancel(ctx)
	var qch chan struct{}
	if cfg.EnableMetrics {
		qch = prom.Initialize(pctx, cfg)
	}
	exporter := otel.NewNoopExporter()
	if cfg.Traces.Enabled {
		sender := _default.New().WithHeaders(map[string]string{
			"X-API-Key":    cfg.APIKey,
			"X-API-Secret": cfg.APISecret,
		})
		exporter = otel.NewExporter(cfg.URL, sender.Client())
	}
	tp := otel.NewProvider(exporter, cfg)
	otelsdk.SetTracerProvider(tp)

	return func() {
		cancel()
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
		if cfg.EnableMetrics {
			<-qch
			close(qch)
		}
	}
}
