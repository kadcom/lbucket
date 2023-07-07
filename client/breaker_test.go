package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCircuitBreakerRoundTripper(t *testing.T) {
	count := 0
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		count++
		if count <= 3 {
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			rw.WriteHeader(http.StatusOK)
		}
	}))

	defer server.Close()

	cb := &CircuitBreakerRoundTripper{
		RoundTripper: http.DefaultTransport,
		FailureCount: 0,
		Open:         false,
		ResetTimeout: 1 * time.Second,
	}

	client := &http.Client{
		Transport: cb,
	}

	// First 3 requests should fail and open the circuit breaker
	for i := 0; i < 3; i++ {
		resp, _ := client.Get(server.URL)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	}

	// The next request should fail because the circuit breaker is open
	resp, err := client.Get(server.URL)
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.True(t, cb.Open)

	// Wait for the circuit breaker to reset
	time.Sleep(2 * time.Second)

	// The next request should succeed because the circuit breaker has reset
	resp, err = client.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.False(t, cb.Open)
}
