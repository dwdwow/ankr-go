package ankr

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// HTTPClient represents the HTTP client for Ankr Advanced API
type HTTPClient struct {
	uri         string
	httpClient  *http.Client
	rateLimiter *RateLimiter
}

type HTTPClientConfig struct {
	APIKey          string
	OnLimitExceeded RateLimitBehavior `default:"block"`
}

// NewHTTPClient creates a new HTTP client with the given configuration
func NewHTTPClient(config *HTTPClientConfig) *HTTPClient {
	if config == nil {
		config = &HTTPClientConfig{}
	}

	// Apply defaults to config
	_ = ApplyDefaults(config)

	// Create rate limiter
	rateLimiter, _ := NewRateLimiter(1000, time.Minute, config.OnLimitExceeded)

	// Create HTTP client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			IdleConnTimeout: 90 * time.Second,
		},
	}

	return &HTTPClient{
		uri:         "https://rpc.ankr.com/multichain/" + config.APIKey,
		httpClient:  httpClient,
		rateLimiter: rateLimiter,
	}
}

// post makes a JSON-RPC post request and returns the result with generic type
func post[T any](client *HTTPClient, ctx context.Context, method string, params any) (T, error) {
	var result T

	// Rate limiting
	acquired, err := client.rateLimiter.Acquire(ctx, 1, nil)
	if err != nil {
		return result, fmt.Errorf("rate limit error: %w", err)
	}
	if !acquired {
		return result, ErrRateLimitExceeded
	}

	// Create JSON-RPC request
	request := RPCReqBody{
		ID:      1,
		JSONRPC: JSONRPC,
		Method:  method,
		Params:  params,
	}

	// Marshal request body
	requestBody, err := json.Marshal(request)
	if err != nil {
		return result, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", client.uri, bytes.NewReader(requestBody))
	if err != nil {
		return result, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, fmt.Errorf("failed to read response: %w", err)
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(body))
	}

	// Parse JSON response
	var apiResponse RPCRespBody[T]
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return result, fmt.Errorf("failed to parse response: %w", err)
	}

	return apiResponse.Result, nil
}

func postWithRetries[T any](client *HTTPClient, ctx context.Context, method string, params any, retries int) (result T, err error) {
	for range retries {
		result, err = post[T](client, ctx, method, params)
		if err == nil {
			return
		}
		slog.Error("ankr: failed to post", "error", err)
		time.Sleep(time.Second)
		continue

	}
	err = fmt.Errorf("ankr: failed to post after %d retries, last error: %w", retries, err)
	return
}

type paginationReqData interface {
	PageToken() string
	SetPageToken(pageToken string)
}

type paginationRespData interface {
	NextPageToken() string
}

func postPages[Req paginationReqData, Resp paginationRespData](client *HTTPClient, ctx context.Context, method string, params Req, retries int) (results []Resp, err error) {
	for {
		var pageToken string
		if len(results) > 0 {
			pageToken = results[len(results)-1].NextPageToken()
			if pageToken == "" {
				return
			}
			params.SetPageToken(pageToken)
		}
		var result Resp
		result, err = postWithRetries[Resp](client, ctx, method, params, retries)
		if err != nil {
			return
		}
		results = append(results, result)
	}
}
