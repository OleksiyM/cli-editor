package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Assistant provides AI-powered code assistance
type Assistant struct {
	apiKey string
}

// NewAssistant creates a new AI assistant
func NewAssistant() (*Assistant, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable not set")
	}
	
	return &Assistant{
		apiKey: apiKey,
	}, nil
}

// CompletionRequest represents a request to the OpenAI API
type CompletionRequest struct {
	Model       string  `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
}

// Message represents a message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// CompletionResponse represents a response from the OpenAI API
type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GetCompletion gets a completion from the OpenAI API
func (a *Assistant) GetCompletion(prompt string) (string, error) {
	reqBody := CompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful coding assistant. Provide concise, accurate code advice.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: 0.7,
		MaxTokens:   1000,
	}
	
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.apiKey)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	
	var completionResp CompletionResponse
	if err := json.Unmarshal(body, &completionResp); err != nil {
		return "", err
	}
	
	if len(completionResp.Choices) == 0 {
		return "", fmt.Errorf("no completions returned")
	}
	
	return completionResp.Choices[0].Message.Content, nil
}

// GetCodeExplanation explains the provided code
func (a *Assistant) GetCodeExplanation(code string) (string, error) {
	prompt := fmt.Sprintf("Explain the following code:\n\n%s", code)
	return a.GetCompletion(prompt)
}

// GetCodeSuggestion provides suggestions for improving the code
func (a *Assistant) GetCodeSuggestion(code string) (string, error) {
	prompt := fmt.Sprintf("Suggest improvements for the following code:\n\n%s", code)
	return a.GetCompletion(prompt)
}

// GetDocumentation generates documentation for the code
func (a *Assistant) GetDocumentation(code string) (string, error) {
	prompt := fmt.Sprintf("Generate documentation for the following code:\n\n%s", code)
	return a.GetCompletion(prompt)
}