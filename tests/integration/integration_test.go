package integration

import (
	"aviasales-shared-utils/config"
	"aviasales-shared-utils/http"
	"aviasales-shared-utils/providers"
	"aviasales-shared-utils/validation"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestFullWorkflow tests a complete workflow using all utilities
func TestFullWorkflow(t *testing.T) {
	// Initialize utilities
	timeProvider := providers.NewSystemTimeProvider()
	idGenerator := providers.NewPrefixedIDGenerator("req")
	validator := validation.NewFieldValidator()
	config := config.NewConfig()

	// Setup configuration
	config.Set("api_timeout", "30s")
	config.Set("max_retries", "3")
	config.Set("debug", "true")

	// Create request data
	type RequestData struct {
		ID        string `validate:"required"`
		Message   string `validate:"required,min=5,max=100"`
		Email     string `validate:"required,email"`
		Timestamp time.Time
	}

	request := RequestData{
		ID:        idGenerator.Generate(),
		Message:   "Hello from integration test",
		Email:     "test@example.com",
		Timestamp: timeProvider.Now(),
	}

	// Validate request
	if err := validator.Validate(request); err != nil {
		t.Fatalf("Request validation failed: %v", err)
	}

	// Get configuration
	timeout := config.GetDurationWithDefault("api_timeout", 10*time.Second)
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout 30s, got %v", timeout)
	}

	maxRetries := config.GetIntWithDefault("max_retries", 1)
	if maxRetries != 3 {
		t.Errorf("Expected max_retries 3, got %v", maxRetries)
	}

	debug := config.GetBoolWithDefault("debug", false)
	if !debug {
		t.Errorf("Expected debug true, got %v", debug)
	}

	// Create HTTP client and make request
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(200)
		http.WriteJSONResponse(w, 200, map[string]interface{}{
			"status":  "success",
			"id":      request.ID,
			"message": "Request processed",
		})
	}))
	defer server.Close()

	client := http.NewClientWithTimeout(server.URL, timeout)
	resp, err := client.Post("/process", request, nil)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %v", resp.StatusCode)
	}

	// Parse response
	var response map[string]interface{}
	if err := http.ParseJSONResponse(resp, &response); err != nil {
		t.Fatalf("Response parsing failed: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status success, got %v", response["status"])
	}

	if response["id"] != request.ID {
		t.Errorf("Expected ID %v, got %v", request.ID, response["id"])
	}
}

// TestTimeProviderIntegration tests time provider with other components
func TestTimeProviderIntegration(t *testing.T) {
	fixedTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	timeProvider := providers.NewFixedTimeProvider(fixedTime)

	// Test with configuration
	config := config.NewConfig()
	config.Set("created_at", timeProvider.Now().Format(time.RFC3339))

	retrievedTime, err := time.Parse(time.RFC3339, config.Get("created_at"))
	if err != nil {
		t.Fatalf("Time parsing failed: %v", err)
	}

	if !retrievedTime.Equal(fixedTime) {
		t.Errorf("Expected time %v, got %v", fixedTime, retrievedTime)
	}
}

// TestIDGeneratorIntegration tests ID generator with validation
func TestIDGeneratorIntegration(t *testing.T) {
	generator := providers.NewUUIDGenerator()
	validator := validation.NewFieldValidator()

	type IDRequest struct {
		ID string `validate:"required,min=10"`
	}

	// Generate and validate multiple IDs
	for i := 0; i < 10; i++ {
		request := IDRequest{
			ID: generator.Generate(),
		}

		if err := validator.Validate(request); err != nil {
			t.Errorf("ID validation failed for iteration %d: %v", i, err)
		}

		if len(request.ID) < 10 {
			t.Errorf("ID too short: %v", request.ID)
		}
	}
}

// TestHTTPClientWithConfig tests HTTP client using configuration
func TestHTTPClientWithConfig(t *testing.T) {
	config := config.NewConfig()
	config.Set("api_base_url", "https://api.example.com")
	config.Set("timeout", "15s")
	config.Set("retry_count", "2")

	baseURL := config.Get("api_base_url")
	timeout := config.GetDurationWithDefault("timeout", 30*time.Second)
	retryCount := config.GetIntWithDefault("retry_count", 1)

	client := http.NewClientWithTimeout(baseURL, timeout)

	if client == nil {
		t.Fatal("HTTP client creation failed")
	}

	// Verify configuration was applied
	if timeout != 15*time.Second {
		t.Errorf("Expected timeout 15s, got %v", timeout)
	}

	if retryCount != 2 {
		t.Errorf("Expected retry count 2, got %v", retryCount)
	}
}

// TestValidationWithHTTPResponse tests validation with HTTP responses
func TestValidationWithHTTPResponse(t *testing.T) {
	validator := validation.NewFieldValidator()

	type APIResponse struct {
		Status  string `validate:"required"`
		Message string `validate:"required,min=5"`
		Code    int    `validate:"min=200,max=599"`
	}

	// Create test server
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		response := APIResponse{
			Status:  "success",
			Message: "Operation completed successfully",
			Code:    200,
		}

		// Validate response before sending
		if err := validator.Validate(response); err != nil {
			http.WriteErrorResponse(w, 500, "Invalid response format")
			return
		}

		http.WriteJSONResponse(w, 200, response)
	}))
	defer server.Close()

	client := http.NewClient(server.URL)
	resp, err := client.Get("/test", nil)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	var apiResponse APIResponse
	if err := http.ParseJSONResponse(resp, &apiResponse); err != nil {
		t.Fatalf("Response parsing failed: %v", err)
	}

	// Validate received response
	if err := validator.Validate(apiResponse); err != nil {
		t.Errorf("Response validation failed: %v", err)
	}
}

// TestErrorHandlingWorkflow tests error handling across all utilities
func TestErrorHandlingWorkflow(t *testing.T) {
	config := config.NewConfig()
	validator := validation.NewFieldValidator()

	// Test configuration error handling
	_, err := config.GetInt("nonexistent_key")
	if err == nil {
		t.Error("Expected error for nonexistent config key")
	}

	// Test validation error handling
	type InvalidData struct {
		Email string `validate:"email"`
	}

	invalidData := InvalidData{Email: "invalid-email"}
	if err := validator.Validate(invalidData); err == nil {
		t.Error("Expected validation error for invalid email")
	}

	// Test HTTP error handling
	server := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		http.WriteErrorResponse(w, 400, "Bad request")
	}))
	defer server.Close()

	client := http.NewClient(server.URL)
	resp, err := client.Get("/error", nil)
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 400 {
		t.Errorf("Expected status 400, got %v", resp.StatusCode)
	}
}
