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

// newCmd defines the 'shiroha new <project-name>' command
var newCmd = &cobra.Command{
	Use:     "new [project-name]",
	Aliases: []string{"generate", "init"},
	Short:   "Create a new Go project structure",
	Long:    "Create a new Go project structure using the specified layered architecture template.",
	Args:    cobra.ExactArgs(1), // Project name must be provided
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		fmt.Printf("Creating project: %s...\n", projectName)

		// Create the project structure
		if err := createProjectStructure(projectName); err != nil {
			return err
		}

		fmt.Printf("✅ Project '%s' generated successfully!\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("cd %s\n", projectName)
		fmt.Println("go mod tidy")
		fmt.Println("go run cmd/server.go")

		// Ask the user whether to automatically execute 'go mod tidy' and start the project
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nDo you want to enter the directory, run 'go mod tidy' and start the project? (y/n): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			projectDir := filepath.Join(".", projectName)

			// run go mod tidy
			fmt.Println("\nExecuting 'go mod tidy'...")
			tidyCmd := exec.Command("go", "mod", "tidy")
			tidyCmd.Dir = projectDir
			tidyCmd.Stdout = os.Stdout
			tidyCmd.Stderr = os.Stderr
			if err := tidyCmd.Run(); err != nil {
				return fmt.Errorf("failed to run 'go mod tidy': %w", err)
			}
			fmt.Println("✅ 'go mod tidy' completed successfully")

			// Start the project (displaying Gin logs in real-time)
			fmt.Println("\nStarting the project...")
			runCmd := exec.Command("go", "run", "cmd/server.go")
			runCmd.Dir = projectDir
			runCmd.Stdout = os.Stdout
			runCmd.Stderr = os.Stderr
			if err := runCmd.Run(); err != nil {
				return fmt.Errorf("failed to start project: %w", err)
			}
		}

		return nil
	},
}

// createProjectStructure is the core function: creates directories and files
func createProjectStructure(name string) error {
	dirs := []string{
		"cmd",
		"internal/database",
		"internal/handler",
		"internal/model",
		"internal/repository",
		"internal/request",
		"internal/response",
		"internal/router",
		"internal/service",
		"pkg/common",
		"pkg/jwt",
		"pkg/utils",
		"config",
		"sql",
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(name, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	// go.mod file content
	goModContent := fmt.Sprintf(`module %s

go 1.18

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/spf13/viper v1.18.2
)`, name)
	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// cmd/server.go file content
	serverContent := fmt.Sprintf(`package main

import (
	"fmt"
	"log"

	"%s/config"
	"%s/internal/router"

	"github.com/gin-gonic/gin"
)

func init() {
	config.LoadConfig()
}

func main() {
	r := router.InitRouter()

	r.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello Shiroha",
		})
	})

	port := config.Cfg.Server.Port
	fmt.Printf("Server running at http://localhost:%%d\\n", port)
	if err := r.Run(fmt.Sprintf(":%%d", port)); err != nil {
		log.Printf("server start error: %%v", err)
	}
}`, name, name)
	if err := os.WriteFile(filepath.Join(name, "cmd", "server.go"), []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to write cmd/server.go: %w", err)
	}

	// internal/router/main_router.go file content
	routerContent := `package router

import (
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	return r
}`
	if err := os.WriteFile(filepath.Join(name, "internal", "router", "main_router.go"), []byte(routerContent), 0644); err != nil {
		return fmt.Errorf("failed to write internal/router/main_router.go: %w", err)
	}

	// README.md file content
	readmeContent := fmt.Sprintf(`# %s

This project uses Go language with a layered architecture.

## Structure Explanation
- internal/handler: Interface layer (Controller)
- internal/service: Business logic layer
- internal/repository: Data access layer
- internal/request/response: DTOs

## Startup
1. cd %s
2. go mod tidy
3. go run cmd/server.go
`, name, name)
	if err := os.WriteFile(filepath.Join(name, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// config.yaml file content
	configYamlContent := `server:
  port: 8080`
	if err := os.WriteFile(filepath.Join(name, "config.yaml"), []byte(configYamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

	// config/config.go file content
	configGoContent := `package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server *ServerConfig ` + "`mapstructure:\"server\"`" + `
}

type ServerConfig struct {
	Port int ` + "`mapstructure:\"port\"`" + `
}

var Cfg *Config

func LoadConfig() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")

	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("../config")
	v.AddConfigPath("..")

	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %%v", err)
	}

	Cfg = &Config{}
	if err := v.Unmarshal(Cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %%v", err)
	}
}`
	if err := os.WriteFile(filepath.Join(name, "config", "config.go"), []byte(configGoContent), 0644); err != nil {
		return fmt.Errorf("failed to write config/config.go: %w", err)
	}

	return nil
}
