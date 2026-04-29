package usecase

import (
	"appointment-service/internal/model"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HTTPDoctorClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPDoctorClient(baseURL string) *HTTPDoctorClient {
	return &HTTPDoctorClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPDoctorClient) CheckDoctorExists(doctorID string) (*model.Doctor, error) {
	url := fmt.Sprintf("%s/doctors/%s", c.baseURL, doctorID)
	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to call doctor service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("doctor not found or service unavailable")
	}

	var doctor model.Doctor
	if err := json.NewDecoder(resp.Body).Decode(&doctor); err != nil {
		return nil, fmt.Errorf("failed to decode doctor response: %w", err)
	}

	return &doctor, nil
}
