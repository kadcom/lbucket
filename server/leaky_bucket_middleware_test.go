package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestLeakyBucketMiddleware(t *testing.T) {
	capacity := 1
	rate := 1 * time.Second
	bucket := CreateLeakyBucket(capacity, rate)

	r := chi.NewRouter()
	r.Use(LeakyBucketMiddleware(&bucket))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	server := httptest.NewServer(r)
	defer server.Close()

	// First request should be allowed.
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Second request should be rate limited.
	resp, err = http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusTooManyRequests, resp.StatusCode)

	// After waiting for the refill duration, we should be able to make another request.
	time.Sleep(rate)
	resp, err = http.Get(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
