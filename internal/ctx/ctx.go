package ctx

import (
	"os"

	"github.com/jpwallace22/seed/cmd/flags"
	"github.com/jpwallace22/seed/pkg/logger"
	"github.com/spf13/cobra"
)

type SeedContext struct {
	Logger logger.Logger
	Cobra  *cobra.Command
	Flags  flags.Flags
}

func New(cobra *cobra.Command, flags flags.Flags) *SeedContext {
	return &SeedContext{
		Cobra:  cobra,
		Logger: logger.NewLogger(os.Stdin, os.Stderr, flags.Root.Silent),
		Flags:  flags,
	}
}
