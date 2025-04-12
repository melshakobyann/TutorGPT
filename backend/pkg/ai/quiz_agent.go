package ai

import (
	"context"
	"fmt"
	"strings"

	"tutorgpt/pkg/models"
)

type QuizGenerator struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewQuizGenerator(openAIClient *OpenAIClient, prompts *Prompts) *QuizGenerator {
	return &QuizGenerator{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (vg *QuizGenerator) GenerateQuiz(ctx context.Context, sessionHistory []models.Message, concept string) (string, error) {
	prompt := fmt.Sprintf(vg.prompts.QuizGeneratorPrompt, concept)

	response, err := vg.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating Quiz: %v", err)
	}

	response = cleanQuizCode(response)

	return response, nil
}

func cleanQuizCode(code string) string {
	code = removeMarkdownCodeBlocks(code)

	isFullHTML := strings.Contains(strings.ToLower(code), "<!doctype html>") ||
		strings.Contains(strings.ToLower(code), "<html")

	if !isFullHTML && !strings.Contains(code, "Quiz-wrapper") {
		code = fmt.Sprintf("<div class=\"Quiz-wrapper\">\n%s\n</div>", code)
	}

	if isFullHTML {
		if !strings.Contains(code, "<base") {
			headEnd := strings.Index(code, "</head>")
			if headEnd != -1 {
				code = code[:headEnd] + "<base target=\"_blank\">" + code[headEnd:]
			}
		}

		if !strings.Contains(code, "viewport") {
			headEnd := strings.Index(code, "</head>")
			if headEnd != -1 {
				code = code[:headEnd] + "<meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">" + code[headEnd:]
			}
		}
	}

	return code
}