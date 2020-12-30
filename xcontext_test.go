package xcontext

import (
	"context"
	"testing"
	"time"
)

func TestDelayContext(t *testing.T) {
	var timeout, delay, cancel, period time.Duration

	timeout, delay, cancel = 1*time.Second, 2*time.Second, 0
	t.Logf("test with timeout=%v delay=%v cancel=%v", timeout, delay, cancel)
	period = testDelayContext(timeout, delay, cancel)
	if period < delay {
		t.Errorf("test failed with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	}

	timeout, delay, cancel = 2*time.Second, 1*time.Second, 0
	t.Logf("test with timeout=%v delay=%v cancel=%v", timeout, delay, cancel)
	period = testDelayContext(timeout, delay, cancel)
	if period < timeout {
		t.Errorf("test failed with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	}

	timeout, delay, cancel = 2*time.Second, 2*time.Second, 1*time.Second
	t.Logf("test with timeout=%v delay=%v cancel=%v", timeout, delay, cancel)
	period = testDelayContext(timeout, delay, cancel)
	if period >= timeout {
		t.Errorf("test failed with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v delay=%v cancel=%v period=%v", timeout, delay, cancel, period)
	}
}

func testDelayContext(timeout, delay, cancel time.Duration) time.Duration {
	startTm := time.Now()
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()
	newCtx, newCtxCancel := DelayContext(ctx, delay)
	defer newCtxCancel()
	if cancel > 0 {
		go func() {
			<-time.After(cancel)
			newCtxCancel()
		}()
	}
	<-newCtx.Done()
	return time.Now().Sub(startTm)
}
