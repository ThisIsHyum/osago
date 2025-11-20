package osago

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ThisIsHyum/osago/dto"
)

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http %d: %s", e.StatusCode, e.Message)
}

func NewHttpError(statusCode int, message string) error {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
	}
}

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(url string, timeout time.Duration) *Client {
	return &Client{
		baseURL:    url,
		httpClient: &http.Client{Timeout: timeout},
	}
}

func (c *Client) doReq(ctx context.Context, method, path string, body, v any) error {
	var req *http.Request
	var err error
	if body != nil {
		b, marshalErr := json.Marshal(body)
		if marshalErr != nil {
			return marshalErr
		}
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	}
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errorResponse dto.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("unable to decode error response: %w", err)
		}
		return NewHttpError(errorResponse.StatusCode, errorResponse.Error)
	}
	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("unable to decode response: %w", err)
		}
	}
	return nil
}
