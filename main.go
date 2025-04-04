package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	model = "learnlm-1.5-pro-experimental"
)

// WEATHER FUNCTION START

type WeatherRequest struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func doGetWeather(request WeatherRequest) string {
	return fmt.Sprintf("Simüle edilmiş hava durumu: %s, %s için hava güneşli ve 25 derece.", request.City, request.Country)
}

// WEATHER FUNCTION END

// PROCESS PARTS START

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
	if call.Name == "getWeather" {
		var request WeatherRequest

		argsJSON, err := json.Marshal(call.Args)
		if err != nil {
			fmt.Println("Error marshaling args:", err)
			return
		}

		if err := json.Unmarshal(argsJSON, &request); err != nil {
			fmt.Println("Error unmarshaling request:", err)
			return
		}

		response := doGetWeather(request)

		fmt.Println(response)

		funcResp := &genai.FunctionResponse{
			Name: "getWeather",
			Response: map[string]interface{}{
				"weather": response,
			},
		}

		session.SendMessage(ctx, funcResp)
	}
}

// PROCESS PARTS END

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

	model := client.GenerativeModel(model)

	model.SetTemperature(0.8)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	model.Tools = []*genai.Tool{
		{
			FunctionDeclarations: []*genai.FunctionDeclaration{
				{
					Name:        "getWeather",
					Description: "gets the weather for a requested city",
					Parameters: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"city": {
								Type: genai.TypeString,
							},
							"country": {
								Type: genai.TypeString,
							},
						},
						Required: []string{
							"city",
							"country",
						},
					},
				},
			},
		},
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
