package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/google/generative-ai-go/genai"
)

var (
	GetWeatherName = "getWeather"
)

type WeatherRequest struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func GetWeather(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall) {
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

	celcius := rand.Intn(30)
	response := fmt.Sprintf("Simüle edilmiş hava durumu: %s, %s için hava %d derece.", request.City, request.Country, celcius)

	fmt.Println(response)

	funcResp := &genai.FunctionResponse{
		Name: GetWeatherName,
		Response: map[string]interface{}{
			"weather": response,
		},
	}

	session.SendMessage(ctx, funcResp)
}

func NewWeatherProcess() Process {
	getWeatherFunctionDeclaration := genai.FunctionDeclaration{
		Name:        GetWeatherName,
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
	}

	return Process{
		Name:                GetWeatherName,
		Execute:             GetWeather,
		FunctionDeclaration: getWeatherFunctionDeclaration,
	}
}
