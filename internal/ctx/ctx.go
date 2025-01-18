package ctx

import (
	"os"

	"github.com/jpwallace22/seed/pkg/logger"
	"github.com/spf13/cobra"
)

type Config struct {
	Silent bool
}

type SeedContext struct {
	Logger logger.Logger
	Cobra  *cobra.Command
	Config Config
}

func Build(cobra *cobra.Command, silent bool) *SeedContext {
	return &SeedContext{
		Cobra:  cobra,
		Logger: logger.NewLogger(os.Stdin, os.Stderr, silent),
		Config: Config{
			Silent: silent,
		},
	}
}
