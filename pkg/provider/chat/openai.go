package chat

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"llm-mock-server/pkg/utils"
)

type openAiProvider struct {
}

func (p *openAiProvider) ShouldHandleRequest(context *gin.Context) bool {
	return true
}

func (p *openAiProvider) HandleChatCompletions(context *gin.Context) {
	var chatRequest chatCompletionRequest
	if err := context.ShouldBindJSON(&chatRequest); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := utils.Validate.Struct(chatRequest); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			context.JSON(http.StatusBadRequest, gin.H{"error": fieldError.Error()})
			return
		}
	}
	prompt := "This is a mock server."
	if chatRequest.Messages[len(chatRequest.Messages)-1].IsStringContent() {
		prompt = chatRequest.Messages[len(chatRequest.Messages)-1].StringContent()
	}
	response := utils.Prompt2Response(prompt)

	if chatRequest.Stream {
		utils.SetEventStreamHeaders(context)
		dataChan := make(chan string)
		stopChan := make(chan bool)
		streamResponse := chatCompletionResponse{
			Id:      completionIdPrefix + uuid.New().String(),
			Object:  objectChatCompletionChunk,
			Created: time.Now().Unix(),
			Model:   chatRequest.Model,
		}
		streamResponseChoice := chatCompletionChoice{Delta: &chatMessage{}}
		go func() {
			for i, s := range response {
				streamResponseChoice.Delta.Content = string(s)
				if i == len(response)-1 {
					streamResponseChoice.FinishReason = stopReason
				}
				streamResponse.Choices = []chatCompletionChoice{streamResponseChoice}
				jsonStr, _ := json.Marshal(streamResponse)
				dataChan <- string(jsonStr)
			}
			stopChan <- true
		}()

		context.Stream(func(w io.Writer) bool {
			select {
			case data := <-dataChan:
				context.Render(-1, streamEvent{Data: "data: " + data})
				return true
			case <-stopChan:
				context.Render(-1, streamEvent{Data: "data: [DONE]"})
				return false
			}
		})
	} else {
		completion := createCompletion(chatRequest.Model, response)
		context.JSON(http.StatusOK, completion)
	}
}

func createCompletion(model, response string) chatCompletionResponse {
	return chatCompletionResponse{
		Id:      completionIdPrefix + uuid.New().String(),
		Object:  objectChatCompletion,
		Created: time.Now().Unix(),
		Model:   model,
		Choices: []chatCompletionChoice{
			{
				Index: 0,
				Message: &chatMessage{
					Role:    roleAssistant,
					Content: response,
				},
				FinishReason: stopReason,
			},
		},
		Usage: usage{
			PromptTokens:     9,
			CompletionTokens: 1,
			TotalTokens:      10,
		},
	}
}
