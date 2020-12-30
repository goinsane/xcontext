package xcontext

import (
	"context"
	"testing"
	"time"
)

func TestDelayContext(t *testing.T) {
	var timeout, delay, period time.Duration

	timeout, delay = 1*time.Second, 2*time.Second
	t.Logf("test with timeout=%v delay=%v", timeout, delay)
	period = testDelayContext(timeout, delay)
	if period < delay {
		t.Errorf("test failed with timeout=%v delay=%v period=%v", timeout, delay, period)
	} else {
		t.Logf("test ok with timeout=%v delay=%v period=%v", timeout, delay, period)
	}

	timeout, delay = 2*time.Second, 1*time.Second
	t.Logf("test with timeout=%v delay=%v", timeout, delay)
	period = testDelayContext(timeout, delay)
	if period < timeout {
		t.Errorf("test failed with timeout=%v delay=%v period=%v", timeout, delay, period)
	} else {
		t.Logf("test ok with timeout=%v delay=%v period=%v", timeout, delay, period)
	}
}

func testDelayContext(timeout, delay time.Duration) time.Duration {
	startTm := time.Now()
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()
	newCtx, newCtxCancel := DelayContext(ctx, delay)
	defer newCtxCancel()
	<-newCtx.Done()
	return time.Now().Sub(startTm)
}
