# Shiroha CLI

A powerful command-line interface (CLI) tool for bootstrapping and managing Gin-based web APIs.

Shiroha is designed to accelerate Go web development by automatically setting up a clean, layered (e.g., domain, service, repository) project architecture using the popular Gin framework. It ensures maintainability, testability, and scalability for your RESTful APIs right from the start.

## ✨ Features

- **Layered Architecture**: Automatically generates a robust project structure (handler, service, repository, model, etc.) following best practices.
- **Rapid Scaffolding**: Create a complete, runnable project with basic routing and configuration in a single command.
- **Dynamic Go Version Support**: Detects the user's local Go environment version and uses it for go.mod, ensuring compatibility, with a minimum requirement of Go 1.18.
- **Cross-Platform Build**: Easily compile your project for deployment on Linux, macOS, and Windows with a simple interactive command.

## 🚀 Installation

To install shiroha on your system, you need to have Go (1.18 or higher) installed and configured.

### Via Go Install

If you use go to install this tool, you should make sure that you have already added go env to your environment.

```bash
go install github.com/NexusXian/shiroha@latest
```

> Note: The CLI tool is named shiroha. Ensure your $GOPATH/bin or $GOBIN is in your system's $PATH environment variable to run the command globally.

### Via AUR (Arch-based Linux)

For Arch-based Linux distributions (including Manjaro, ArcoLinux, etc.), you can easily install shiroha-cli from the Arch User Repository (AUR). Use your preferred AUR helper (such as yay or paru):

```bash
yay -S shiroha-cli
# OR
paru -S shiroha-cli
```

## 💻 Usage

### 1. Creating a New Project (shiroha new)

Use the new command to generate a new project structure. Replace `<project-name>` with your desired module name (e.g., my-awesome-api).

```bash
shiroha new <project-name>
```

The command will:
- Create the project directory and file structure.
- Generate essential files (main.go, config.yaml, README.md, etc.).
- Crucially, it will set the go.mod file's Go version based on your installed Go environment.
- Optionally prompt you to run `go mod tidy` and start the server immediately.

### 2. Running the Project (shiroha run)

Use the run command to start the application server directly from the project root directory.

```bash
shiroha run
```

This command executes `go run cmd/server.go`.

### 3. Building for Deployment (shiroha build)

Use the build command to compile your project into a standalone executable for deployment targets.

```bash
shiroha build
```
### 4. Generate interface doc (shiroha doc)

Use the doc command to generate your project interface document.

```bash
shiroha doc
```

## 📂 Project Structure & Directory Functions

The project uses a standard layered architecture for maintainability and scalability.

| Directory               | Layer                  | Function                                                                 |
|-------------------------|------------------------|--------------------------------------------------------------------------|
| cmd                     | Entrypoint             | Contains application entrypoints, such as server.go for starting the web service. |
| internal/router         | Presentation           | Initialization and configuration of the Gin router and definition of all HTTP routes. |
| internal/handler        | Interface (Controller) | Handles HTTP requests, validates input from request DTOs, calls the service layer, and formats output using response DTOs. |
| internal/service        | Business Logic         | Contains the core business rules and transactional logic, coordinating between the handler and repository layers. |
| internal/repository     | Data Access            | Handles communication with the persistence layer (e.g., database), translating model entities to/from the data store. |
| internal/model          | Domain/Data Model      | Defines the core data structures (structs) used throughout the application, often representing database tables. |
| internal/request        | Data Transfer Object (DTO) | Input structs used to bind and validate incoming request data (e.g., JSON body, query params). |
| internal/response       | Data Transfer Object (DTO) | Output structs used to format and serialize data sent back to the client. |
| internal/database       | Database Config        | Holds database connection setup and migration logic.                     |
| config                  | Configuration          | Contains configuration loading logic (config.go) using Viper, reading settings from config.yaml. |
| pkg/utils               | Shared Utilities       | Common, reusable non-domain-specific functions (e.g., formatting, string manipulation). |
| pkg/jwt                 | Shared Utilities       | Logic for JWT token generation, parsing, and authentication.             |
| sql                     | Database               | Stores database migration scripts or schema definitions.                 |
| bin                     | Build Output           | Destination for the compiled executable binary files after running `shiroha build`. |

## ⚙️ Go Version Support Policy

Shiroha is designed to be compatible with modern Go development practices.

| Component               | Policy                                                                 |
|-------------------------|------------------------------------------------------------------------|
| Minimum Required Version | Go 1.18                                                                |
| go.mod Generation       | The CLI uses the currently installed Go version (from runtime.Version()) to set the go directive in the generated go.mod file. |
| Fallback                | If the detected Go version is older than 1.18, the CLI will automatically enforce the minimum version (go 1.18) to ensure all dependencies and modern features are available. |

---
