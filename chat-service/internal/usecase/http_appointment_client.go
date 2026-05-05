package usecase

import (
	"chat-service/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPAppointmentClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPAppointmentClient(baseURL string) *HTTPAppointmentClient {
	return &HTTPAppointmentClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPAppointmentClient) GetAppointment(id string, authorizationHeader string) (*model.Appointment, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/appointments/%s", c.baseURL, id), nil)
	if err != nil {
		return nil, err
	}
	if authorizationHeader != "" {
		req.Header.Set("Authorization", authorizationHeader)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call appointment service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("appointment not found or service unavailable")
	}

	appointment := &model.Appointment{}
	if err := json.NewDecoder(resp.Body).Decode(appointment); err != nil {
		return nil, fmt.Errorf("failed to decode appointment response: %w", err)
	}

	return appointment, nil
}
