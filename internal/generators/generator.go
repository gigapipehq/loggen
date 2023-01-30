package generators

import (
	"context"

	"github.com/gigapipehq/loggen/internal/generators/loki"
)

type Generator struct {
	rate   int
	labels map[string]string
}

func New(rate int, labels map[string]string) *Generator {
	labels["job"] = "loggen"
	return &Generator{rate: rate, labels: labels}
}

func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	return loki.GenerateLokiLogs(ctx, g.rate, g.labels)
}

func (g *Generator) Rate() int {
	return g.rate
}
