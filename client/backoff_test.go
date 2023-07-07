package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoffRoundTripper(t *testing.T) {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		count++
		if count <= 2 {
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(http.StatusOK)
		}
	}))

	defer server.Close()

	client := &http.Client{
		Transport: &ExponentialBackoffRoundTripper{
			RoundTripper: http.DefaultTransport,
			MaxRetries:   3,
			BackoffTime:  1 * time.Second,
		},
	}

	// Add delay to simulate real network call
	time.Sleep(2 * time.Second)

	resp, err := client.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
