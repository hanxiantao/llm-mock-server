package chat

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"llm-mock-server/provider"
)

type requestHandler interface {
	provider.CommonRequestHandler

	HandleChatCompletions(context *gin.Context)
}

var chatCompletionsHandlers = []requestHandler{
	&openAiProvider{},
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
