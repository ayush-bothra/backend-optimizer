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

const GroqEndPoint = "https://api.groq.com/openai/v1/chat/completions"


func NewRequest(messages []Message) *ReqBody {
	model := "llama3-8b-8192"
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