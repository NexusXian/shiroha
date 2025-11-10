package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// docCmd generates Swagger documentation by scanning project annotations.
var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate Swagger API documentation",
	Long:  "Scan Go annotation comments across the project and generate Swagger API documentation automatically.",
	RunE: func(cmd *cobra.Command, args []string) error {

		// Detect project root
		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get working directory: %w", err)
		}

		// If currently inside ./cmd then go to root
		if filepath.Base(projectRoot) == "cmd" {
			projectRoot = filepath.Dir(projectRoot)
		}

		fmt.Println("📄 Generating Swagger documentation...")

		// Ensure swag CLI exists
		if _, err := exec.LookPath("swag"); err != nil {
			fmt.Println("⚠️ swag command not found. Installing...")

			install := exec.Command("go", "install", "github.com/swaggo/swag/cmd/swag@latest")
			install.Stdout = os.Stdout
			install.Stderr = os.Stderr
			if err := install.Run(); err != nil {
				return fmt.Errorf("failed to install swag CLI: %w", err)
			}

			fmt.Println("✅ swag installed successfully")
		}

		// Execute swag init
		generate := exec.Command("swag", "init", "-g", "cmd/server.go", "-o", "docs")
		generate.Dir = projectRoot
		generate.Stdout = os.Stdout
		generate.Stderr = os.Stderr

		if err := generate.Run(); err != nil {
			return fmt.Errorf("failed to generate documentation: %w", err)
		}

		fmt.Println("✅ Swagger documentation generated in /docs")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(docCmd)
}
