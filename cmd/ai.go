package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OllamaRequest represents a request to the Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// OllamaResponse represents a response from the Ollama API
// The Ollama API might return the response in different formats
// We'll handle multiple possible response structures
type OllamaResponse struct {
	Response string `json:"response"`
	Content  string `json:"content"` // Some versions use content instead of response
}

// GeminiRequest represents a request to the Gemini API
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

// GeminiContent represents the content portion of a Gemini request
type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

// GeminiPart represents a content part in a Gemini request
type GeminiPart struct {
	Text string `json:"text"`
}

// GeminiResponse represents a response from the Gemini API
type GeminiResponse struct {
	Candidates []GeminiCandidate `json:"candidates"`
}

// GeminiCandidate represents a candidate response from Gemini
type GeminiCandidate struct {
	Content GeminiContent `json:"content"`
}

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

// generateWithOllama generates a commit message using the Ollama API
func generateWithOllama(prompt, model, apiURL string) (string, error) {
	// Prepare request to Ollama API
	ollamaReq := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false, // We want the complete response, not streamed
	}

	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	// Send request to Ollama API
	resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to call Ollama API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the full response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse response
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(bodyBytes, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}

	// Check which field has the content
	var commitMsg string
	if ollamaResp.Response != "" {
		commitMsg = strings.TrimSpace(ollamaResp.Response)
	} else if ollamaResp.Content != "" {
		commitMsg = strings.TrimSpace(ollamaResp.Content)
	} else {
		// Try to find any relevant text in the response
		if strings.Contains(string(bodyBytes), "response") || strings.Contains(string(bodyBytes), "content") {
			// Try to extract the value manually
			for _, line := range strings.Split(string(bodyBytes), ",") {
				if strings.Contains(line, "response") || strings.Contains(line, "content") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) > 1 {
						commitMsg = strings.TrimSpace(parts[1])
						// Remove quotes
						commitMsg = strings.Trim(commitMsg, "\"' ")
						break
					}
				}
			}
		}

		// If still empty, use the entire response as a fallback
		if commitMsg == "" {
			commitMsg = strings.TrimSpace(string(bodyBytes))
		}
	}

	// Remove quotes if they're wrapping the message
	if (strings.HasPrefix(commitMsg, "\"") && strings.HasSuffix(commitMsg, "\"")) ||
		(strings.HasPrefix(commitMsg, "'") && strings.HasSuffix(commitMsg, "'")) {
		commitMsg = commitMsg[1 : len(commitMsg)-1]
	}

	return commitMsg, nil
}

// generateWithGemini generates a commit message using the Gemini API
func generateWithGemini(prompt, model, apiURL, apiKey string) (string, error) {
	// Prepare Gemini request
	geminiReq := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{
						Text: prompt,
					},
				},
			},
		},
	}

	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal Gemini request: %v", err)
	}

	// Add API key to URL
	if !strings.Contains(apiURL, "key=") {
		if strings.Contains(apiURL, "?") {
			apiURL += "&key=" + apiKey
		} else {
			apiURL += "?key=" + apiKey
		}
	}

	// Send request to Gemini API
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create Gemini request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Gemini API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Gemini API returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read Gemini response body: %v", err)
	}

	// Parse Gemini response
	var geminiResp GeminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return "", fmt.Errorf("failed to parse Gemini response: %v", err)
	}

	// Extract commit message from response
	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from Gemini API")
	}

	// Extract text from the first candidate's first part
	var commitMsg string
	commitMsg = geminiResp.Candidates[0].Content.Parts[0].Text

	if commitMsg == "" {
		return "", fmt.Errorf("could not extract commit message from Gemini response")
	}

	return strings.TrimSpace(commitMsg), nil
}
