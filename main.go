package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/enteresanlikk/go-gemini/functions"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	model   = "learnlm-1.5-pro-experimental"
	process = make(map[string]functions.Process)
)

func loadProcesses() {
	weatherProcess := functions.NewWeatherProcess()
	process[weatherProcess.Name] = weatherProcess

	fmt.Println(len(process), "processes loaded successfully")
}

func processParts(ctx context.Context, session *genai.ChatSession, parts []genai.Part) {
	for _, part := range parts {
		if text, ok := part.(genai.Text); ok {
			fmt.Println(text)
		} else if call, ok := part.(genai.FunctionCall); ok {
			processFunctionCall(ctx, session, call)
		}
	}
}

func processFunctionCall(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall) {
	if processFunction, ok := process[call.Name]; ok {
		processFunction.Execute(ctx, session, call)
	} else {
		fmt.Println("Process not found:", call.Name)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}
	defer client.Close()

	loadProcesses()

	model := client.GenerativeModel(model)

	model.SetTemperature(0.8)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	for _, process := range process {
		model.Tools = append(model.Tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{&process.FunctionDeclaration},
		})
	}

	model.ToolConfig = &genai.ToolConfig{FunctionCallingConfig: &genai.FunctionCallingConfig{
		Mode: genai.FunctionCallingAny,
	}}

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text("türkiye'de istanbul'da hava durumu nasıl?"))
	if err != nil {
		panic(err)
	}

	processParts(ctx, session, resp.Candidates[0].Content.Parts)
}
