package ai

import (
	"context"
	"fmt"

	"tutorgpt/pkg/models"
)

type LectureBuilder struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewLectureBuilder(openAIClient *OpenAIClient, prompts *Prompts) *LectureBuilder {
	return &LectureBuilder{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (lb *LectureBuilder) GenerateLecture(ctx context.Context, sessionHistory []models.Message, topic string) (string, error) {

	prompt := fmt.Sprintf(lb.prompts.LecturePrompt, topic)

	response, err := lb.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating lecture: %v", err)
	}

	return response, nil
}