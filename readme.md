# AI-Commit

A flexible CLI tool that uses various AI models (Ollama, Gemini) to automatically generate git commit messages based on your changes.

## Features

- Support for multiple AI providers:
  - Ollama (local AI models)
  - Gemini (Google's API)
- Customizable prompt templates
- Confirmation before committing
- Configuration file for persistent settings
- Cross-platform (Windows, macOS, Linux)

## Requirements

- Go 1.18+ installed
- Git installed
- At least one of:
  - Ollama running locally (default: http://localhost:11434)
  - Gemini API key

## Installation

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/mrandiw/wizard-commit.git
   cd wizard-commit
   ```

2. Build the executable:
   ```bash
   go build -o wizard-commit .
   ```

3. Move the executable to your PATH:

   **Linux/macOS**:
   ```bash
   sudo mv wizard-commit /usr/local/bin/
   # OR for a user-local installation
   mkdir -p ~/bin
   mv wizard-commit ~/bin/
   # Make sure ~/bin is in your PATH
   ```

   **Windows**:
   ```powershell
   # Create a directory for the executable (if it doesn't exist)
   mkdir -p $env:USERPROFILE\bin

   # Move the executable
   move wizard-commit.exe $env:USERPROFILE\bin\

   # Add to PATH (may need admin PowerShell)
   $currentPath = [Environment]::GetEnvironmentVariable("Path", "User")
   [Environment]::SetEnvironmentVariable("Path", $currentPath + ";$env:USERPROFILE\bin", "User")
   ```

### Cross-Compilation

You can build for different platforms from any OS:

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o wizard-commit-linux .

# Build for macOS 
GOOS=darwin GOARCH=amd64 go build -o wizard-commit-macos .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o wizard-commit.exe .
```

### Using Go Install

You can also install directly using Go:

```bash
go install github.com/mrandiw/wizard-commit@latest
```

## Usage

Basic usage:
```bash
# Show the generated commit message without committing
wizard-commit
```

Automatically commit with the generated message:
```bash
wizard-commit -a
```

Use specific AI provider:
```bash
# Use Ollama
wizard-commit -provider ollama -model llama3

# Use Gemini
wizard-commit -provider gemini -api-key YOUR_API_KEY
```

## Configuration

You can configure wizard-commit using a configuration file. The tool looks for configuration in the following locations:

1. `./wizard-commit.json` (current directory)
2. `~/.wizard-commit.json` (home directory)

You can create a configuration file manually or use the `-save-config` flag to save your current settings:

```bash
# Save your current settings to ~/.wizard-commit.json
wizard-commit -provider ollama -model codellama -save-config

# Or save Gemini configuration
wizard-commit -provider gemini -api-key YOUR_API_KEY -save-config
```

### Configuration File Format

The configuration file is in JSON format:

```json
{
  "provider" : "gemini",
  "ollamaApiUrl": "http://localhost:11434/api/generate",
  "geminiApiUrl": "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent",
  "geminiApiKey": "your-api-key-here",
  "defaultModel": "llama3",
  "promptTemplate": "Generate a concise and descriptive git commit message based on the following changes.\nFollow best practices for git commit messages: use imperative mood, keep it under 50 characters for the first line,\nand add more details in a body if necessary.\n\nRespond ONLY with the commit message, no other text, explanation, or quotes.\nJust the commit message that would be used with 'git commit -m'.\n\nChanges:\n%s"
}
```

Command-line flags will override the configuration file settings.

## Available Flags

- `-a`: Automatically commit using the generated message
- `-y`: Skip confirmation prompt (used with -a)
- `-save-config`: Save current settings as your default configuration
- `-provider string`: AI provider to use (ollama, gemini)
- `-model string`: Model to use with Ollama (default from config or "llama3")
- `-url string`: Ollama API URL (default from config or "http://localhost:11434/api/generate")
- `-api-key string`: API key for Gemini

## Provider-Specific Information

### Ollama

Ollama needs to be running on your machine. By default, the tool connects to http://localhost:11434. Make sure you have the model you want to use already loaded.

### Gemini

You'll need a Gemini API key from Google. You can provide it either through the `-api-key` flag or by storing it in your configuration file.

## Example

```bash
# Show the generated commit message without committing
$ wizard-commit
Generated commit message:
------------------------
feat: add user authentication and password reset functionality
------------------------
Use -a flag to automatically commit with this message

# Commit with confirmation using Ollama
$ wizard-commit -provider ollama -a
Generated commit message:
------------------------
feat: add user authentication and password reset functionality
------------------------
Are you sure you want to use this commit message? (y/n): y
Changes committed successfully!

# Commit using Gemini
$ wizard-commit -provider gemini -api-key YOUR_API_KEY -a -y
Generated commit message:
------------------------
feat: implement user profile page with avatar upload
------------------------
Changes committed successfully!
```

## License

MIT