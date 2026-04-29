package main

import (
	"doctor-service/internal/app"
	"log"
)

func main() {
	r, port, err := app.SetupDoctorService()
	if err != nil {
		log.Fatalf("failed to setup doctor service: %v", err)
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run doctor service: %v", err)
	}
}
