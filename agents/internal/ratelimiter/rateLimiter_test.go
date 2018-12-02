package ratelimiter

import (
	"context"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func Test(t *testing.T) {

	rate := New(1.0/2.0 - 2.0/10.0)

	ctx := context.Background()
	for i := 0; i < 4; i++ {
		rate.Take(ctx)
		t.Log(time.Now())
	}
}

func Test2(t *testing.T) {
	rate := New(1.0)
	logrus.SetLevel(logrus.DebugLevel)
	ctx := context.Background()
	for i := 0; i < 4; i++ {
		logrus.Debug("before call:", time.Now())
		rate.TakeMore(ctx, 1, 5)
		logrus.Debug("after call: ", time.Now())
	}
}
