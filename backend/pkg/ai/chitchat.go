package ai

import (
	"context"
	"fmt"

	"tutorgpt/pkg/models"
)

type ChitchatAgent struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewChitchatAgent(openAIClient *OpenAIClient, prompts *Prompts) *ChitchatAgent {
	return &ChitchatAgent{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (chitchat *ChitchatAgent) Chitchat(ctx context.Context, sessionHistory []models.Message, message string) (string, error) {

	prompt := fmt.Sprintf(chitchat.prompts.ChitchatAgentPrompt, message)

	response, err := chitchat.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error answering question: %v", err)
	}
	return response, nil
}
