package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"tutorgpt/pkg/ai"
	"tutorgpt/pkg/models"
	"tutorgpt/pkg/utils"
)


type ChatHandler struct {
	config                 *utils.Config
	openAIClient           *ai.OpenAIClient
	prompts                *ai.Prompts
	orchestrator           *ai.Orchestrator
	lectureBuilder         *ai.LectureBuilder
	vagueExplainer         *ai.VagueExplainer
	taskManager            *ai.TaskManager
	qaAgent                *ai.QAAgent
	visualizationGenerator *ai.VisualizationGenerator
	chitchat 			   *ai.ChitchatAgent
	quizGenerator          *ai.QuizGenerator

	firstTimeUsers map[string]bool
}

func NewChatHandler(config *utils.Config) *ChatHandler {
	handler := &ChatHandler{
		config:         config,
		firstTimeUsers: make(map[string]bool),
	}
	handler.initAIModules()
	return handler
}

func (h *ChatHandler) initAIModules() {
	h.openAIClient = ai.NewOpenAIClient(h.config.OpenAIAPIKey)

	h.prompts = ai.NewPrompts()

	h.orchestrator = ai.NewOrchestrator(h.openAIClient, h.prompts)

	h.lectureBuilder = ai.NewLectureBuilder(h.openAIClient, h.prompts)
	h.vagueExplainer = ai.NewVagueExplainer(h.openAIClient, h.prompts)
	h.taskManager = ai.NewTaskManager(h.openAIClient, h.prompts)
	h.qaAgent = ai.NewQAAgent(h.openAIClient, h.prompts)
	h.chitchat = ai.NewChitchatAgent(h.openAIClient, h.prompts)
	h.visualizationGenerator = ai.NewVisualizationGenerator(h.openAIClient, h.prompts)
	h.quizGenerator = ai.NewQuizGenerator(h.openAIClient, h.prompts)
}

func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.ChatRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&req); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Received %s request with content: %s", req.MessageType, req.Content)
	log.Printf("Session history length: %d", len(req.SessionHistory))

	if len(req.SessionHistory) == 0 || req.MessageType == "welcome" {
		log.Printf("First-time user or welcome request detected, sending welcome message")
		welcomeMsg := h.orchestrator.GetWelcomeMessage()
		response := &models.ChatResponse{
			Response: welcomeMsg,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error encoding welcome message: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		} else {
			log.Printf("Welcome message sent successfully")
		}
		return
	}

	log.Printf("Forwarding to decision engine")
	response, err := h.processRequest(req)
	if err != nil {
		log.Printf("Error processing request: %v", err)

		errorResponse := &models.ChatResponse{
			Error: err.Error(),
		}

		if response != nil && response.Response != "" {
			errorResponse.Response = response.Response
		} else {
			errorResponse.Response = "I apologize, but I encountered an issue processing your request. Please try again later."
		}

		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		if encodeErr := encoder.Encode(errorResponse); encodeErr != nil {
			log.Printf("Error encoding error response: %v", encodeErr)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		} else {
			log.Printf("Error response sent successfully")
		}
		return
	}

	log.Printf("Sending response back to client: %s", response.Response[:min(10, len(response.Response))]+"...")
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	} else {
		log.Printf("Response sent successfully")
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *ChatHandler) processRequest(req models.ChatRequest) (*models.ChatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	log.Printf("Calling orchestrator to determine intent")
	orchestratorResult, err := h.orchestrator.DetermineIntent(ctx, req.SessionHistory, req.Content)
	if err != nil {
		log.Printf("Error determining intent: %v", err)
		log.Printf("Falling back to direct chat handler")
		return h.handleChatRequest(ctx, req)
	}

	log.Printf("Orchestrator selected tool: %s with instructions: %s", orchestratorResult.Tool, orchestratorResult.Instructions)

	if orchestratorResult.Instructions == "" {
		orchestratorResult.Instructions = req.Content
		log.Printf("Empty instructions detected, using original content: %s", orchestratorResult.Instructions)
	}

	var result *models.ChatResponse
	var resultErr error

	switch orchestratorResult.Tool {
	case "lecture_builder":
		log.Printf("Routing to lecture builder")
		result, resultErr = h.handleLectureRequest(ctx, req, orchestratorResult.Instructions)
	case "vague_explainer":
		log.Printf("Routing to vague explainer")
		result, resultErr = h.handleVagueRequest(ctx, req, orchestratorResult.Instructions)
	case "qa_agent":
		log.Printf("Routing to QA agent")
		result, resultErr = h.handleQARequest(ctx, req, orchestratorResult.Instructions)
	case "task_manager":
		log.Printf("Routing to task manager")
		result, resultErr = h.handleTaskRequest(ctx, req, orchestratorResult.Instructions)
	case "task_checker":
		log.Printf("Routing to task checker")
		result, resultErr = h.handleTaskCheckRequest(ctx, req, orchestratorResult.Instructions)
	case "visualization_generator":
		log.Printf("Routing to visualization generator")
		result, resultErr = h.handleVisualizationRequest(ctx, req, orchestratorResult.Instructions)
	case "chitchat":
		log.Printf("Routing to chitchat")
		result, resultErr = h.handleChitchatRequest(ctx, req, orchestratorResult.Instructions)
	case "quiz":
		log.Printf("Routing to quiz")
		result, resultErr = h.handleQuizRequest(ctx, req, orchestratorResult.Instructions)
	default:
		log.Printf("Unknown tool '%s' selected by orchestrator, defaulting to lecture", orchestratorResult.Tool)
		result, resultErr = h.handleChatRequest(ctx, req)
	}

	if result == nil {
		log.Printf("No result returned from tool %s", orchestratorResult.Tool)
		result = &models.ChatResponse{
			Response: fmt.Sprintf("I attempted to process your request about '%s' but encountered a technical issue. Please try rephrasing your question.", orchestratorResult.Instructions),
		}

		return result, fmt.Errorf("error from %s : %v", orchestratorResult.Tool, resultErr)
	}

	if result.Response == "" {
		log.Printf("Empty response from tool %s", orchestratorResult.Tool)
		result.Response = fmt.Sprintf("I processed your request about '%s' but couldn't generate a proper response.", orchestratorResult.Instructions)
	}

	if resultErr != nil {
		log.Printf("Error from selected tool %s: %v", orchestratorResult.Tool, resultErr)
		return result, fmt.Errorf("error from %s: %v", orchestratorResult.Tool, resultErr)
	}

	log.Printf("Successfully processed request with %s", orchestratorResult.Tool)
	log.Printf("Response from %s: %s", orchestratorResult.Tool, result.Response[:min(100, len(result.Response))]+"...")
	return result, nil
}

func (h *ChatHandler) handleLectureRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	response, err := h.lectureBuilder.GenerateLecture(ctx, req.SessionHistory, instructions)
	if err != nil {
		if response != "" {
			return &models.ChatResponse{
				Response: response,
				Error:    err.Error(),
			}, err
		}
		return nil, fmt.Errorf("error generating lecture: %v", err)
	}

	return &models.ChatResponse{
		Response: response,
	}, nil
}

func (h *ChatHandler) handleVagueRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	response, err := h.vagueExplainer.GenerateVagueExplanation(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error generating vague explanation: %v", err)
	}

	return &models.ChatResponse{
		Response: response,
	}, nil
}

func (h *ChatHandler) handleQARequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	response, err := h.qaAgent.AnswerQuestion(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error answering question: %v", err)
	}

	return &models.ChatResponse{
		Response: response,
	}, nil
}

func (h *ChatHandler) handleChitchatRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	response, err := h.chitchat.Chitchat(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error in response to: %v", err)
	}

	return &models.ChatResponse{
		Response: response,
	}, nil
}

func (h *ChatHandler) handleChatRequest(ctx context.Context, req models.ChatRequest) (*models.ChatResponse, error) {
	var response string
	var err error

	response, err = h.chitchat.Chitchat(ctx, req.SessionHistory, req.Content)
	if err != nil {
		if response != "" {
			return &models.ChatResponse{
				Response: response,
				Error:    err.Error(),
			}, err
		}

		return nil, fmt.Errorf("error generating chat response: %v", err)
	}

	return &models.ChatResponse{
		Response: response,
	}, nil
}

func (h *ChatHandler) handleVisualizationRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	visualizationHTML, err := h.visualizationGenerator.GenerateVisualization(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error generating visualization: %v", err)
	}

	return &models.ChatResponse{
		Response:             "Here's the visualization:",
		VisualizationPayload: visualizationHTML,
	}, nil
}

func (h *ChatHandler) handleQuizRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	visualizationHTML, err := h.quizGenerator.GenerateQuiz(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error generating visualization: %v", err)
	}

	return &models.ChatResponse{
		Response:             "Here's your quiz:",
		VisualizationPayload: visualizationHTML,
	}, nil
}

func (h *ChatHandler) handleTaskRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	task, err := h.taskManager.GenerateTask(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error generating task: %v", err)
	}

	return &models.ChatResponse{
		Response: task,
	}, nil
}

func (h *ChatHandler) handleTaskCheckRequest(ctx context.Context, req models.ChatRequest, instructions string) (*models.ChatResponse, error) {
	task, err := h.taskManager.CheckSubmission(ctx, req.SessionHistory, instructions)
	if err != nil {
		return nil, fmt.Errorf("error generating task: %v", err)
	}

	return &models.ChatResponse{
		Response: task,
	}, nil
}