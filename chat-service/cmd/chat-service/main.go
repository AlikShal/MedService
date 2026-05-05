package main

import (
	"chat-service/internal/app"
	"log"
)

func main() {
	application, port, err := app.SetupChatService()
	if err != nil {
		log.Fatalf("failed to setup chat service: %v", err)
	}

	if err := application.Run(":" + port); err != nil {
		log.Fatalf("failed to run chat service: %v", err)
	}
}
