package ai

import (
	"bytes"
	"net/http"
	"os"
)

type Message struct {
	Role string `json:"role"`
	Content string `json:"content"`
}

type ReqBody struct {
	Model string `json:"model"`
	Messages []Message `json:"messages"`
}

type Choice struct {
	Index int `json:"index"`
	Message Message `json:"message"`
	FinishReason string `json:"finish_reason"`
} 

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens int `json:"total_tokens"`
}

type RespBody struct {
	ID string `json:"id"`
	Model string `json:"model"`
	Choices []Choice `json:"choices"`
	Usage Usage `json:"usage"`
}

const GroqEndPoint = "https://api.groq.com/openai/v1/chat/completions"


func NewRequest(messages []Message) *ReqBody {
	model := "llama-3.1-8b-instant"
	req := ReqBody{Model: model, Messages: messages}
	return &req
}

func NewQuery(jsonBody []byte) (*http.Request, error) {
	postreq, err := http.NewRequest("POST", GroqEndPoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	postreq.Header.Set("Content-Type", "application/json")
	postreq.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))
	return postreq, nil;
}