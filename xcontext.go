// Package xcontext provides extended utilities for context.
package xcontext

import (
	"context"
	"sync"
	"time"
)

// DelayContext creates a new cancel context not inherited from ctx.
// It would be cancelled if ctx is done, after the delay started with the call of DelayContext is done.
// The new context doesn't inherit ctx.
// Calling the cancel function is not necessary.
func DelayContext(ctx context.Context, delay time.Duration) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(context.Background())
	delayTimer := time.NewTimer(delay)
	go func() {
		parentCtx := ctx
		parentCtxOk := context.Background()
		delayCh := delayTimer.C
		delayOkCh := make(<-chan time.Time)
		for done := false; !done; {
			select {
			case <-parentCtx.Done():
				parentCtx = parentCtxOk
				if delayCh == delayOkCh {
					done = true
				}
			case <-delayCh:
				delayCh = delayOkCh
				if parentCtx == parentCtxOk {
					done = true
				}
			case <-newCtx.Done():
				done = true
			}
		}
		newCtxCancel()
		delayTimer.Stop()
	}()
	return newCtx, newCtxCancel
}

// DelayContext2 is similar with DelayContext except not returns cancel function.
func DelayContext2(ctx context.Context, delay time.Duration) context.Context {
	newCtx, _ := DelayContext(ctx, delay)
	return newCtx
}

// DelayAfterContext creates a new cancel context not inherited from ctx.
// It would be cancelled after the delay started when ctx is done.
// The new context doesn't inherit ctx.
// Calling the cancel function is not necessary.
func DelayAfterContext(ctx context.Context, delay time.Duration) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(context.Background())
	delayTimer := time.NewTimer(0)
	if !delayTimer.Stop() {
		<-delayTimer.C
	}
	go func() {
		parentCtx := ctx
		parentCtxOk := context.Background()
		for done := false; !done; {
			select {
			case <-parentCtx.Done():
				parentCtx = parentCtxOk
				delayTimer.Reset(delay)
			case <-delayTimer.C:
				if parentCtx == parentCtxOk {
					done = true
				}
			case <-newCtx.Done():
				done = true
			}
		}
		newCtxCancel()
		delayTimer.Stop()
	}()
	return newCtx, newCtxCancel
}

// DelayAfterContext2 is similar with DelayAfterContext except not returns cancel function.
func DelayAfterContext2(ctx context.Context, delay time.Duration) context.Context {
	newCtx, _ := DelayAfterContext(ctx, delay)
	return newCtx
}

// MultiContext creates a new cancel context inherited from parent.
// It would be cancelled when at least one of the sub contexts is done.
// Calling the cancel function is not necessary.
func MultiContext(parent context.Context, subs ...context.Context) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(parent)
	for _, sub := range subs {
		go func(sub context.Context) {
			select {
			case <-sub.Done():
			case <-newCtx.Done():
			}
			newCtxCancel()
		}(sub)
	}
	return newCtx, newCtxCancel
}

// MultiContext2 is similar with MultiContext except not returns cancel function.
func MultiContext2(parent context.Context, subs ...context.Context) context.Context {
	newCtx, _ := MultiContext(parent, subs...)
	return newCtx
}

// WaitContext creates a new cancel context inherited from parent.
// It would be cancelled when all the sub contexts are done.
// Calling the cancel function is not necessary.
func WaitContext(parent context.Context, subs ...context.Context) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(parent)
	var wg sync.WaitGroup
	for _, sub := range subs {
		wg.Add(1)
		go func(sub context.Context) {
			select {
			case <-sub.Done():
			case <-newCtx.Done():
			}
			wg.Done()
		}(sub)
	}
	go func() {
		wg.Wait()
		newCtxCancel()
	}()
	return newCtx, newCtxCancel
}

// WaitContext2 is similar with WaitContext except not returns cancel function.
func WaitContext2(parent context.Context, subs ...context.Context) context.Context {
	newCtx, _ := WaitContext(parent, subs...)
	return newCtx
}

// Or creates a new cancel context to cancel when at least one of the contexts is done.
// The new context doesn't inherit any context in ctxs.
func Or(ctxs ...context.Context) (context.Context, context.CancelFunc) {
	return MultiContext(context.Background(), ctxs...)
}

// Or2 is similar with Or except not returns cancel function.
func Or2(ctxs ...context.Context) context.Context {
	return MultiContext2(context.Background(), ctxs...)
}

// And creates a new cancel context to cancel when all the contexts are done.
// The new context doesn't inherit any context in ctxs.
func And(ctxs ...context.Context) (context.Context, context.CancelFunc) {
	return WaitContext(context.Background(), ctxs...)
}

// And2 is similar with And except not returns cancel function.
func And2(ctxs ...context.Context) context.Context {
	return WaitContext2(context.Background(), ctxs...)
}

// AutoCancel cancels the underlying context with cancel function when it was done through parent, deadline, timeout or in any way.
// It returns ctx and cancel.
func AutoCancel(ctx context.Context, cancel context.CancelFunc) (context.Context, context.CancelFunc) {
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx, cancel
}

// WithCancel is similar with context.WithCancel.
// But it cancels the context when it was done through parent.
func WithCancel(parent context.Context) (context.Context, context.CancelFunc) {
	return AutoCancel(context.WithCancel(parent))
}

// WithCancel2 is similar with WithCancel.
// It returns only a new context inherited from parent.
func WithCancel2(parent context.Context) context.Context {
	ctx, _ := WithCancel(parent)
	return ctx
}

// WithDeadline is similar with context.WithDeadline.
// But it cancels the context when it was done through parent or deadline.
func WithDeadline(parent context.Context, d time.Time) (context.Context, context.CancelFunc) {
	return AutoCancel(context.WithDeadline(parent, d))
}

// WithDeadline2 is similar with WithDeadline.
// It returns only a new context inherited from parent.
func WithDeadline2(parent context.Context, d time.Time) context.Context {
	ctx, _ := WithDeadline(parent, d)
	return ctx
}

// WithTimeout is similar with context.WithTimeout.
// But it cancels the context when it was done through parent or timeout.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return AutoCancel(context.WithTimeout(parent, timeout))
}

// WithTimeout2 is similar with WithTimeout.
// It returns only a new context inherited from parent.
func WithTimeout2(parent context.Context, timeout time.Duration) context.Context {
	ctx, _ := WithTimeout(parent, timeout)
	return ctx
}
