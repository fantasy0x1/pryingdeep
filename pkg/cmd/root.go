package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/pryingbytez/pryingdeep/pkg/cmd/exporter"
)

var rootCmd = &cobra.Command{
	Use:   "pryingdeep",
	Short: "Pryingdeep is a dark web osint intelligence tool.",
	Long: `Pryingdeep specializes in collecting information about dark-web/clearnet websites.
		This tool was specifically built to extract as much information as possible from a .onion website`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(exporter.ExporterCMD)
	rootCmd.AddCommand(crawlCmd)
}
