package ai

import (
	"context"
	"fmt"

	"tutorgpt/pkg/models"
)

type VagueExplainer struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewVagueExplainer(openAIClient *OpenAIClient, prompts *Prompts) *VagueExplainer {
	return &VagueExplainer{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (ve *VagueExplainer) GenerateVagueExplanation(ctx context.Context, sessionHistory []models.Message, topic string) (string, error) {
	prompt := fmt.Sprintf(ve.prompts.VagueExplainerPrompt, topic)

	response, err := ve.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating vague explanation: %v", err)
	}

	return response, nil
}
