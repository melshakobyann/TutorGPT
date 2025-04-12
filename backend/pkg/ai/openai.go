package ai

import (
	"context"
	"fmt"
	"time"

	"tutorgpt/pkg/models"

	"github.com/sashabaranov/go-openai"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	fmt.Println("Running with OpenAI API")

	if len(apiKey) > 10 {
		apiKeyPrefix := apiKey[:5]
		apiKeySuffix := apiKey[len(apiKey)-5:]
		fmt.Printf("Using API Key: %s...%s\n", apiKeyPrefix, apiKeySuffix)
	}

	client := openai.NewClient(apiKey)

	return &OpenAIClient{
		client: client,
	}
}

func (c *OpenAIClient) GenerateResponse(ctx context.Context, sessionHistory []models.Message, prompt string) (string, error) {
	messages := make([]openai.ChatCompletionMessage, 0, len(sessionHistory)+1)

	for _, msg := range sessionHistory {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    "user",
		Content: prompt,
	})

	modelToUse := "o3-mini-2025-01-31"
	fmt.Printf("Using model: %s\n", modelToUse)

	request := openai.ChatCompletionRequest{
		Model:               modelToUse,
		Messages:            messages,
		MaxCompletionTokens: 100000,
	}

	// Create a channel to handle timeouts more gracefully
	resultChan := make(chan struct {
		response *openai.ChatCompletionResponse
		err      error
	})

	// Set a longer timeout for the API call (2 minutes)
	timeoutCtx, cancel := context.WithTimeout(ctx, 500*time.Second)
	defer cancel()

	// Make the API call in a goroutine
	go func() {
		fmt.Printf("Making OpenAI API request with model: %s\n", modelToUse)
		resp, err := c.client.CreateChatCompletion(timeoutCtx, request)
		resultChan <- struct {
			response *openai.ChatCompletionResponse
			err      error
		}{response: &resp, err: err}
	}()

	// Wait for either a result or timeout
	var response *openai.ChatCompletionResponse
	var err error

	select {
	case result := <-resultChan:
		response = result.response
		err = result.err
	case <-timeoutCtx.Done():
		fmt.Println("OpenAI API request timed out, attempting fallback")
		err = fmt.Errorf("request timed out after 120 seconds")
	}

	if err != nil {
		errorMsg := fmt.Sprintf("Unable to generate a response. Primary error: %v", err)
		return "I apologize, but I'm currently having trouble. Please try again later.", fmt.Errorf(errorMsg)
	}

	if response != nil && len(response.Choices) > 0 {
		return response.Choices[0].Message.Content, nil
	}

	return "I apologize, but I couldn't generate a response. Please try again with a different query.", fmt.Errorf("no response from OpenAI API")
}
