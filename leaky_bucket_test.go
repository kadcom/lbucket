package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLeakyBucket_TryTake(t *testing.T) {
	capacity := 2
	rate := 1 * time.Second
	bucket := CreateLeakyBucket(capacity, rate)

	// Should be able to take twice immediately.
	assert.True(t, bucket.TryTake(), "Expected TryTake to return true, got false")
	assert.True(t, bucket.TryTake(), "Expected TryTake to return true, got false")

	// But not a third time.
	assert.False(t, bucket.TryTake(), "Expected TryTake to return false, got true")

	// After waiting for the refill duration, we should be able to take again.
	time.Sleep(rate)
	assert.True(t, bucket.TryTake(), "Expected TryTake to return true, got false")
}
