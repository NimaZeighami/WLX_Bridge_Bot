package http

import (
	log "bridgebot/internal/utils/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Get sends an HTTP GET request with support for headers, query params, context, retries, and JSON decoding.
func Get[T any](ctx context.Context, baseURL string, headers map[string]string, queryParams map[string]string) (*T, error) {
	finalURL, err := addQueryParams(baseURL, queryParams)
	if err != nil {
		return nil, fmt.Errorf("error parsing URL with query params: %v", err)
	}

	// log.Infof("Sending GET request to URL: %s", finalURL)

	req, err := http.NewRequestWithContext(ctx, "GET", finalURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating GET request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Infof("Attempt %d to call %s", attempt, finalURL)
		res, err := client.Do(req)
		if err != nil {
			log.Errorf("Request error: %v", err)
			if attempt < maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("error sending GET request: %v", err)
		}

		parsedResponse, err := HandleResponse(res)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, err
		}

		var result T
		if err := parsedResponse.ParseJSON(&result); err != nil {
			return nil, err
		}

		// log.Infof("GET request successful: %s", finalURL)
		return &result, nil
	}

	return nil, fmt.Errorf("failed to complete GET request after %d attempts", maxRetries)
}

// Post sends an HTTP POST request with support for headers, context, retries, JSON encoding, and response decoding.
func Post[T any](ctx context.Context, url string, headers map[string]string, body interface{}) (*T, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error encoding JSON body: %v", err)
	}

	// log.Infof("Sending POST request to URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, fmt.Errorf("error creating POST request: %v", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	const maxRetries = 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Infof("Attempt %d to call %s", attempt, url)
		res, err := client.Do(req)
		if err != nil {
			log.Errorf("Request error: %v", err)
			if attempt < maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, fmt.Errorf("error sending POST request: %v", err)
		}

		parsedResponse, err := HandleResponse(res)
		if err != nil {
			if attempt < maxRetries {
				time.Sleep(2 * time.Second)
				continue
			}
			return nil, err
		}

		var result T
		if err := parsedResponse.ParseJSON(&result); err != nil {
			return nil, err
		}

		// log.Infof("POST request successful: %s", url)
		return &result, nil
	}

	return nil, fmt.Errorf("failed to complete POST request after %d attempts", maxRetries)
}

// addQueryParams adds query parameters to a base URL.
func addQueryParams(baseURL string, params map[string]string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for key, value := range params {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

// package http

// import (
// 	log "bridgebot/internal/utils/logger"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"time"
// )

// // Get sends an HTTP GET request with support for headers, query params, context, retries, and JSON decoding.
// func Get[T any](ctx context.Context, baseURL string, headers map[string]string, queryParams map[string]string) (*T, error) {
// 	finalURL, err := addQueryParams(baseURL, queryParams)
// 	if err != nil {
// 		return nil, fmt.Errorf("error parsing URL with query params: %v", err)
// 	}

// 	log.Infof("Sending GET request to URL: %s", finalURL)

// 	req, err := http.NewRequestWithContext(ctx, "GET", finalURL, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating GET request: %v", err)
// 	}

// 	for key, value := range headers {
// 		req.Header.Set(key, value)
// 	}

// 	client := &http.Client{
// 		Timeout: 5 * time.Second,
// 	}

// 	const maxRetries = 3
// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		log.Infof("Attempt %d to call %s", attempt, finalURL)
// 		res, err := client.Do(req)
// 		if err != nil {
// 			log.Errorf("Request error: %v", err)
// 			if attempt < maxRetries {
// 				time.Sleep(time.Second * 2)
// 				continue
// 			}
// 			return nil, fmt.Errorf("error sending GET request: %v", err)
// 		}
// 		defer res.Body.Close()

// 		if res.StatusCode < 200 || res.StatusCode >= 300 {
// 			bodyBytes, _ := io.ReadAll(res.Body)
// 			log.Errorf("Received non-success status code: %d with body: %s", res.StatusCode, string(bodyBytes))
// 			if attempt < maxRetries {
// 				time.Sleep(time.Second * 2)
// 				continue
// 			}
// 			return nil, fmt.Errorf("error: received status code %d, response: %s", res.StatusCode, string(bodyBytes))
// 		}

// 		var result T
// 		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 			log.Errorf("JSON parsing error: %v", err)
// 			return nil, fmt.Errorf("error parsing json: %v", err)
// 		}

// 		log.Infof("Request successful for URL: %s", finalURL)
// 		return &result, nil
// 	}

// 	return nil, fmt.Errorf("failed to complete GET request after %d attempts", maxRetries)
// }

// // Post sends an HTTP POST request with support for headers, context, retries, JSON encoding, and response decoding.
// func Post[T any](ctx context.Context, url string, headers map[string]string, body interface{}) (*T, error) {
// 	jsonBody, err := json.Marshal(body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error encoding JSON body: %v", err)
// 	}

// 	log.Infof("Sending POST request to URL: %s", url)

// 	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonBody)))
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating POST request: %v", err)
// 	}

// 	for key, value := range headers {
// 		req.Header.Set(key, value)
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	client := &http.Client{
// 		Timeout: 5 * time.Second,
// 	}

// 	const maxRetries = 3
// 	for attempt := 1; attempt <= maxRetries; attempt++ {
// 		log.Infof("Attempt %d to call %s", attempt, url)
// 		res, err := client.Do(req)
// 		if err != nil {
// 			log.Errorf("Request error: %v", err)
// 			if attempt < maxRetries {
// 				time.Sleep(time.Second * 2)
// 				continue
// 			}
// 			return nil, fmt.Errorf("error sending POST request: %v", err)
// 		}
// 		defer res.Body.Close()

// 		if res.StatusCode < 200 || res.StatusCode >= 300 {
// 			bodyBytes, _ := io.ReadAll(res.Body)
// 			log.Errorf("Received non-success status code: %d with body: %s", res.StatusCode, string(bodyBytes))
// 			if attempt < maxRetries {
// 				time.Sleep(time.Second * 2)
// 				continue
// 			}
// 			return nil, fmt.Errorf("error: received status code %d, response: %s", res.StatusCode, string(bodyBytes))
// 		}

// 		var result T
// 		if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
// 			log.Errorf("JSON parsing error: %v", err)
// 			return nil, fmt.Errorf("error parsing json: %v", err)
// 		}

// 		log.Infof("POST request successful for URL: %s", url)
// 		return &result, nil
// 	}

// 	return nil, fmt.Errorf("failed to complete POST request after %d attempts", maxRetries)
// }

// // addQueryParams adds query parameters to a base URL.
// func addQueryParams(baseURL string, params map[string]string) (string, error) {
// 	u, err := url.Parse(baseURL)
// 	if err != nil {
// 		return "", err
// 	}
// 	q := u.Query()
// 	for key, value := range params {
// 		q.Set(key, value)
// 	}
// 	u.RawQuery = q.Encode()
// 	return u.String(), nil
// }

// ! If speed mattes, consider using gnet or fasthttp instead of the standard library.
// ! If you need to send a lot of requests, consider using a connection pool.
// ! If you need to send a lot of requests to the same host, consider using a connection pool.
