package ai

import (
	"context"
	"fmt"

	"tutorgpt/pkg/models"
)

type TaskManager struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewTaskManager(openAIClient *OpenAIClient, prompts *Prompts) *TaskManager {
	return &TaskManager{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (tm *TaskManager) GenerateTask(ctx context.Context, sessionHistory []models.Message, topic string) (string, error) {
	prompt := fmt.Sprintf(tm.prompts.TaskGenerationPrompt, topic)

	response, err := tm.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating task: %v", err)
	}

	return response, nil
}

func (tm *TaskManager) CheckSubmission(ctx context.Context, sessionHistory []models.Message, instructions string) (string, error) {
	prompt := fmt.Sprintf(tm.prompts.TaskSubmissionPrompt, instructions)

	response, err := tm.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error checking submission: %v", err)
	}

	return response, nil
}
