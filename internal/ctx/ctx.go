package ctx

import (
	"os"

	"github.com/jpwallace22/seed/pkg/logger"
	"github.com/spf13/cobra"
)

type GlobalFlags struct {
	Silent bool
}

type SeedContext struct {
	Logger      logger.Logger
	Cobra       *cobra.Command
	GlobalFlags GlobalFlags
}

func Build(cobra *cobra.Command, silent bool) *SeedContext {
	return &SeedContext{
		Cobra:  cobra,
		Logger: logger.NewLogger(os.Stdin, os.Stderr, silent),
		GlobalFlags: GlobalFlags{
			Silent: silent,
		},
	}
}
