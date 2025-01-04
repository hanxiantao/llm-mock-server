package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"llm-mock-server/provider/openai"
)

type ChatCompletionsHandler interface {
	ShouldHandleRequest(context *gin.Context) bool

	HandleChatCompletions(context *gin.Context)
}

var chatCompletionsHandlers = []ChatCompletionsHandler{
	&openai.Provider{},
}

func HandleChatCompletions(context *gin.Context) {
	for _, handler := range chatCompletionsHandlers {
		if handler.ShouldHandleRequest(context) {
			handler.HandleChatCompletions(context)
			return
		}
	}
	context.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
}
