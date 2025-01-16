package main

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/internal/runner"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "seed [string]",
	Short:   "Plant the seeds of your directory tree 🌱.",
	Long:    "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via string or clipboard.",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		runner := runner.NewRootRunner(flags)
		return runner.Run(flags.FromClipboard, args)
	},
}

var flags runner.RootFlags

func init() {
	rootCmd.Flags().BoolVarP(&flags.FromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
	rootCmd.Flags().BoolVarP(&flags.Silent, "silent", "s", false, "If true, suppresses all non-essential console output.")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
