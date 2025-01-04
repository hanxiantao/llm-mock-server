package provider

import (
	"github.com/gin-gonic/gin"
	"llm-mock-server/provider/openai"
)

type ChatCompletionsHandler interface {
	ShouldHandleRequest(context *gin.Context) bool

	HandleChatCompletions(context *gin.Context)
}

func NewChatCompletionsHandlers() []ChatCompletionsHandler {
	return []ChatCompletionsHandler{
		&openai.Provider{},
	}
}
