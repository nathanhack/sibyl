package flowLimiter

import "context"

type FlowLimiter struct {
	c chan bool
}

func New(maxTakes int) *FlowLimiter {
	channel := make(chan bool, maxTakes)
	for i := 0; i < maxTakes; i++ {
		channel <- true
	}
	return &FlowLimiter{channel}
}

func (fl *FlowLimiter) Take(ctx context.Context) {
	//the intent is to block
	// until the limiter allows it
	select {
	case <-ctx.Done():
	case <-fl.c:
	}
}

func (fl *FlowLimiter) Return() {
	//we make sure this never blocks
	select {
	case fl.c <- true:
	default:
	}
}
