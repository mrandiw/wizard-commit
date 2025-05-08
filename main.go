package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mrandiw/wizard-commit/cmd"
)

func main() {
	// Load configuration
	config := cmd.LoadConfig()

	// Define flags with defaults from config
	autoCommit := flag.Bool("a", false, "Automatically commit using the generated message")
	model := flag.String("model", config.DefaultModel, "AI model to use")
	noConfirm := flag.Bool("y", false, "Skip confirmation prompt")
	saveConfig := flag.Bool("save-config", false, "Save current settings to config file")
	ollamaURL := flag.String("url", config.OllamaAPIURL, "Ollama API URL")
	provider := flag.String("provider", config.Provider, "API type to use (ollama, gemini, or deepseek)")
	geminiKey := flag.String("gemini-key", config.GeminiAPIKey, "Gemini API key")
	geminiURL := flag.String("gemini-url", config.GeminiAPIURL, "Gemini API URL")
	deepseekKey := flag.String("deepseek-key", config.DeepSeekAPIKey, "DeepSeek API key")
	deepseekURL := flag.String("deepseek-url", config.DeepSeekAPIURL, "DeepSeek API URL")
	flag.Parse()

	// Save configuration if requested
	if *saveConfig {
		config.DefaultModel = *model
		config.OllamaAPIURL = *ollamaURL
		config.Provider = *provider
		config.GeminiAPIKey = *geminiKey
		config.GeminiAPIURL = *geminiURL
		config.DeepSeekAPIKey = *deepseekKey
		config.DeepSeekAPIURL = *deepseekURL

		// Convert config to JSON
		configJSON, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config JSON: %v\n", err)
			os.Exit(1)
		}

		// Write to home directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}

		configPath := filepath.Join(homeDir, ".wizard-commit.json")
		if err := os.WriteFile(configPath, configJSON, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration saved to %s\n", configPath)
		os.Exit(0)
	}

	// Update config with command line arguments for the current run
	config.Provider = *provider
	config.GeminiAPIKey = *geminiKey
	config.GeminiAPIURL = *geminiURL
	config.DeepSeekAPIKey = *deepseekKey
	config.DeepSeekAPIURL = *deepseekURL

	// Validate API configuration
	if config.Provider == "gemini" && config.GeminiAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: Gemini API requires an API key. Use --gemini-key or save it in the configuration.\n")
		os.Exit(1)
	}

	if config.Provider == "deepseek" && config.DeepSeekAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: DeepSeek API requires an API key. Use --deepseek-key or save it in the configuration.\n")
		os.Exit(1)
	}

	// Get git diff
	gitDiff, err := cmd.GetGitDiff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting git diff: %v\n", err)
		os.Exit(1)
	}

	if gitDiff == "" {
		fmt.Println("No changes to commit")
		os.Exit(0)
	}

	// Show which API is being used
	fmt.Printf("Using %s API with model %s\n", config.Provider, *model)

	// Generate commit message using selected API
	commitMsg, err := cmd.GenerateCommitMessage(gitDiff, *model, *ollamaURL, config.PromptTemplate, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating commit message: %v\n", err)
		os.Exit(1)
	}

	// Print the generated commit message
	fmt.Println("Generated commit message:")
	fmt.Println("------------------------")
	fmt.Println(commitMsg)
	fmt.Println("------------------------")

	// If auto-commit flag is set
	if *autoCommit {
		// Skip confirmation if -y flag is provided
		if !*noConfirm {
			confirmed := cmd.ConfirmCommit(commitMsg)
			if !confirmed {
				fmt.Println("Commit aborted.")
				os.Exit(0)
			}
		}

		err = cmd.ExecuteGitCommit(commitMsg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error executing git commit: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Changes committed successfully!")
	} else {
		fmt.Println("Use -a flag to automatically commit with this message")
	}
}
