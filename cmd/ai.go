package cmd

import (
	"fmt"
)

// GenerateCommitMessage generates a commit message using either Ollama or Gemini API
func GenerateCommitMessage(gitDiff, model, apiURL, promptTemplate string, config Config) (string, error) {
	// Prepare prompt
	prompt := fmt.Sprintf(promptTemplate, gitDiff)

	// Decide which API to use
	if config.Provider == "gemini" {
		return generateWithGemini(prompt, model, config.GeminiAPIURL, config.GeminiAPIKey)
	}

	// Default to Ollama
	return generateWithOllama(prompt, model, apiURL)
}
