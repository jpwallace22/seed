package main

import (
	"fmt"
	"os"

	cmdFlags "github.com/jpwallace22/seed/cmd/flags"
	"github.com/jpwallace22/seed/internal/ctx"
	"github.com/jpwallace22/seed/internal/runner"
	"github.com/spf13/cobra"
)

var flags cmdFlags.Flags

func init() {
	// Persistent Flags
	rootCmd.PersistentFlags().BoolVarP(&flags.Root.Silent, "silent", "s", false, "If true, suppresses all non-essential console output.")

	// Command Flags
	rootCmd.Flags().BoolVarP(&flags.Root.FromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
	rootCmd.Flags().StringVarP(&flags.Root.FilePath, "file", "f", "", "Use tree structure from a file.")
	rootCmd.Flags().VarP(&flags.Root.Format, "format", "F", "Format of the input [tree, json, yaml]")
}

var rootCmd = &cobra.Command{
	Version: "0.1.1", // unreleased
	Use:     "seed [string]",
	Short:   "Plant the seeds of your directory tree 🌱.",
	Long:    "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via string or clipboard.",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := ctx.New(cmd, flags)
		runner := runner.NewRootRunner(cmd, ctx)
		return runner.Run(args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
