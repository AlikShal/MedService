package main

import (
	"log"
	"notification-service/internal/app"
)

func main() {
	r, port, err := app.SetupNotificationService()
	if err != nil {
		log.Fatalf("failed to setup notification service: %v", err)
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run notification service: %v", err)
	}
}
