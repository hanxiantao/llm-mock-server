package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"llm-mock-server/provider"
)

var port int

func init() {
	flag.IntVar(&port, "port", 3000, "Port to run the server on")
	flag.Parse()
}

func main() {
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}
	server := gin.Default()
	server.Use(CORS())

	handlers := provider.NewChatCompletionsHandlers()
	server.POST("/v1/chat/completions", func(context *gin.Context) {
		for _, handler := range handlers {
			if handler.ShouldHandleRequest(context) {
				handler.HandleChatCompletions(context)
				return // 提前返回，避免不必要的循环
			}
		}
		context.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
	})

	log.Printf("Starting server on port %d", port)
	if err := server.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
