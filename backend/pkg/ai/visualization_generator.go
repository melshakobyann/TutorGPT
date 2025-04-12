package ai

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"tutorgpt/pkg/models"
)

type VisualizationGenerator struct {
	openAIClient *OpenAIClient
	prompts      *Prompts
}

func NewVisualizationGenerator(openAIClient *OpenAIClient, prompts *Prompts) *VisualizationGenerator {
	return &VisualizationGenerator{
		openAIClient: openAIClient,
		prompts:      prompts,
	}
}

func (vg *VisualizationGenerator) GenerateVisualization(ctx context.Context, sessionHistory []models.Message, concept string) (string, error) {
	prompt := fmt.Sprintf(vg.prompts.VisualizationPrompt, concept)

	response, err := vg.openAIClient.GenerateResponse(ctx, sessionHistory, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating visualization: %v", err)
	}

	response = cleanVisualizationCode(response)

	return response, nil
}

func cleanVisualizationCode(code string) string {
	code = removeMarkdownCodeBlocks(code)

	isFullHTML := strings.Contains(strings.ToLower(code), "<!doctype html>") ||
		strings.Contains(strings.ToLower(code), "<html")

	if !isFullHTML && !strings.Contains(code, "visualization-wrapper") {
		code = fmt.Sprintf("<div class=\"visualization-wrapper\">\n%s\n</div>", code)
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

func removeMarkdownCodeBlocks(code string) string {
	markdownPattern := regexp.MustCompile("(?s)^```(?:html|javascript|js|css)?(.*?)```$")
	if markdownPattern.MatchString(code) {
		matches := markdownPattern.FindStringSubmatch(code)
		if len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	code = strings.TrimPrefix(code, "```")
	code = strings.TrimSuffix(code, "```")

	return strings.TrimSpace(code)
}
