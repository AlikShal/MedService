package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"
)

type endpoint struct {
	Name         string
	Method       string
	Path         string
	RequiresAuth bool
}

type doctor struct {
	ID string `json:"id"`
}

type authPayload struct {
	Token string `json:"token"`
}

type loadContext struct {
	Token    string
	DoctorID string
}

type requestResult struct {
	Endpoint     string `json:"endpoint"`
	Method       string `json:"method"`
	Path         string `json:"path"`
	StatusCode   int    `json:"statusCode"`
	ResponseMs   int64  `json:"responseMs"`
	Error        string `json:"error,omitempty"`
	Round        int    `json:"round"`
	Worker       int    `json:"worker"`
	TimestampUTC string `json:"timestampUtc"`
}

type endpointSummary struct {
	Endpoint      string  `json:"endpoint"`
	TotalRequests int     `json:"totalRequests"`
	SuccessCount  int     `json:"successCount"`
	FailCount     int     `json:"failCount"`
	AverageMs     float64 `json:"averageMs"`
	P95Ms         int64   `json:"p95Ms"`
}

type outputFile struct {
	TimestampUTC string            `json:"timestampUtc"`
	BaseURL      string            `json:"baseUrl"`
	Rounds       int               `json:"rounds"`
	Workers      int               `json:"workers"`
	Summaries    []endpointSummary `json:"summaries"`
	Results      []requestResult   `json:"results"`
}

var endpoints = []endpoint{
	{
		Name:   "GET /api/appointments",
		Method: http.MethodGet,
		Path:   "/api/appointments",
		RequiresAuth: true,
	},
	{
		Name:   "POST /api/appointments",
		Method: http.MethodPost,
		Path:   "/api/appointments",
		RequiresAuth: true,
	},
	{
		Name:   "GET /api/doctors",
		Method: http.MethodGet,
		Path:   "/api/doctors",
	},
	{
		Name:   "GET /api/patients/me",
		Method: http.MethodGet,
		Path:   "/api/patients/me",
		RequiresAuth: true,
	},
	{
		Name:   "GET /api/health/auth",
		Method: http.MethodGet,
		Path:   "/api/health/auth",
	},
}

func main() {
	baseURL := flag.String("url", "http://localhost:8080", "base URL for the API gateway or service")
	rounds := flag.Int("rounds", 3, "number of load test rounds")
	workers := flag.Int("workers", 50, "concurrent requests per endpoint per round")
	flag.Parse()

	if *rounds < 1 {
		fmt.Fprintln(os.Stderr, "--rounds must be at least 1")
		os.Exit(1)
	}
	if *workers < 1 {
		fmt.Fprintln(os.Stderr, "--workers must be at least 1")
		os.Exit(1)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	base := normalizeBaseURL(*baseURL)
	ctx, err := bootstrapLoadContext(client, base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to bootstrap load context: %v\n", err)
		os.Exit(1)
	}

	results := runLoadTest(client, base, ctx, *rounds, *workers)
	summaries := summarize(results)
	printReport(summaries, results)

	if err := saveResults("load_test_results.json", base, *rounds, *workers, summaries, results); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save load_test_results.json: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("Saved full results to load_test_results.json")
}

func runLoadTest(client *http.Client, baseURL string, ctx loadContext, rounds int, workers int) []requestResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]requestResult, 0, rounds*workers*len(endpoints))

	for round := 1; round <= rounds; round++ {
		for worker := 1; worker <= workers; worker++ {
			for _, target := range endpoints {
				wg.Add(1)
				go func(round int, worker int, target endpoint) {
					defer wg.Done()
					result := sendRequest(client, baseURL, ctx, round, worker, target)

					mu.Lock()
					results = append(results, result)
					mu.Unlock()
				}(round, worker, target)
			}
		}
		wg.Wait()
	}

	return results
}

func sendRequest(client *http.Client, baseURL string, ctx loadContext, round int, worker int, target endpoint) requestResult {
	result := requestResult{
		Endpoint:     target.Name,
		Method:       target.Method,
		Path:         target.Path,
		Round:        round,
		Worker:       worker,
		TimestampUTC: time.Now().UTC().Format(time.RFC3339),
	}

	var body *bytes.Reader
	requestBody := buildRequestBody(target, ctx, round, worker)
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			result.Error = err.Error()
			return result
		}
		body = bytes.NewReader(payload)
	} else {
		body = bytes.NewReader(nil)
	}

	req, err := http.NewRequest(target.Method, baseURL+target.Path, body)
	if err != nil {
		result.Error = err.Error()
		return result
	}
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if target.RequiresAuth {
		req.Header.Set("Authorization", "Bearer "+ctx.Token)
	}

	start := time.Now()
	resp, err := client.Do(req)
	result.ResponseMs = time.Since(start).Milliseconds()
	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	result.StatusCode = resp.StatusCode
	return result
}

func buildRequestBody(target endpoint, ctx loadContext, round int, worker int) map[string]string {
	if target.Method != http.MethodPost || target.Path != "/api/appointments" {
		return nil
	}

	scheduledAt := time.Now().UTC().
		Add(time.Duration(round*worker) * time.Minute).
		Format(time.RFC3339)

	return map[string]string{
		"title":        fmt.Sprintf("Load test appointment r%d-w%d", round, worker),
		"description":  "Automated load test request",
		"doctor_id":    ctx.DoctorID,
		"scheduled_at": scheduledAt,
	}
}

func bootstrapLoadContext(client *http.Client, baseURL string) (loadContext, error) {
	email := fmt.Sprintf("loadtest-%d@medsync.local", time.Now().UTC().UnixNano())
	password := "loadtest123"

	registerBody := map[string]string{
		"full_name": "Load Test Patient",
		"email":     email,
		"password":  password,
	}

	payload, err := doJSONRequest[authPayload](client, http.MethodPost, baseURL+"/api/auth/register", registerBody, "")
	if err != nil {
		return loadContext{}, fmt.Errorf("register user: %w", err)
	}

	profileBody := map[string]string{
		"full_name":     "Load Test Patient",
		"phone":         "+77000000000",
		"date_of_birth": "1999-01-01",
		"notes":         "generated for load testing",
	}

	if _, err := doJSONRequest[map[string]any](client, http.MethodPut, baseURL+"/api/patients/me", profileBody, payload.Token); err != nil {
		return loadContext{}, fmt.Errorf("create patient profile: %w", err)
	}

	doctors, err := doJSONRequest[[]doctor](client, http.MethodGet, baseURL+"/api/doctors", nil, "")
	if err != nil {
		return loadContext{}, fmt.Errorf("fetch doctors: %w", err)
	}
	if len(doctors) == 0 {
		return loadContext{}, fmt.Errorf("no doctors available for appointment creation")
	}

	return loadContext{
		Token:    payload.Token,
		DoctorID: doctors[0].ID,
	}, nil
}

func doJSONRequest[T any](client *http.Client, method string, url string, body any, token string) (T, error) {
	var zero T

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			return zero, err
		}
		reader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return zero, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		return zero, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, fmt.Errorf("status %d: %s", resp.StatusCode, string(respBody))
	}

	if len(respBody) == 0 {
		return zero, nil
	}

	var decoded T
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return zero, err
	}

	return decoded, nil
}

func summarize(results []requestResult) []endpointSummary {
	summaries := make([]endpointSummary, 0, len(endpoints))

	for _, target := range endpoints {
		times := make([]int64, 0)
		summary := endpointSummary{Endpoint: target.Name}

		for _, result := range results {
			if result.Endpoint != target.Name {
				continue
			}

			summary.TotalRequests++
			times = append(times, result.ResponseMs)

			if result.Error == "" && result.StatusCode >= 200 && result.StatusCode < 300 {
				summary.SuccessCount++
			} else {
				summary.FailCount++
			}
		}

		summary.AverageMs = average(times)
		summary.P95Ms = p95(times)
		summaries = append(summaries, summary)
	}

	return summaries
}

func average(values []int64) float64 {
	if len(values) == 0 {
		return 0
	}

	var total int64
	for _, value := range values {
		total += value
	}

	return float64(total) / float64(len(values))
}

func p95(values []int64) int64 {
	if len(values) == 0 {
		return 0
	}

	sorted := append([]int64(nil), values...)
	sort.Slice(sorted, func(i int, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(math.Ceil(float64(len(sorted))*0.95)) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

func printReport(summaries []endpointSummary, results []requestResult) {
	fmt.Println("Healthcare Microservices Load Test Results")
	fmt.Println()
	fmt.Printf("%-28s %10s %10s %10s %12s %10s\n", "Endpoint", "Total", "Success", "Fail", "Avg ms", "P95 ms")
	fmt.Println("------------------------------------------------------------------------------------")

	totalRequests := 0
	totalSuccess := 0
	slowest := endpointSummary{}

	for i, summary := range summaries {
		totalRequests += summary.TotalRequests
		totalSuccess += summary.SuccessCount
		if i == 0 || summary.AverageMs > slowest.AverageMs {
			slowest = summary
		}

		fmt.Printf(
			"%-28s %10d %10d %10d %12.2f %10d\n",
			summary.Endpoint,
			summary.TotalRequests,
			summary.SuccessCount,
			summary.FailCount,
			summary.AverageMs,
			summary.P95Ms,
		)
	}

	successRate := 0.0
	if totalRequests > 0 {
		successRate = (float64(totalSuccess) / float64(totalRequests)) * 100
	}

	fmt.Println()
	fmt.Println("Final Summary")
	fmt.Printf("Total requests: %d\n", totalRequests)
	fmt.Printf("Overall success rate: %.2f%%\n", successRate)
	fmt.Printf("Slowest endpoint: %s (avg %.2f ms)\n", slowest.Endpoint, slowest.AverageMs)

	failures := 0
	for _, result := range results {
		if result.Error != "" || result.StatusCode < 200 || result.StatusCode >= 300 {
			failures++
		}
	}
	fmt.Printf("Total failures: %d\n", failures)
}

func saveResults(path string, baseURL string, rounds int, workers int, summaries []endpointSummary, results []requestResult) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(outputFile{
		TimestampUTC: time.Now().UTC().Format(time.RFC3339),
		BaseURL:      baseURL,
		Rounds:       rounds,
		Workers:      workers,
		Summaries:    summaries,
		Results:      results,
	})
}

func normalizeBaseURL(baseURL string) string {
	for len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	return baseURL
}
