# Golang Development Environment Quick Start Guide

A comprehensive guide for setting up Go development environments on macOS and Windows workstations.

## Table of Contents

- [Introduction to Go](#introduction-to-go)
- [Installation](#installation)
  - [macOS Installation](#macos-installation)
  - [Windows Installation](#windows-installation)
- [Environment Setup](#environment-setup)
- [Essential Go Concepts](#essential-go-concepts)
- [Development Tools](#development-tools)
- [First Steps](#first-steps)
- [Resources and Documentation](#resources-and-documentation)

## Introduction to Go

Go (also known as Golang) is an open-source programming language developed by Google in 2007. It was designed to be simple, efficient, and reliable for building scalable software systems.

### Key Features of Go:
- **Compiled Language**: Go compiles to native machine code, resulting in fast execution
- **Garbage Collected**: Automatic memory management reduces memory leaks
- **Concurrent Programming**: Built-in support for goroutines and channels
- **Static Typing**: Type safety with compile-time error checking
- **Cross-Platform**: Single codebase can compile for multiple operating systems
- **Fast Compilation**: Rapid build times even for large projects
- **Simple Syntax**: Clean, readable code with minimal boilerplate

### Common Use Cases:
- **Web Services & APIs**: HTTP servers, REST APIs, microservices
- **Cloud & Network Programming**: Distributed systems, containerization tools
- **DevOps Tools**: CLI utilities, automation scripts, deployment tools
- **System Programming**: Operating system utilities, database systems
- **Concurrent Applications**: Real-time systems, data processing pipelines

## Installation

### macOS Installation

#### Prerequisites
Ensure you have Homebrew installed. If not, install it first:
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

#### Install Go using Homebrew
```bash
# Update Homebrew to latest version
brew update

# Install the latest stable version of Go
brew install go

# Verify installation
go version
```

#### Alternative: Install Specific Go Version
```bash
# Install a specific version (if needed)
brew install go@1.21

# Link the specific version
brew link go@1.21
```

### Windows Installation

#### Prerequisites
Ensure you have Chocolatey installed. If not, install it first by running PowerShell as Administrator:
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

#### Install Go using Chocolatey
```powershell
# Open PowerShell as Administrator and run:
choco install golang

# Verify installation
go version
```

#### Alternative: Manual Installation
1. Download Go from the [official website](https://golang.org/dl/)
2. Run the MSI installer
3. Follow the installation wizard
4. Verify installation in Command Prompt: `go version`

## Environment Setup

### Configure GOPATH and GOROOT (Optional for Go 1.11+)

Starting with Go 1.11, Go modules are the recommended way to manage dependencies, making GOPATH optional. However, understanding these environment variables is still important:

#### macOS/Linux Environment Variables
Add to your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.profile`):
```bash
# Go environment variables (optional for Go 1.11+)
export GOROOT="/usr/local/go"  # Go installation directory
export GOPATH="$HOME/go"       # Go workspace directory
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"
```

Reload your shell configuration:
```bash
source ~/.zshrc  # or ~/.bashrc
```

#### Windows Environment Variables
```powershell
# Set environment variables (run in PowerShell as Administrator)
[Environment]::SetEnvironmentVariable("GOPATH", "$env:USERPROFILE\go", "User")
[Environment]::SetEnvironmentVariable("PATH", "$env:PATH;$env:GOPATH\bin", "User")
```

### Verify Environment Setup
```bash
# Check Go installation
go version

# Check Go environment
go env

# Key environment variables to verify
go env GOROOT
go env GOPATH
go env GOMODCACHE
```

## Essential Go Concepts

### Go Modules (Dependency Management)
Go modules are the modern way to manage dependencies in Go projects:

```bash
# Initialize a new module
go mod init example.com/myproject

# Add dependencies
go get github.com/gin-gonic/gin

# Download dependencies
go mod download

# Clean up unused dependencies
go mod tidy

# Vendor dependencies (optional)
go mod vendor
```

### Project Structure
Recommended Go project structure:
```
myproject/
├── go.mod              # Module definition
├── go.sum              # Dependency checksums
├── main.go             # Main application entry point
├── cmd/                # Application entry points
│   └── myapp/
│       └── main.go
├── internal/           # Private application code
│   ├── handlers/
│   ├── models/
│   └── services/
├── pkg/                # Public library code
│   └── utils/
├── api/                # API definitions (OpenAPI, gRPC)
├── web/                # Web application assets
├── configs/            # Configuration files
├── scripts/            # Build and deployment scripts
├── test/               # Integration tests
└── docs/               # Documentation
```

### Go Workspace Commands
```bash
# Build the current package
go build

# Build and install
go install

# Run the application
go run main.go

# Test the package
go test

# Format code
go fmt ./...

# Lint code (requires golint)
golint ./...

# Generate documentation
go doc

# Clean build cache
go clean

# Check for security vulnerabilities
go list -json -deps ./... | nancy sleuth
```

## Development Tools

### Essential VS Code Extensions
If using Visual Studio Code, install these extensions:
```bash
# Install VS Code Go extension
code --install-extension golang.go
```

### Recommended CLI Tools
```bash
# macOS installation using Homebrew
brew install golangci-lint    # Advanced Go linter
brew install delve            # Go debugger
brew install gore             # Go REPL

# Windows installation using Chocolatey
choco install golangci-lint
choco install delve

# Install using Go
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/motemen/gore/cmd/gore@latest
```

### IDE Options
- **VS Code**: Free, excellent Go support with official extension
- **GoLand**: JetBrains IDE specifically for Go (paid)
- **Vim/Neovim**: With vim-go plugin
- **Emacs**: With go-mode
- **Atom**: With go-plus package

## First Steps

### Create Your First Go Program

1. **Create a new directory and initialize a module:**
```bash
mkdir hello-go
cd hello-go
go mod init hello-go
```

2. **Create main.go:**
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```

3. **Run the program:**
```bash
go run main.go
```

4. **Build an executable:**
```bash
go build -o hello
./hello  # macOS/Linux
hello.exe  # Windows
```

### Cross-Compilation Example
Go makes it easy to build for different platforms:
```bash
# Build for Windows from macOS/Linux
GOOS=windows GOARCH=amd64 go build -o hello.exe

# Build for macOS from Windows/Linux
GOOS=darwin GOARCH=amd64 go build -o hello-mac

# Build for Linux from Windows/macOS
GOOS=linux GOARCH=amd64 go build -o hello-linux

# See all supported platforms
go tool dist list
```

### Testing Your Setup
Create a simple test file (`main_test.go`):
```go
package main

import "testing"

func TestHello(t *testing.T) {
    expected := "Hello, Go!"
    if got := hello(); got != expected {
        t.Errorf("hello() = %q, want %q", got, expected)
    }
}

func hello() string {
    return "Hello, Go!"
}
```

Run the test:
```bash
go test
```

## Resources and Documentation

### Official Resources
- **Go Official Website**: https://golang.org/
- **Go Documentation**: https://golang.org/doc/
- **Go Tour (Interactive Tutorial)**: https://tour.golang.org/
- **Go Playground**: https://play.golang.org/
- **Go Package Repository**: https://pkg.go.dev/

### Learning Resources
- **Effective Go**: https://golang.org/doc/effective_go.html
- **Go by Example**: https://gobyexample.com/
- **Go Web Examples**: https://gowebexamples.com/
- **Awesome Go**: https://github.com/avelino/awesome-go

### Community and Support
- **Go Forum**: https://forum.golangbridge.org/
- **Go Reddit**: https://www.reddit.com/r/golang/
- **Gopher Slack**: https://gophers.slack.com/
- **Go Discord**: https://discord.gg/golang

### Best Practices and Style Guides
- **Go Code Review Comments**: https://github.com/golang/go/wiki/CodeReviewComments
- **Uber Go Style Guide**: https://github.com/uber-go/guide
- **Google Go Style Guide**: https://google.github.io/styleguide/go/

### Package Management
- **Go Modules Reference**: https://golang.org/ref/mod
- **Go Module Proxy**: https://proxy.golang.org/
- **Popular Go Packages**: https://pkg.go.dev/search?q=most+imported

## Troubleshooting

### Common Issues and Solutions

**Issue**: `go: command not found`
**Solution**: Ensure Go is installed and added to your PATH environment variable.

**Issue**: Module-related errors
**Solution**: Ensure you're in a directory with a `go.mod` file or run `go mod init`.

**Issue**: Permission denied when installing packages
**Solution**: Check file permissions and avoid using `sudo` with Go commands.

**Issue**: Proxy errors when downloading modules
**Solution**: 
```bash
# Disable module proxy temporarily
go env -w GOPROXY=direct

# Or configure proxy settings
go env -w GOPROXY=https://proxy.golang.org,direct
```

### Performance Tips
- Use `go build -ldflags="-s -w"` to reduce binary size
- Profile your applications with `go tool pprof`
- Use `go mod vendor` for faster builds in CI/CD
- Enable module caching: `go env -w GOPROXY=https://proxy.golang.org`

## Next Steps

1. **Complete the Go Tour**: https://tour.golang.org/
2. **Read Effective Go**: https://golang.org/doc/effective_go.html
3. **Build a simple CLI tool or web service**
4. **Explore popular Go frameworks** (Gin, Echo, Fiber for web development)
5. **Learn about Go's concurrency patterns** (goroutines, channels)
6. **Contribute to open-source Go projects**

---

*This guide provides the foundation for Go development. As you progress, explore advanced topics like performance optimization, advanced concurrency patterns, and Go's reflection capabilities.*
