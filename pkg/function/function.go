package pkgFunction

import (
	"context"

	"github.com/google/generative-ai-go/genai"
)

type ExecuteFunction func(ctx context.Context, session *genai.ChatSession, call genai.FunctionCall)

type Function struct {
	Name                string
	Execute             ExecuteFunction
	FunctionDeclaration genai.FunctionDeclaration
}
