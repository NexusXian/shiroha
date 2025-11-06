package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// newCmd defines the 'shiroha new <project-name>' command
var newCmd = &cobra.Command{
	Use:     "new [project-name]",
	Aliases: []string{"generate", "init"},
	Short:   "Create a new Go project structure",
	Long:    "Create a new Go project structure using the specified layered architecture template.",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]

		fmt.Printf("Creating project: %s...\n", projectName)

		if err := createProjectStructure(projectName); err != nil {
			return err
		}

		fmt.Printf("✅ Project '%s' generated successfully!\n", projectName)
		fmt.Println("\nNext steps:")
		fmt.Printf("cd %s\n", projectName)
		fmt.Println("go mod tidy")
		fmt.Println("go run cmd/server.go")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\nDo you want to enter the directory, run 'go mod tidy' and start the project? (y/n): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			projectDir := filepath.Join(".", projectName)

			fmt.Println("\nExecuting 'go mod tidy'...")
			tidyCmd := exec.Command("go", "mod", "tidy")
			tidyCmd.Dir = projectDir
			tidyCmd.Stdout = os.Stdout
			tidyCmd.Stderr = os.Stderr
			if err := tidyCmd.Run(); err != nil {
				return fmt.Errorf("failed to run 'go mod tidy': %w", err)
			}
			fmt.Println("✅ 'go mod tidy' completed successfully")

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

// createProjectStructure creates the folder structure and starter files
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

	// 获取当前 Go 版本
	goVersion := runtime.Version() // e.g., "go1.21.3"
	goVersion = strings.TrimPrefix(goVersion, "go")

	// 只保留主版本号和次版本号，例如 1.21.3 -> 1.21
	parts := strings.Split(goVersion, ".")
	if len(parts) >= 2 {
		goVersion = parts[0] + "." + parts[1]
	} else {
		goVersion = "1.18" // fallback
	}

	// 确保最低为 1.18
	if compareGoVersion(goVersion, "1.18") < 0 {
		goVersion = "1.18"
	}

	// 生成 go.mod
	goModContent := fmt.Sprintf(`module %s

go %s

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/spf13/viper v1.18.2
)`, name, goVersion)
	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	// 其余内容保持不变 ↓
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
	fmt.Printf("Server running at http://localhost:%%d\n", port)
	if err := r.Run(fmt.Sprintf(":%%d", port)); err != nil {
		log.Printf("server start error: %%v", err)
	}
}`, name, name)
	if err := os.WriteFile(filepath.Join(name, "cmd", "server.go"), []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to write cmd/server.go: %w", err)
	}

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

	readmeContent := fmt.Sprintf(`# %s

This project uses Go language with a layered architecture.

## Structure
- internal/handler: Interface layer (Controller)
- internal/service: Business logic layer
- internal/repository: Data access layer

## Startup
1. cd %s
2. go mod tidy
3. go run cmd/server.go
`, name, name)
	if err := os.WriteFile(filepath.Join(name, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	configYamlContent := `server:
  port: 8080`
	if err := os.WriteFile(filepath.Join(name, "config.yaml"), []byte(configYamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

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

// compareGoVersion compares two Go version strings (e.g. "1.18" < "1.21" → -1)
func compareGoVersion(v1, v2 string) int {
	parse := func(v string) (int, int) {
		parts := strings.Split(v, ".")
		major := 0
		minor := 0
		if len(parts) >= 1 {
			fmt.Sscanf(parts[0], "%d", &major)
		}
		if len(parts) >= 2 {
			fmt.Sscanf(parts[1], "%d", &minor)
		}
		return major, minor
	}

	m1, n1 := parse(v1)
	m2, n2 := parse(v2)

	if m1 != m2 {
		if m1 < m2 {
			return -1
		}
		return 1
	}
	if n1 < n2 {
		return -1
	} else if n1 > n2 {
		return 1
	}
	return 0
}
