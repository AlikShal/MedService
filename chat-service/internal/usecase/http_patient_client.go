package usecase

import (
	"chat-service/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type HTTPPatientClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPPatientClient(baseURL string) *HTTPPatientClient {
	return &HTTPPatientClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *HTTPPatientClient) GetAuthorizedPatient(authorizationHeader string) (*model.Patient, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/patients/me", c.baseURL), nil)
	if err != nil {
		return nil, err
	}
	if authorizationHeader != "" {
		req.Header.Set("Authorization", authorizationHeader)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call patient service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("patient profile not found or service unavailable")
	}

	patient := &model.Patient{}
	if err := json.NewDecoder(resp.Body).Decode(patient); err != nil {
		return nil, fmt.Errorf("failed to decode patient response: %w", err)
	}

	return patient, nil
}
