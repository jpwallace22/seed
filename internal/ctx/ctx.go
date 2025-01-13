package ctx

import (
	"os"

	"github.com/jpwallace22/seed/pkg/logger"
)

type SeedContext struct {
	Logger logger.Logger
}

func Build(silent bool) *SeedContext {
	return &SeedContext{
		Logger: logger.NewLogger(os.Stdin, os.Stderr, silent),
	}
}
