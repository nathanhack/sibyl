package ratelimiter

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

//TODO consider making this also perform rate limiting with priority, so more important events take priority; also consider using or making its ICD like golang.org/x/time/rate
type RateLimiter struct {
	sync.Mutex
	last            time.Time
	timePerRequests time.Duration
}

func New(requestsPerSecond float64) *RateLimiter {
	return &RateLimiter{
		timePerRequests: time.Duration(float64(time.Second) / requestsPerSecond),
	}
}

func (rl *RateLimiter) TimePerRequests() time.Duration {
	return rl.timePerRequests
}

func (t *RateLimiter) Take(ctx context.Context) error {
	return t.TakeMore(ctx, 1, 0)
}

func (t *RateLimiter) TakeMore(ctx context.Context, rateBlocksBefore, rateBlocksAfter int) error {
	t.Lock()
	defer t.Unlock()

	now := time.Now()

	// if this is our first request
	// set the last time to now
	if t.last.IsZero() {
		t.last = now.Add(time.Duration(float64(rateBlocksAfter) * float64(t.timePerRequests)))
		logrus.Debugf("RateLimiter(%v) released at %v", t.timePerRequests, now)
		return nil
	}

	// calculate how long to sleep for based on requestsPerSecond
	sleep := time.Duration(float64(rateBlocksBefore)*float64(t.timePerRequests)) - now.Sub(t.last)

	//if sleepFor is <= 0 we set last to
	// now and return

	if sleep <= 0 {
		t.last = now.Add(time.Duration(float64(rateBlocksBefore) * float64(t.timePerRequests)))
		logrus.Debugf("RateLimiter(%v) released at %v", t.timePerRequests, now)
		return nil
	}

	//else we do a context based sleep

	select {
	case <-ctx.Done():
		return fmt.Errorf("context canceled")
	case <-time.After(sleep):
	}

	now = time.Now()
	t.last = now.Add(time.Duration(float64(rateBlocksAfter) * float64(t.timePerRequests)))
	logrus.Debugf("RateLimiter(%v) released at %v", t.timePerRequests, now)
	return nil
}
