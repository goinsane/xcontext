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
	if period >= timeout || period < cancel {
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

func TestMultiContext(t *testing.T) {
	var timeout, timeout1, timeout2, cancel, period time.Duration

	timeout, timeout1, timeout2, cancel = 1*time.Second, 2*time.Second, 3*time.Second, 0
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testMultiContext(timeout, timeout1, timeout2, cancel)
	if period < timeout || period >= timeout1 {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}

	timeout, timeout1, timeout2, cancel = 3*time.Second, 2*time.Second, 1*time.Second, 0
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testMultiContext(timeout, timeout1, timeout2, cancel)
	if period < timeout2 || period >= timeout1 {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}

	timeout, timeout1, timeout2, cancel = 3*time.Second, 2*time.Second, 2*time.Second, 1*time.Second
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testMultiContext(timeout, timeout1, timeout2, cancel)
	if period >= timeout1 || period < cancel {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}
}

func testMultiContext(timeout, timeout1, timeout2, cancel time.Duration) time.Duration {
	startTm := time.Now()
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()
	ctx1, ctxCancel1 := context.WithTimeout(context.Background(), timeout1)
	defer ctxCancel1()
	ctx2, ctxCancel2 := context.WithTimeout(context.Background(), timeout2)
	defer ctxCancel2()
	newCtx, newCtxCancel := MultiContext(ctx, ctx1, ctx2)
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

func TestWaitContext(t *testing.T) {
	var timeout, timeout1, timeout2, cancel, period time.Duration

	timeout, timeout1, timeout2, cancel = 1*time.Second, 2*time.Second, 3*time.Second, 0
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testWaitContext(timeout, timeout1, timeout2, cancel)
	if period < timeout || period >= timeout1 {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}

	timeout, timeout1, timeout2, cancel = 3*time.Second, 2*time.Second, 1*time.Second, 0
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testWaitContext(timeout, timeout1, timeout2, cancel)
	if period < timeout1 || period >= timeout {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}

	timeout, timeout1, timeout2, cancel = 3*time.Second, 2*time.Second, 2*time.Second, 1*time.Second
	t.Logf("test with timeout=%v timeout1=%v timeout2=%v cancel=%v", timeout, timeout1, timeout2, cancel)
	period = testWaitContext(timeout, timeout1, timeout2, cancel)
	if period >= timeout1 || period < cancel {
		t.Errorf("test failed with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	} else {
		t.Logf("test ok with timeout=%v timeout1=%v timeout2=%v cancel=%v period=%v", timeout, timeout1, timeout2, cancel, period)
	}
}

func testWaitContext(timeout, timeout1, timeout2, cancel time.Duration) time.Duration {
	startTm := time.Now()
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)
	defer ctxCancel()
	ctx1, ctxCancel1 := context.WithTimeout(context.Background(), timeout1)
	defer ctxCancel1()
	ctx2, ctxCancel2 := context.WithTimeout(context.Background(), timeout2)
	defer ctxCancel2()
	newCtx, newCtxCancel := WaitContext(ctx, ctx1, ctx2)
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
