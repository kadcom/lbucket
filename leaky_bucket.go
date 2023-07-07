package main

import "time"

type LeakyBucket struct {
	capacity  int
	remaining int
	reset     time.Time
	rate      time.Duration
}

func CreateLeakyBucket(capacity int, rate time.Duration) LeakyBucket {
	return LeakyBucket{
		capacity:  capacity,
		rate:      rate,
		reset:     time.Now().Add(rate),
		remaining: capacity,
	}
}

func (bucket *LeakyBucket) Refill() {
	now := time.Now()

	if now.After(bucket.reset) {
		bucket.remaining = bucket.capacity
		bucket.reset = now.Add(bucket.rate)
	}
}

func (bucket *LeakyBucket) TryTake() bool {
	bucket.Refill()

	if bucket.remaining <= 0 {
		return false
	}

	bucket.remaining--
	return true
}
