package generators

import (
	"github.com/gigapipehq/loggen/internal/config"
	"github.com/gigapipehq/loggen/internal/generators/loki"
)

type Generator struct {
	config *config.Config
}

func New(cfg *config.Config) *Generator {
	return &Generator{config: cfg}
}

func (g *Generator) Generate() ([]byte, error) {
	return loki.GenerateLokiLogs(g.config)
}

func (g *Generator) Rate() int {
	return g.config.Rate
}

func Example(logConfig config.LogConfig) []byte {
	return loki.GenerateLokiExampleLog(logConfig)
}
