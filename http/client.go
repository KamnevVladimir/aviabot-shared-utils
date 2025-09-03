package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client provides HTTP client utilities
type Client struct {
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new HTTP client with default settings
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// NewClientWithTimeout creates a new HTTP client with custom timeout
func NewClientWithTimeout(baseURL string, timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL: baseURL,
	}
}

// Get performs a GET request
func (c *Client) Get(endpoint string, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GET request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.httpClient.Do(req)
}

// Post performs a POST request with JSON body
func (c *Client) Post(endpoint string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(endpoint)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(http.MethodPost, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create POST request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.setHeaders(req, headers)

	return c.httpClient.Do(req)
}

// Put performs a PUT request with JSON body
func (c *Client) Put(endpoint string, body interface{}, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(endpoint)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal JSON body: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(http.MethodPut, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create PUT request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	c.setHeaders(req, headers)

	return c.httpClient.Do(req)
}

// Delete performs a DELETE request
func (c *Client) Delete(endpoint string, headers map[string]string) (*http.Response, error) {
	url := c.buildURL(endpoint)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create DELETE request: %w", err)
	}

	c.setHeaders(req, headers)
	return c.httpClient.Do(req)
}

// buildURL constructs the full URL
func (c *Client) buildURL(endpoint string) string {
	if c.baseURL == "" {
		return endpoint
	}

	baseURL, err := url.Parse(c.baseURL)
	if err != nil {
		return endpoint
	}

	endpointURL, err := url.Parse(endpoint)
	if err != nil {
		return c.baseURL + endpoint
	}

	return baseURL.ResolveReference(endpointURL).String()
}

// setHeaders sets request headers
func (c *Client) setHeaders(req *http.Request, headers map[string]string) {
	for key, value := range headers {
		req.Header.Set(key, value)
	}
}

// ParseJSONResponse parses JSON response into provided struct
func ParseJSONResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return nil
}

// WriteJSONResponse writes JSON response to http.ResponseWriter
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data == nil {
		return nil
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}

	return nil
}

// WriteErrorResponse writes error response in standard format
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message string) error {
	errorResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
			"code":    statusCode,
		},
	}

	return WriteJSONResponse(w, statusCode, errorResponse)
}
