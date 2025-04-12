package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	SessionHistory []Message `json:"session_history"`
	MessageType    string    `json:"message_type"`
	Content        string    `json:"content"`
}

type ChatResponse struct {
	Response             string `json:"response"`
	VisualizationPayload string `json:"visualization_payload,omitempty"`
	Error                string `json:"error,omitempty"`
}
