package generators

import (
	"context"

	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators/loki"
)

type Generator struct {
	logConfig config.LogConfig
	rate      int
	labels    map[string]string
}

func New(logConfig config.LogConfig, rate int, labels map[string]string) *Generator {
	labels["job"] = "loggen"
	labels["format"] = logConfig.Format
	return &Generator{logConfig: logConfig, rate: rate, labels: labels}
}

func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	return loki.GenerateLokiLogs(ctx, g.logConfig, g.rate, g.labels)
}

func (g *Generator) Rate() int {
	return g.rate
}
