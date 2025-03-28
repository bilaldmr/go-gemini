package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

const (
	model = "learnlm-1.5-pro-experimental"
)

var getWeatherTool = &genai.Tool{
	FunctionDeclarations: []*genai.FunctionDeclaration{
		&genai.FunctionDeclaration{
			Name:        "getWeather",
			Description: "gets the weather for a requested city",
			Parameters: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"city": &genai.Schema{
						Type: genai.TypeString,
					},
					"country": &genai.Schema{
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
}

type Response struct {
	GetWeather map[string]interface{} `json:"getWeather"`
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				str := fmt.Sprintf("%v", part)

				str = strings.Trim(str, "{}")
				if strings.HasPrefix(str, "getWeather") {
					mapStr := strings.TrimPrefix(str, "getWeather ")
					mapStr = strings.Trim(mapStr, "map[]")

					pairs := strings.Split(mapStr, " ")
					weatherResp := make(map[string]string)

					for _, pair := range pairs {
						if strings.Contains(pair, ":") {
							kv := strings.Split(pair, ":")
							if len(kv) == 2 {
								weatherResp[kv[0]] = kv[1]
							}
						}
					}

					// Convert to JSON
					prettyJSON, _ := json.MarshalIndent(weatherResp, "", "  ")
					fmt.Println(string(prettyJSON))
				}
			}
		}
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

	model := client.GenerativeModel(model)

	model.SetTemperature(0.8)
	model.SetTopK(64)
	model.SetTopP(0.95)
	model.SetMaxOutputTokens(8192)
	model.ResponseMIMEType = "text/plain"

	model.Tools = []*genai.Tool{
		getWeatherTool,
	}

	model.ToolConfig = &genai.ToolConfig{FunctionCallingConfig: &genai.FunctionCallingConfig{
		Mode: genai.FunctionCallingAny,
	}}

	session := model.StartChat()
	session.History = []*genai.Content{
		{
			Role: "user",
			Parts: []genai.Part{
				genai.Text("türkiye'de istanbul'da hava durumu nasıl?"),
			},
		},
	}

	resp, err := session.SendMessage(ctx, genai.Text("you are a rest api, you are given a request and you need to respond with the correct response"))
	if err != nil {
		log.Fatalf("Error sending message: %v", err)
	}

	printResponse(resp)
}
