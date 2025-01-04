package main

import (
	"flag"
	"log"
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
	server.POST("/v1/chat/completions", provider.HandleChatCompletions)

	log.Printf("Starting server on port %d", port)
	if err := server.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
