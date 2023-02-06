package cmd

import (
	"context"
	"fmt"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators"
	"github.com/gigapipehq/loggen/internal/otel"
	"github.com/gigapipehq/loggen/internal/senders"
	_default "github.com/gigapipehq/loggen/internal/senders/default"
)

type progressTracker interface {
	Add(int)
}

func Do(cfg *config.Config, opName string, progress progressTracker) error {
	ctx, span := otel.Tracer.Start(context.Background(), opName)
	defer span.End()

	gen := func(ctx context.Context) senders.Generator {
		_, span := otel.Tracer.Start(ctx, "create new generator")
		defer span.End()
		return generators.New(cfg.Format, cfg.Rate, cfg.Labels)
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
