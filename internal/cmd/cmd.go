package cmd

import (
	"context"
	"fmt"

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

func Do(ctx context.Context, cfg *config.Config, opName string, progress progressTracker) error {
	shutdownMT := configureTracesAndTraces(ctx, cfg)
	defer shutdownMT()

	gen := generators.New(cfg)
	s, err := _default.New().WithHeaders(map[string]string{
		"X-API-Key":    cfg.APIKey,
		"X-API-Secret": cfg.APISecret,
	}).WithURL(fmt.Sprintf("%s/loki/api/v1/push", cfg.URL))
	if err != nil {
		return fmt.Errorf("unable to create sender: %v", err)
	}

	go func() {
		for i := range s.Progress() {
			progress.Add(i)
		}
	}()
	ctx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()
	senders.Start(ctx, s, gen)
	return nil
}

func configureTracesAndTraces(ctx context.Context, cfg *config.Config) func() {
	if cfg.EnableMetrics {
		prom.Initialize(ctx, cfg)
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
		if tp != nil {
			_ = tp.Shutdown(context.Background())
		}
	}
}
