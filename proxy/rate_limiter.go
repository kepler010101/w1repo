package main

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	mutex     sync.Mutex
	requests  map[string][]time.Time
	limit     int
	windowSec int
}

func NewSlidingWindowLimiter(limit int, windowSec int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		requests:  make(map[string][]time.Time),
		limit:     limit,
		windowSec: windowSec,
	}
}

func (l *SlidingWindowLimiter) Allow(key string) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	windowStart := now.Add(-time.Duration(l.windowSec) * time.Second)

	requests, exists := l.requests[key]
	if !exists {
		l.requests[key] = []time.Time{now}
		return true
	}

	validRequests := make([]time.Time, 0)
	for _, t := range requests {
		if t.After(windowStart) {
			validRequests = append(validRequests, t)
		}
	}

	if len(validRequests) < l.limit {
		validRequests = append(validRequests, now)
		l.requests[key] = validRequests
		return true
	}

	return false
}
