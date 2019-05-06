package internal

import (
	"context"
	"github.com/nathanhack/sibyl/core"
	"time"
)

func tomorrowAt6AM() time.Time {
	return core.NewDateTypeFromTime(time.Now()).
		AddDate(0, 0, 1).
		Time().Add(6 * time.Hour)
}

//areWeDone is a nonblocking context.Done() check
func areWeDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
	}
	return false
}
