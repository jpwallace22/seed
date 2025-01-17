package main

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/internal/runner"
	"github.com/spf13/cobra"
)

var (
	flags  runner.RootFlags
	config runner.Config
)

func init() {
	// Config
	rootCmd.Flags().BoolVarP(&config.Silent, "silent", "s", false, "If true, suppresses all non-essential console output.")

	// Flags
	rootCmd.Flags().BoolVarP(&flags.FromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
	rootCmd.Flags().StringVarP(&flags.FilePath, "file", "f", "", "Use tree structure from a file.")
	// rootCmd.Flags().VarP(&flags.Format, "format", "F", "Format of the input [tree, json, yamlj]") // TODO: implement different parsers
}

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "seed [string]",
	Short:   "Plant the seeds of your directory tree 🌱.",
	Long:    "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via string or clipboard.",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		runner := runner.NewRootRunner(cmd, config.Silent)
		return runner.Run(flags, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
