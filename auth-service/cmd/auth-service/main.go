package main

import (
	"auth-service/internal/app"
	"log"
)

func main() {
	application, port, err := app.SetupAuthService()
	if err != nil {
		log.Fatalf("failed to setup auth service: %v", err)
	}

	if err := application.Run(":" + port); err != nil {
		log.Fatalf("failed to run auth service: %v", err)
	}
}
