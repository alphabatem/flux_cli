package httpclient

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/alphabatem/flux_cli/dto"
)

type Client struct {
	BaseURL    string
	APIKey     string
	HeaderName string // Header name for API key (e.g. "Authorization", "X-API-KEY")
	HTTPClient *http.Client
}

func New(baseURL, apiKey, headerName string) *Client {
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		APIKey:     apiKey,
		HeaderName: headerName,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get performs a GET request and unmarshals the response into target.
func (c *Client) Get(path string, target interface{}) error {
	req, err := http.NewRequest("GET", c.BaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	c.applyHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, target)
}

// Post performs a POST request with a JSON body and unmarshals the response.
func (c *Client) Post(path string, body io.Reader, target interface{}) error {
	req, err := http.NewRequest("POST", c.BaseURL+path, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.applyHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, target)
}

// PostRaw performs a POST and returns the raw response body.
func (c *Client) PostRaw(path string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest("POST", c.BaseURL+path, body)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.applyHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}

	return data, resp.StatusCode, nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string, target interface{}) error {
	req, err := http.NewRequest("DELETE", c.BaseURL+path, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	c.applyHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, target)
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, body io.Reader, target interface{}) error {
	req, err := http.NewRequest("PUT", c.BaseURL+path, body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	c.applyHeaders(req)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, target)
}

func (c *Client) applyHeaders(req *http.Request) {
	if c.APIKey != "" && c.HeaderName != "" {
		req.Header.Set(c.HeaderName, c.APIKey)
	}
	req.Header.Set("User-Agent", "flux-cli/0.1.0")
}

func (c *Client) handleResponse(resp *http.Response, target interface{}) error {
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(data),
		}
	}

	if target != nil {
		if err := json.Unmarshal(data, target); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}

// APIError represents an HTTP error response from an API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error (HTTP %d): %s", e.StatusCode, e.Body)
}

// ExitCodeForError maps an API error to a CLI exit code.
func ExitCodeForError(err error) int {
	if apiErr, ok := err.(*APIError); ok {
		switch {
		case apiErr.StatusCode == 401 || apiErr.StatusCode == 403:
			return dto.ExitAuthError
		case apiErr.StatusCode == 503:
			return dto.ExitServiceDown
		default:
			return dto.ExitAPIError
		}
	}
	return dto.ExitGeneralError
}
