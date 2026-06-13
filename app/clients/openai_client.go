package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const openAIEndpoint = "https://api.openai.com/v1/chat/completions"

type OpenAIClient struct {
	apiKey   string
	model    string
	endpoint string
	http     *http.Client
}

func NewOpenAIClient(apiKey, model string, timeoutSecs int) *OpenAIClient {
	return &OpenAIClient{
		apiKey:   apiKey,
		model:    model,
		endpoint: openAIEndpoint,
		http:     &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
	}
}

type openAIRequest struct {
	Model          string          `json:"model"`
	Messages       []openAIMessage `json:"messages"`
	Temperature    float64         `json:"temperature"`
	ResponseFormat responseFormat  `json:"response_format"`
}

type responseFormat struct {
	Type string `json:"type"`
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIResponse struct {
	Choices []struct {
		Message openAIMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// ChatComplete sends a system + user prompt and returns the assistant's text.
func (c *OpenAIClient) ChatComplete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	body, err := json.Marshal(openAIRequest{
		Model: c.model,
		Messages: []openAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.3,
		// Força o modelo a devolver um objeto JSON válido (sem texto/markdown
		// em volta). A palavra "JSON" precisa aparecer nos prompts — já aparece.
		ResponseFormat: responseFormat{Type: "json_object"},
	})
	if err != nil {
		return "", fmt.Errorf("marshal openai request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create openai request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("openai request: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read openai response: %w", err)
	}

	var result openAIResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("unmarshal openai response: %w", err)
	}
	if result.Error != nil {
		return "", fmt.Errorf("openai api error: %s", result.Error.Message)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("openai returned no choices")
	}
	return result.Choices[0].Message.Content, nil
}
