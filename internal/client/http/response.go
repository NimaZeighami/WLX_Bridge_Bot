package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "bridgebot/internal/utils/logger"
)

// Response represents a structured HTTP response.
type Response struct {
	StatusCode int
	Body       []byte
	Headers    http.Header
}

// ParseJSON decodes the response body into the given target struct.
func (r *Response) ParseJSON(target interface{}) error {
	if err := json.Unmarshal(r.Body, target); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}
	return nil
}

// HandleResponse processes the HTTP response and returns a structured Response object.
func HandleResponse(res *http.Response) (*Response, error) {
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Errorf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	response := &Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Headers:    res.Header,
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Errorf("HTTP Error: %d - %s", res.StatusCode, string(body))
		return response, fmt.Errorf("HTTP error: received status code %d, response: %s", res.StatusCode, string(body))
	}

	//? Uncomment the line below to log successful HTTP responses 
	// log.Infof("Successful HTTP Response [%d]: %s", res.StatusCode, string(body))
	return response, nil
}
