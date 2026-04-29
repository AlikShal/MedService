package main

import (
	"log"
	"patient-service/internal/app"
)

func main() {
	application, port, err := app.SetupPatientService()
	if err != nil {
		log.Fatalf("failed to setup patient service: %v", err)
	}

	if err := application.Run(":" + port); err != nil {
		log.Fatalf("failed to run patient service: %v", err)
	}
}
