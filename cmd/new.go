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
		fmt.Println("swag init -g cmd/server.go -o docs")
		fmt.Println("go run cmd/server.go")

		reader := bufio.NewReader(os.Stdin)
		// Updated prompt to include Swagger generation
		fmt.Print("\nDo you want to enter the directory, run 'go mod tidy', generate Swagger docs, and start the project? (Y/n, default Y): ") // 提示语更新
		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(strings.ToLower(input))
		// 默认回车 (input == "") 视为同意，与 "y" 或 "yes" 相同
		if input == "y" || input == "yes" || input == "" {
			projectDir := filepath.Join(".", projectName)

			// 1. Run go mod tidy
			fmt.Println("\nExecuting 'go mod tidy'...")
			tidyCmd := exec.Command("go", "mod", "tidy")
			tidyCmd.Dir = projectDir
			tidyCmd.Stdout = os.Stdout
			tidyCmd.Stderr = os.Stderr
			if err := tidyCmd.Run(); err != nil {
				return fmt.Errorf("failed to run 'go mod tidy': %w", err)
			}
			fmt.Println("✅ 'go mod tidy' completed successfully")

			// 2. Run swag init to generate documentation
			fmt.Println("\nExecuting 'swag init -g cmd/server.go -o docs' to generate API documentation...")
			swagCmd := exec.Command("swag", "init", "-g", "cmd/server.go", "-o", "docs")
			swagCmd.Dir = projectDir
			swagCmd.Stdout = os.Stdout
			swagCmd.Stderr = os.Stderr
			if err := swagCmd.Run(); err != nil {
				// Non-fatal error for swag init, often means `swag` tool isn't installed.
				// We proceed to run the server, but inform the user.
				fmt.Printf("⚠️ Warning: Failed to run 'swag init'. Please ensure the 'swag' tool is installed (go install github.com/swaggo/swag/cmd/swag@latest): %v\n", err)
			} else {
				fmt.Println("✅ Swagger documentation generated successfully")
			}

			// 3. Start the project
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
		"docs", // Added docs folder for swagger output
	}

	for _, dir := range dirs {
		fullPath := filepath.Join(name, dir)
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", fullPath, err)
		}
	}

	// get current go version
	goVersion := runtime.Version() // go1.21.3
	goVersion = strings.TrimPrefix(goVersion, "go")

	parts := strings.Split(goVersion, ".")
	if len(parts) >= 2 {
		goVersion = parts[0] + "." + parts[1]
	} else {
		goVersion = "1.18"
	}

	if compareGoVersion(goVersion, "1.18") < 0 {
		goVersion = "1.18"
	}

	// -------- ✅ go.mod (latest versions) --------
	goModContent := fmt.Sprintf(`module %s

go %s

require (
    github.com/gin-gonic/gin latest
    github.com/spf13/viper latest
    github.com/swaggo/files latest
    github.com/swaggo/gin-swagger latest
    github.com/swaggo/swag latest
)
`, name, goVersion)

	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goModContent), 0644); err != nil {
		return fmt.Errorf("failed to write go.mod: %w", err)
	}

	serverContent := fmt.Sprintf(`package main

import (
    _ "%s/docs" // swagger docs

    "fmt"
    "log"

    "%s/config"
    "%s/internal/handler" // 导入新的 handler 包
    "%s/internal/router"

)

// @title Shiroha API
// @version 1.0
// @description Auto-generated API documentation by Shiroha CLI.
// @host localhost:8080
// @BasePath /

func init() {
    config.LoadConfig()
}

func main() {
    r := router.InitRouter()

    // 路由调用 internal/handler/test_handler.go 中的函数
    r.GET("/test", handler.TestEndpoint)

    port := config.Cfg.Server.Port
    fmt.Printf("Server running at http://localhost:%%d\n", port)
    if err := r.Run(fmt.Sprintf(":%%d", port)); err != nil {
       log.Printf("server start error: %%v", err)
    }
}
`, name, name, name, name)

	if err := os.WriteFile(filepath.Join(name, "cmd", "server.go"), []byte(serverContent), 0644); err != nil {
		return fmt.Errorf("failed to write cmd/server.go: %w", err)
	}

	// -------- ✅ test_handler.go (NEW FILE) --------
	testHandlerContent := `package handler

import (
    "github.com/gin-gonic/gin"
)

// TestEndpoint
// @Summary Test endpoint
// @Description Returns a simple greeting message to verify server status.
// @Tags Health
// @Accept  json
// @Produce json
// @Success 200 {object} map[string]interface{} "Returns {message: 'Hello Shiroha'}"
// @Router /test [get]
func TestEndpoint(c *gin.Context) {
    c.JSON(200, gin.H{
       "message": "Hello Shiroha",
    })
}
`

	if err := os.WriteFile(filepath.Join(name, "internal", "handler", "test_handler.go"), []byte(testHandlerContent), 0644); err != nil {
		return fmt.Errorf("failed to write internal/handler/test_handler.go: %w", err)
	}

	routerContent := `package router

import (
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"

    "github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
    r := gin.Default()

    // swagger UI
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}
`

	if err := os.WriteFile(filepath.Join(name, "internal", "router", "main_router.go"), []byte(routerContent), 0644); err != nil {
		return fmt.Errorf("failed to write internal/router/main_router.go: %w", err)
	}

	// -------- ✅ README --------
	readmeContent := fmt.Sprintf(`# %s

This project uses Go with a clean layered architecture.

## Swagger Docs

Generate:
    swag init -g cmd/server.go -o docs

Visit:
    http://localhost:8080/swagger/index.html

`, name)

	if err := os.WriteFile(filepath.Join(name, "README.md"), []byte(readmeContent), 0644); err != nil {
		return fmt.Errorf("failed to write README.md: %w", err)
	}

	// -------- ✅ config.yaml --------
	configYamlContent := `server:
  port: 8080`

	if err := os.WriteFile(filepath.Join(name, "config.yaml"), []byte(configYamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config.yaml: %w", err)
	}

	// -------- ✅ config.go --------
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
       log.Fatalf("Failed to read config file: %v", err)
    }

    Cfg = &Config{}
    if err := v.Unmarshal(Cfg); err != nil {
       log.Fatalf("Failed to unmarshal config: %v", err)
    }
}
`

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
