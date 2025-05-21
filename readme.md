[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
[![ChatGPT](https://img.shields.io/badge/ChatGPT-74aa9c?logo=openai&logoColor=white)](#)
[![Deepseek](https://custom-icon-badges.demolab.com/badge/Deepseek-4D6BFF?logo=deepseek&logoColor=fff)](#)
[![Google Gemini](https://img.shields.io/badge/Google%20Gemini-886FBF?logo=googlegemini&logoColor=fff)](#)
[![Claude](https://img.shields.io/badge/Claude-D97757?logo=claude&logoColor=fff)](#)
[![LinkedIn](https://custom-icon-badges.demolab.com/badge/LinkedIn-0A66C2?logo=linkedin-white&logoColor=fff)](https://www.linkedin.com/in/mrandiw/)
[![YouTube](https://img.shields.io/badge/YouTube-%23FF0000.svg?logo=YouTube&logoColor=white)](https://www.youtube.com/@CodeWithAndiw)

## Wizard Commit

A flexible CLI tool that uses various AI models (Ollama, Gemini, Deepseek, OpenAI, Claude, Groq) to automatically generate git commit messages based on your changes.

## Features

- Support for multiple AI providers:
  - Ollama (local AI models)
  - Gemini (Google's API)
  - Deepseek (Deepseek API)
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
  - Deepseek API Key

## Installation

### Download from Release Page
https://github.com/mrandiw/wizard-commit/releases/download/v1.0.0/app.zip

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/mrandiw/wizard-commit.git
   cd wizard-commit
   ```

2. Build the executable:
   ```bash
   go build -o wizard-commit . (MacOS / Linux)
   go build -o wizard-commit.exe . (Windows)
   ```

3. Move the executable & Config file to your Project Folder then run:

   **Linux/macOS**:
   ```bash
   ./wizard-commit -a
   ```

   **Windows**:
   ```bash
   ./wizard-commit.exe -a
   ```

### Cross-Compilation

You can build for different platforms from any OS:

```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o wizard-commit .

# Build for macOS 
GOOS=darwin GOARCH=amd64 go build -o wizard-commit .

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
wizard-commit -provider gemini -gemini-key YOUR_API_KEY

# Use Deepseek
wizard-commit -provider deepseek -deepseek-key YOUR_API_KEY -model deepseek-chat
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
wizard-commit -provider gemini -gemini-key YOUR_API_KEY -save-config

# Deepseek
wizard-commit -provider deepseek -deepseek-key YOUR_API_KEY -model deepseek-chat -save-config
```

### Configuration File Format

The configuration file is in JSON format:

```json
{
  "provider" : "gemini",
  "ollamaApiUrl": "http://localhost:11434/api/generate",
  "geminiApiUrl": "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash-lite:generateContent",
  "geminiApiKey": "your-api-key-here",
  "deepseekApiUrl" : "https://api.deepseek.com/chat/completions",
  "deepseekApiKey" : "your-api-key-here",
  "defaultModel": "llama3",
  "promptTemplate": "Act as a software developer.\nGive commit message based on code changes no more than two sentenses. \n\nContex:\n%s"
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

You'll need a Gemini API key from Google. You can provide it either through the `-gemini-key` flag or by storing it in your configuration file.

### Deepseek

You'll need a Deepseek API key from Deepseek. You can provide it either through the `-deepseek-key` flag or by storing it in your configuration file.


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
$ wizard-commit -provider gemini -gemini-key YOUR_API_KEY -a -y
Generated commit message:
------------------------
feat: implement user profile page with avatar upload
------------------------
Changes committed successfully!

# Commit using Deepseek
$ wizard-commit -provider deepseek -deepseek-key YOUR_API_KEY -model deepseek-chat -a -y
Generated commit message:
------------------------
feat: implement user profile page with avatar upload
------------------------
Changes committed successfully!
```

## License

MIT
