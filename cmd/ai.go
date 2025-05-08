package cmd

import (
	"fmt"
)

// GenerateCommitMessage generates a commit message using either Ollama, Gemini, or DeepSeek API
func GenerateCommitMessage(gitDiff, model, apiURL, promptTemplate string, config Config) (string, error) {
	// Prepare prompt
	prompt := fmt.Sprintf(promptTemplate, gitDiff)

	// Decide which API to use
	if config.Provider == "gemini" {
		return generateWithGemini(prompt, model, config.GeminiAPIURL, config.GeminiAPIKey)
	} else if config.Provider == "deepseek" {
		return generateWithDeepSeek(prompt, model, config.DeepSeekAPIURL, config.DeepSeekAPIKey)
	}

	// Default to Ollama
	return generateWithOllama(prompt, model, apiURL)
}
