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

	ctx, span := otel.Tracer.Start(ctx, opName)
	defer span.End()

	gen := func(ctx context.Context) senders.Generator {
		_, span := otel.Tracer.Start(ctx, "create new generator")
		defer span.End()
		return generators.New(cfg.LogConfig, cfg.Rate, cfg.Labels)
	}(ctx)

	s, err := func(ctx context.Context) (senders.Sender, error) {
		_, span := otel.Tracer.Start(ctx, "create new sender")
		defer span.End()
		auth := map[string]string{
			"X-API-Key":    cfg.APIKey,
			"X-API-Secret": cfg.APISecret,
		}
		curl := fmt.Sprintf("%s/loki/api/v1/push", cfg.URL)
		return _default.New().WithHeaders(auth).WithURL(curl)
	}(ctx)
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
	if cfg.EnableTraces {
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
