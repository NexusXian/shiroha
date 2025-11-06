package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Supported target platform mapping
var platforms = map[string]struct{ os, arch, desc string }{
	"1": {"linux", "amd64", "Linux (amd64)"},
	"2": {"linux", "arm64", "Linux (arm64)"},
	"3": {"darwin", "amd64", "macOS Intel (amd64)"},
	"4": {"darwin", "arm64", "macOS Apple Silicon (arm64)"},
	"5": {"windows", "amd64", "Windows (amd64)"},
}

// buildCmd defines the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build project for specified platform",
	Long:  "Build the Go project into an executable binary for the selected target platform.",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Display platform selection menu
		fmt.Println("Select target platform:")
		for i := 1; i <= len(platforms); i++ {
			fmt.Printf("%d. %s\n", i, platforms[fmt.Sprintf("%d", i)].desc)
		}
		fmt.Print("Enter your choice (1-5): ")

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}
		choice := strings.TrimSpace(input)

		// Validate the choice
		platform, exists := platforms[choice]
		if !exists {
			return fmt.Errorf("invalid choice: %s, please enter 1-%d", choice, len(platforms))
		}

		// Get the current project root directory (one level above 'cmd')
		projectRoot, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}

		// If the current directory is 'cmd', go up one level
		if filepath.Base(projectRoot) == "cmd" {
			projectRoot = filepath.Dir(projectRoot)
		}

		projectName := filepath.Base(projectRoot)

		// Determine the output file name
		outputName := projectName
		if platform.os == "windows" {
			outputName += ".exe"
		}

		// Create the bin directory
		binDir := filepath.Join(projectRoot, "bin")
		if err := os.MkdirAll(binDir, 0755); err != nil {
			return fmt.Errorf("failed to create bin directory: %w", err)
		}

		outputPath := filepath.Join(binDir, outputName)

		// Execute the build command (using relative path for robustness)
		fmt.Printf("\nBuilding for %s/%s...\n", platform.os, platform.arch)
		buildCmd := exec.Command(
			"go", "build",
			"-o", outputPath,
			"./cmd/server.go",
		)

		// Set cross-compilation environment variables
		buildCmd.Env = append(os.Environ(),
			fmt.Sprintf("GOOS=%s", platform.os),
			fmt.Sprintf("GOARCH=%s", platform.arch),
		)

		// Execute and capture output
		output, err := buildCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("build failed: %w\nOutput: %s", err, string(output))
		}

		fmt.Printf("✅ Build successful! Output file: %s\n", outputPath)
		return nil
	},
}
