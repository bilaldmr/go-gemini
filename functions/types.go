package functions

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

type ProcessFunction func(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall)

type Process struct {
	Name                string
	Execute             ProcessFunction
	FunctionDeclaration genai.FunctionDeclaration
}
