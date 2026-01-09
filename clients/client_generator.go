package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
)

// ClientGenerator generates REST clients with resty
type ClientGenerator struct {
	baseURL    string
	timeout    time.Duration
	debug      bool
	middleware []func(*resty.Client) *resty.Client
}

// GeneratorOption configures the ClientGenerator
type GeneratorOption func(*ClientGenerator)

// WithBaseURL sets the base URL for generated clients
func WithBaseURL(baseURL string) GeneratorOption {
	return func(cg *ClientGenerator) {
		cg.baseURL = baseURL
	}
}

// WithTimeout sets the timeout for generated clients
func WithTimeout(timeout time.Duration) GeneratorOption {
	return func(cg *ClientGenerator) {
		cg.timeout = timeout
	}
}

// WithDebug enables debug mode for generated clients
func WithDebug(debug bool) GeneratorOption {
	return func(cg *ClientGenerator) {
		cg.debug = debug
	}
}

// WithMiddleware adds custom middleware to generated clients
func WithMiddleware(middleware func(*resty.Client) *resty.Client) GeneratorOption {
	return func(cg *ClientGenerator) {
		cg.middleware = append(cg.middleware, middleware)
	}
}

// NewClientGenerator creates a new REST client generator
func NewClientGenerator(opts ...GeneratorOption) *ClientGenerator {
	cg := &ClientGenerator{
		baseURL:    "http://localhost:8080/api",
		timeout:    30 * time.Second,
		debug:      false,
		middleware: make([]func(*resty.Client) *resty.Client, 0),
	}

	for _, opt := range opts {
		opt(cg)
	}

	return cg
}

// GenerateClient creates a new resty client with all configurations
func (cg *ClientGenerator) GenerateClient(name string) *resty.Client {
	client := resty.New()

	// Basic configuration
	client.SetBaseURL(cg.baseURL).
		SetTimeout(cg.timeout).
		SetHeader("User-Agent", fmt.Sprintf("Generated-Client-%s/1.0", name)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	// Debug mode
	if cg.debug {
		client.SetDebug(true)
	}

	// Retry configuration
	client.SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(r *resty.Response, err error) bool {
			// Retry on network errors or 5xx status codes
			return r.StatusCode() >= 500 || err != nil
		})

	// Request/Response middleware
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// Add timestamp to all requests
		req.SetHeader("X-Request-Time", time.Now().UTC().Format(time.RFC3339))
		return nil
	})

	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// Log response time
		if cg.debug {
			fmt.Printf("[%s] %s %s - %d (%s)\n",
				name,
				resp.Request.Method,
				resp.Request.URL,
				resp.StatusCode(),
				resp.Time())
		}
		return nil
	})

	client.OnError(func(req *resty.Request, err error) {
		if cg.debug {
			fmt.Printf("[%s] Request Error: %v\n", name, err)
		}
	})

	// Apply custom middleware
	for _, middleware := range cg.middleware {
		client = middleware(client)
	}

	return client
}

// AuthMiddleware adds authentication headers
func AuthMiddleware(token string) func(*resty.Client) *resty.Client {
	return func(client *resty.Client) *resty.Client {
		client.SetAuthToken(token)
		return client
	}
}

// JSONResponseMiddleware handles JSON responses consistently
func JSONResponseMiddleware() func(*resty.Client) *resty.Client {
	return func(client *resty.Client) *resty.Client {
		client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// Ensure we received JSON
			contentType := resp.Header().Get("Content-Type")
			if contentType != "" && contentType != "application/json" {
				return fmt.Errorf("expected JSON response, got %s", contentType)
			}
			return nil
		})
		return client
	}
}

// RateLimitMiddleware adds rate limiting headers
func RateLimitMiddleware(requestsPerSecond int) func(*resty.Client) *resty.Client {
	return func(client *resty.Client) *resty.Client {
		client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			req.SetHeader("X-RateLimit-Limit", fmt.Sprintf("%d", requestsPerSecond))
			return nil
		})
		return client
	}
}

// HealthCheckMiddleware adds health check capability
func HealthCheckMiddleware(endpoint string) func(*resty.Client) *resty.Client {
	return func(client *resty.Client) *resty.Client {
		// Add health check method to client
		client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			// Add health check headers if endpoint matches
			if req.URL == endpoint {
				req.SetHeader("X-Health-Check", "true")
			}
			return nil
		})
		return client
	}
}

// GenericAPIClient provides common REST operations
type GenericAPIClient struct {
	client   *resty.Client
	basePath string
}

// NewGenericAPIClient creates a generic API client
func NewGenericAPIClient(generator *ClientGenerator, basePath string) *GenericAPIClient {
	client := generator.GenerateClient("Generic")
	return &GenericAPIClient{
		client:   client,
		basePath: basePath,
	}
}

// Get performs a GET request
func (api *GenericAPIClient) Get(path string, result interface{}) error {
	resp, err := api.client.R().
		SetResult(result).
		Get(api.basePath + path)

	if err != nil {
		return fmt.Errorf("GET request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("GET request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Post performs a POST request
func (api *GenericAPIClient) Post(path string, body interface{}, result interface{}) error {
	request := api.client.R()

	if body != nil {
		request = request.SetBody(body)
	}

	if result != nil {
		request = request.SetResult(result)
	}

	resp, err := request.Post(api.basePath + path)

	if err != nil {
		return fmt.Errorf("POST request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("POST request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Put performs a PUT request
func (api *GenericAPIClient) Put(path string, body interface{}, result interface{}) error {
	request := api.client.R()

	if body != nil {
		request = request.SetBody(body)
	}

	if result != nil {
		request = request.SetResult(result)
	}

	resp, err := request.Put(api.basePath + path)

	if err != nil {
		return fmt.Errorf("PUT request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("PUT request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// Delete performs a DELETE request
func (api *GenericAPIClient) Delete(path string) error {
	resp, err := api.client.R().Delete(api.basePath + path)

	if err != nil {
		return fmt.Errorf("DELETE request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return fmt.Errorf("DELETE request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return nil
}

// HealthCheck performs a health check
func (api *GenericAPIClient) HealthCheck() (map[string]interface{}, error) {
	var result map[string]interface{}

	resp, err := api.client.R().
		SetResult(&result).
		Get("/health")

	if err != nil {
		return nil, fmt.Errorf("health check failed: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("health check failed with status %d", resp.StatusCode())
	}

	return result, nil
}

// GetRawJSON performs a GET request and returns raw JSON
func (api *GenericAPIClient) GetRawJSON(path string) (json.RawMessage, error) {
	resp, err := api.client.R().Get(api.basePath + path)

	if err != nil {
		return nil, fmt.Errorf("GET request failed: %v", err)
	}

	if resp.StatusCode() >= 400 {
		return nil, fmt.Errorf("GET request failed with status %d: %s", resp.StatusCode(), resp.String())
	}

	return json.RawMessage(resp.Body()), nil
}
