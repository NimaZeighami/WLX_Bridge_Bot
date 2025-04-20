package http_test

import (
	customhttp "bridgebot/internal/client/http"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseJSON_Success(t *testing.T) {
	res := &customhttp.Response{
		StatusCode: 200,
		Body:       []byte(`{"message": "Parsed successfully"}`),
	}

	type msg struct {
		Message string `json:"message"`
	}
	var result msg
	err := res.ParseJSON(&result)
	assert.NoError(t, err)
	assert.Equal(t, "Parsed successfully", result.Message)
}

func TestHandleResponse_Success(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "hello"}`))
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer res.Body.Close()

	parsed, err := customhttp.HandleResponse(res)
	assert.NoError(t, err)
	assert.Equal(t, 200, parsed.StatusCode)
	assert.Contains(t, string(parsed.Body), "hello")
}

func TestHandleResponse_HTTPError(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`Bad request`))
	}

	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	assert.NoError(t, err)
	defer res.Body.Close()

	parsed, err := customhttp.HandleResponse(res)
	assert.Error(t, err)
	assert.Equal(t, 400, parsed.StatusCode)
	assert.Contains(t, string(parsed.Body), "Bad request")
}

func TestHandleResponse_ReadError(t *testing.T) {
	res := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(brokenReader{}),
	}
	_, err := customhttp.HandleResponse(res)
	assert.Error(t, err)
}

type brokenReader struct{}

func (brokenReader) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (brokenReader) Close() error {
	return nil
}