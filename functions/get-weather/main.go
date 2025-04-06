package functions

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	pkgFunction "github.com/enteresanlikk/go-gemini/pkg/function"
	"github.com/enteresanlikk/go-gemini/pkg/utils"
	"github.com/google/generative-ai-go/genai"
)

type WeatherResponse struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

func Execute(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall) {
	var weatherResponse WeatherResponse

	if err := utils.ArgsToStruct(call.Args, &weatherResponse); err != nil {
		log.Fatal(err)
		return
	}

	celcius := rand.Intn(30)
	prompt := "Simüle edilmiş hava durumu: " + weatherResponse.City + ", " + weatherResponse.Country + " için hava " + strconv.Itoa(celcius) + " derece."

	log.Println(prompt)

	funcResp := &genai.FunctionResponse{
		Name: "getWeather",
		Response: map[string]interface{}{
			"weather": prompt,
		},
	}

	session.SendMessage(ctx, funcResp)
}

func NewGetWeather() pkgFunction.Function {
	return pkgFunction.Function{
		Name:    "getWeather",
		Execute: Execute,
		FunctionDeclaration: genai.FunctionDeclaration{
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
	}
}

func init() {
	pkgFunction.RegisterFunction(NewGetWeather())
}
