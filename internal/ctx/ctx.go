package ctx

import (
	"os"

	"github.com/jpwallace22/seed/pkg/logger"
)

type GlobalFlags struct {
	Silent bool
}

type SeedContext struct {
	Logger      logger.Logger
	GlobalFlags GlobalFlags
}

func Build(silent bool) *SeedContext {
	return &SeedContext{
		Logger: logger.NewLogger(os.Stdin, os.Stderr, silent),
		GlobalFlags: GlobalFlags{
			Silent: silent,
		},
	}
}
