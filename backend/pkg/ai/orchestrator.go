package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"tutorgpt/pkg/models"
)

type OrchestratorResult struct {
	Tool         string `json:"tool"`
	Instructions string `json:"instructions"`
}

type Orchestrator struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewOrchestrator(openAIClient *OpenAIClient, prompts *Prompts) *Orchestrator {
	return &Orchestrator{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (o *Orchestrator) DetermineIntent(ctx context.Context, sessionHistory []models.Message, message string) (*OrchestratorResult, error) {
	prompt := o.prompts.OrchestratorPrompt + "\n\nUser's message: " + message

	fmt.Printf("Sending prompt to Orchestrator: %s\n", prompt[:min(10, len(prompt))]+"...")

	response, err := o.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		fmt.Printf("Error from OpenAI in orchestrator: %v\n", err)
		return nil, fmt.Errorf("error determining intent: %v", err)
	}

	fmt.Printf("Raw orchestrator response: %s\n", response)

	jsonStart := strings.Index(response, "{")
	jsonEnd := strings.LastIndex(response, "}")

	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := response[jsonStart : jsonEnd+1]
		fmt.Printf("Extracted JSON: %s\n", jsonStr)

		var result OrchestratorResult
		if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
			fmt.Printf("Error parsing JSON: %v\n", err)
			return &OrchestratorResult{
				Tool:         "lecture_builder",
				Instructions: message,
			}, fmt.Errorf("error parsing orchestrator response: %v", err)
		}

		if result.Tool == "" {
			result.Tool = "lecture_builder"
			fmt.Printf("WARNING: Tool was empty in orchestrator response, defaulting to lecture_builder\n")
		}
		if result.Instructions == "" {
			result.Instructions = message
			fmt.Printf("WARNING: Instructions were empty in orchestrator response, using original message\n")
		}

		fmt.Printf("Orchestrator selected tool: %s\n", result.Tool)
		return &result, nil
	}

	fmt.Printf("Could not find valid JSON in response. Using original message as instructions.\n")
	return &OrchestratorResult{
		Tool:         "lecture_builder",
		Instructions: message,
	}, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (o *Orchestrator) GetWelcomeMessage() string {
	return o.prompts.WelcomeMessage
}
