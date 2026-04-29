package main

import (
	"appointment-service/internal/app"
	"log"
)

func main() {
	r, port, err := app.SetupAppointmentService()
	if err != nil {
		log.Fatalf("failed to setup appointment service: %v", err)
	}
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to run appointment service: %v", err)
	}
}
