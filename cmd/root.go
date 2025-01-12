/*
Copyright © 2025 Justin Wallace <jpwallace22@gmai.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/jpwallace22/seed/pkg/logger"
	"github.com/jpwallace22/seed/pkg/parser"
	"github.com/spf13/cobra"
	clipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

var (
	filePath      string
	fromClipboard bool
)

func init() {
	rootCmd.Flags().BoolVarP(&fromClipboard, "clipboard", "c", false, "Use tree structure from clipboard.")
}

var rootCmd = &cobra.Command{
	Use:   "seed [path/to/tree]",
	Short: "Plant the seeds of your directory tree.",
	Long:  "Seed is a CLI tool that helps you grow directory structures from a tree representation provided via file or clipboard.",
	Args:  cobra.MaximumNArgs(1),
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	if fromClipboard {
		fmt.Println("Planting from clipboard...")
		text, err := getClipboardContent()
		if err != nil {
			logger.Error("Unable to access clipboard contents: %v", err)
			os.Exit(1)
		}

		err = parser.ParseTreeString(text)
		if err != nil {
			logger.Error("Unable to parse the tree structure: %v", err)
		}
	} else if len(args) > 0 {
		fmt.Println("Sprouting directories from file:", filePath)
		// Logic to read and process the file
	} else {
		fmt.Println("Error: No seeds provided. Provide a path or use --clipboard (-c).")
		os.Exit(1)
	}
	logger.Success("Your directory tree has grown successfully!")
}

func getClipboardContent() (string, error) {
	c := clipboard.New()
	text, err := c.PasteText()
	if err != nil {
		return "", fmt.Errorf("clipboard read error: %w", err)
	}

	return text, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
