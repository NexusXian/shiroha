package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// runCmd defines the 'shiroha run' command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs the project's main server (go run cmd/server.go)",
	Long:  "Executes 'go run cmd/server.go' to start the application server from the project root.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Starting project server (go run cmd/server.go)...")

		// Check if the project file exists to ensure execution from the project root
		if _, err := os.Stat("cmd/server.go"); os.IsNotExist(err) {
			return fmt.Errorf("cmd/server.go not found. Please ensure you are running 'shiroha run' from the project root directory")
		}

		runCmd := exec.Command("go", "run", "cmd/server.go")
		runCmd.Stdout = os.Stdout
		runCmd.Stderr = os.Stderr

		if err := runCmd.Run(); err != nil {
			return fmt.Errorf("failed to start project server: %w", err)
		}

		return nil
	},
}

// GetRunCommand returns the Cobra command for 'shiroha run'
func GetRunCommand() *cobra.Command {
	return runCmd
}
