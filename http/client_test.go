package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	baseURL := "https://api.example.com"
	client := NewClient(baseURL)

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.baseURL != baseURL {
		t.Errorf("NewClient() baseURL = %v, want %v", client.baseURL, baseURL)
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("NewClient() timeout = %v, want %v", client.httpClient.Timeout, 30*time.Second)
	}
}

func TestNewClientWithTimeout(t *testing.T) {
	baseURL := "https://api.example.com"
	timeout := 10 * time.Second
	client := NewClientWithTimeout(baseURL, timeout)

	if client == nil {
		t.Fatal("NewClientWithTimeout() returned nil")
	}

	if client.httpClient.Timeout != timeout {
		t.Errorf("NewClientWithTimeout() timeout = %v, want %v", client.httpClient.Timeout, timeout)
	}
}

func TestClient_Get(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %v", r.Method)
		}
		if r.URL.Path != "/test" {
			t.Errorf("Expected path /test, got %v", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer token" {
			t.Errorf("Expected Authorization header, got %v", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	headers := map[string]string{
		"Authorization": "Bearer token",
	}

	resp, err := client.Get("/test", headers)
	if err != nil {
		t.Fatalf("Client.Get() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Client.Get() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Post(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %v", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %v", r.Header.Get("Content-Type"))
		}

		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if body["key"] != "value" {
			t.Errorf("Expected body key=value, got %v", body)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	requestBody := map[string]interface{}{
		"key": "value",
	}

	resp, err := client.Post("/create", requestBody, nil)
	if err != nil {
		t.Fatalf("Client.Post() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Client.Post() status = %v, want %v", resp.StatusCode, http.StatusCreated)
	}
}

func TestClient_Put(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %v", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	requestBody := map[string]interface{}{
		"id":  123,
		"key": "updated",
	}

	resp, err := client.Put("/update/123", requestBody, nil)
	if err != nil {
		t.Fatalf("Client.Put() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Client.Put() status = %v, want %v", resp.StatusCode, http.StatusOK)
	}
}

func TestClient_Delete(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %v", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	resp, err := client.Delete("/delete/123", nil)
	if err != nil {
		t.Fatalf("Client.Delete() error = %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Client.Delete() status = %v, want %v", resp.StatusCode, http.StatusNoContent)
	}
}

func TestParseJSONResponse(t *testing.T) {
	// Create test response
	responseBody := `{"message": "success", "id": 123}`

	// Create with proper reader
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to get test response: %v", err)
	}

	var result map[string]interface{}
	err = ParseJSONResponse(resp, &result)
	if err != nil {
		t.Fatalf("ParseJSONResponse() error = %v", err)
	}

	if result["message"] != "success" {
		t.Errorf("ParseJSONResponse() message = %v, want success", result["message"])
	}

	if int(result["id"].(float64)) != 123 {
		t.Errorf("ParseJSONResponse() id = %v, want 123", result["id"])
	}
}

func TestWriteJSONResponse(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]interface{}{
		"message": "success",
		"id":      123,
	}

	err := WriteJSONResponse(w, http.StatusOK, data)
	if err != nil {
		t.Fatalf("WriteJSONResponse() error = %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("WriteJSONResponse() status = %v, want %v", w.Code, http.StatusOK)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("WriteJSONResponse() Content-Type = %v, want application/json", w.Header().Get("Content-Type"))
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if result["message"] != "success" {
		t.Errorf("WriteJSONResponse() message = %v, want success", result["message"])
	}
}

func TestWriteErrorResponse(t *testing.T) {
	w := httptest.NewRecorder()

	err := WriteErrorResponse(w, http.StatusBadRequest, "Invalid input")
	if err != nil {
		t.Fatalf("WriteErrorResponse() error = %v", err)
	}

	if w.Code != http.StatusBadRequest {
		t.Errorf("WriteErrorResponse() status = %v, want %v", w.Code, http.StatusBadRequest)
	}

	var result map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal error response: %v", err)
	}

	errorObj := result["error"].(map[string]interface{})
	if errorObj["message"] != "Invalid input" {
		t.Errorf("WriteErrorResponse() message = %v, want 'Invalid input'", errorObj["message"])
	}

	if int(errorObj["code"].(float64)) != http.StatusBadRequest {
		t.Errorf("WriteErrorResponse() code = %v, want %v", errorObj["code"], http.StatusBadRequest)
	}
}
