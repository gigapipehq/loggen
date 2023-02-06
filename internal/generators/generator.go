package generators

import (
	"context"

	"github.com/gigapipehq/loggen/internal/generators/loki"
)

type Generator struct {
	lineFmt string
	rate    int
	labels  map[string]string
}

func New(lineFmt string, rate int, labels map[string]string) *Generator {
	labels["job"] = "loggen"
	return &Generator{lineFmt: lineFmt, rate: rate, labels: labels}
}

func (g *Generator) Generate(ctx context.Context) ([]byte, error) {
	return loki.GenerateLokiLogs(ctx, g.lineFmt, g.rate, g.labels)
}

func (g *Generator) Rate() int {
	return g.rate
}
