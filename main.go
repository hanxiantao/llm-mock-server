package main

import (
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"llm-mock-server/provider/chat"
	"llm-mock-server/provider/embeddings"
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
	server.POST("/v1/chat/completions", chat.HandleChatCompletions)
	server.POST("/v1/embeddings", embeddings.HandleEmbeddings)

	log.Printf("Starting server on port %d", port)
	if err := server.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
