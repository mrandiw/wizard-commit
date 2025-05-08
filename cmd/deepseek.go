package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// DeepSeekRequest represents a request to the DeepSeek API
type DeepSeekRequest struct {
	Model    string            `json:"model"`
	Messages []DeepSeekMessage `json:"messages"`
	Stream   bool              `json:"stream"`
}

// DeepSeekMessage represents a message in the DeepSeek API request
type DeepSeekMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// DeepSeekResponse represents a response from the DeepSeek API
type DeepSeekResponse struct {
	Choices []DeepSeekChoice `json:"choices"`
}

// DeepSeekChoice represents a choice in the DeepSeek API response
type DeepSeekChoice struct {
	Message DeepSeekMessage `json:"message"`
}

// generateWithDeepSeek generates a commit message using the DeepSeek API
func generateWithDeepSeek(prompt, model, apiURL, apiKey string) (string, error) {
	// Prepare DeepSeek request
	deepseekReq := DeepSeekRequest{
		Model: model,
		Messages: []DeepSeekMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false, // We want the complete response, not streamed
	}

	reqBody, err := json.Marshal(deepseekReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DeepSeek request: %v", err)
	}

	// Send request to DeepSeek API
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create DeepSeek request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call DeepSeek API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("DeepSeek API returned non-OK status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read DeepSeek response body: %v", err)
	}

	// Parse DeepSeek response
	var deepseekResp DeepSeekResponse
	if err := json.Unmarshal(bodyBytes, &deepseekResp); err != nil {
		return "", fmt.Errorf("failed to parse DeepSeek response: %v", err)
	}

	// Extract commit message from response
	if len(deepseekResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from DeepSeek API")
	}

	// Extract text from the first choice's message content
	commitMsg := deepseekResp.Choices[0].Message.Content

	if commitMsg == "" {
		return "", fmt.Errorf("could not extract commit message from DeepSeek response")
	}

	return strings.TrimSpace(commitMsg), nil
}
