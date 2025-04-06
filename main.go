package main

import (
	"context"
	"errors"
	"log"
	"os"

	_ "github.com/enteresanlikk/go-gemini/functions"
	pkgFunction "github.com/enteresanlikk/go-gemini/pkg/function"
	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

var (
	ErrorLoadingEnv     = errors.New("error loading .env file")
	ErrorCreatingClient = errors.New("error creating client")
	ErrorGenerating     = errors.New("error generating")
	ErrProcessNotFound  = errors.New("process not found")
)

var (
	model = "learnlm-1.5-pro-experimental"
)

func processParts(ctx context.Context, session *genai.ChatSession, parts []genai.Part) {
	for _, part := range parts {
		if call, ok := part.(genai.FunctionCall); ok {
			processFunctionCall(ctx, session, call)
		}
	}
}

func processFunctionCall(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall) {
	if processFunction, ok := pkgFunction.GetFunction(call.Name); ok {
		processFunction.Execute(ctx, session, call)
	} else {
		log.Fatal(ErrProcessNotFound)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(ErrorLoadingEnv)
	}

	ctx := context.Background()

	apiKey := os.Getenv("GEMINI_API_KEY")

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(ErrorCreatingClient)
	}
	defer client.Close()

	model := client.GenerativeModel(model)

	model.SetTemperature(0.8)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	functions := pkgFunction.GetAllFunctions()

	log.Println(len(functions), "functions loaded successfully")

	for _, function := range functions {
		model.Tools = append(model.Tools, &genai.Tool{
			FunctionDeclarations: []*genai.FunctionDeclaration{&function.FunctionDeclaration},
		})
	}

	model.ToolConfig = &genai.ToolConfig{FunctionCallingConfig: &genai.FunctionCallingConfig{
		Mode: genai.FunctionCallingAny,
	}}

	session := model.StartChat()

	resp, err := session.SendMessage(ctx, genai.Text("türkiye'de istanbul'da hava durumu nasıl?"))
	if err != nil {
		log.Fatal(ErrorGenerating)
	}

	processParts(ctx, session, resp.Candidates[0].Content.Parts)
}
