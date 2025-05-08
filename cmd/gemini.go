package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

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
