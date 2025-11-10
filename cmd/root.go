package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "shiroha",
	Short: "A powerful CLI tool for bootstrapping and managing Gin-based web APIs.",
	Long: `Shiroha is a robust command-line interface (CLI) designed to accelerate Go web development using the Gin framework.

It automatically sets up a clean, layered (e.g., domain, service, repository) project architecture, ensuring maintainability, testability, and scalability for your RESTful APIs.

Usage:
  shiroha new <project-name>   - Quickly generate a new Gin web project with best practices.
  shiroha build                - Compile and package the project for deployment across various platforms.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

// Execute serves as the entry point for all commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("❌ Error occurred:", err)
		os.Exit(1)
	}
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(docCmd)
	rootCmd.AddCommand(runCmd)
}
