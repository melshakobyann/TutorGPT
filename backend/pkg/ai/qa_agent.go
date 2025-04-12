package ai

import (
	"context"
	"fmt"

	"tutorgpt/pkg/models"
)

type QAAgent struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewQAAgent(openAIClient *OpenAIClient, prompts *Prompts) *QAAgent {
	return &QAAgent{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (qa *QAAgent) AnswerQuestion(ctx context.Context, sessionHistory []models.Message, question string) (string, error) {

	prompt := fmt.Sprintf(qa.prompts.QAAgentPrompt, question)

	response, err := qa.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error answering question: %v", err)
	}
	return response, nil
}
