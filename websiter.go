package main

import (
	"context"
	_ "embed"
	"log"

	"google.golang.org/genai"
)

//go:embed sysprompt.txt
var systemPrompt string

//go:embed gemini.txt
var apiGemini string

//go:embed data.csv
var data string

var variant string = "gemini-2.0-flash-lite"

func Makesite(ctx context.Context, in string) (string, error) {
	// ctx??
	agent := NewAgent(ctx)
	text, err := agent.Run(ctx, in)
	if err != nil {
		return "", err
	}
	return text, nil
}

func (a *Agent) Run(ctx context.Context, url string) (string, error) {
	conversation := []*genai.Content{}
	localFile := genai.NewContentFromText(data, "user")
	conversation = append(conversation, localFile)

	userMessage := genai.NewContentFromText(url, "user")
	conversation = append(conversation, userMessage)

	message, err := a.runInference(ctx, conversation)
	if err != nil {
		return "", err
	}
	conversation = append(conversation, message)
	return message.Parts[0].Text, nil
}

// Call API
func (a *Agent) runInference(ctx context.Context, conversation []*genai.Content) (*genai.Content, error) {
	result, err := a.client.Models.GenerateContent(ctx, variant, conversation, &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{Parts: []*genai.Part{{Text: systemPrompt}}},
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if result == nil {
		log.Println("result is nil")
	}
	return result.Candidates[0].Content, err
}

func NewAgent(ctx context.Context) *Agent {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiGemini,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatal(err)
	}
	return &Agent{
		client: client,
	}
}

type Agent struct {
	client *genai.Client
}
