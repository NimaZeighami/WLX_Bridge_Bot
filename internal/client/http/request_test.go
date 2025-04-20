package http_test

import (
	customhttp "bridgebot/internal/client/http"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testResponse struct {
	Message string `json:"message"`
}

func TestGet_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, GET!"}`))
	}))
	defer ts.Close()

	headers := map[string]string{}
	query := map[string]string{}

	resp, err := customhttp.Get[testResponse](context.Background(), ts.URL, headers, query)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Hello, GET!", resp.Message)
}

func TestPost_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello, POST!"}`))
	}))
	defer ts.Close()

	headers := map[string]string{}
	body := map[string]string{"dummy": "data"}

	resp, err := customhttp.Post[testResponse](context.Background(), ts.URL, headers, body)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Hello, POST!", resp.Message)
}

func TestGet_InvalidURL(t *testing.T) {
	_, err := customhttp.Get[testResponse](context.Background(), "http://%%invalid-url", nil, nil)
	assert.Error(t, err)
}

func TestPost_InvalidURL(t *testing.T) {
	_, err := customhttp.Post[testResponse](context.Background(), "http://%%invalid-url", nil, nil)
	assert.Error(t, err)
}
