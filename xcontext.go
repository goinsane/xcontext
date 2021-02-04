// Package xcontext provides extended utilities for context.
package xcontext

import (
	"context"
	"sync"
	"time"
)

// DelayContext creates a new cancel context.
// It would be cancelled when ctx is done except before delay.
// Calling the cancel function is not necessary.
func DelayContext(ctx context.Context, delay time.Duration) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(context.Background())
	go func() {
		parentCtx := ctx
		delayCh := time.After(delay)
		delayOkCh := make(chan time.Time)
		for done := false; !done; {
			select {
			case <-parentCtx.Done():
				parentCtx = context.Background()
				if delayCh == delayOkCh {
					done = true
				}
			case <-delayCh:
				delayCh = delayOkCh
				if parentCtx != ctx {
					done = true
				}
			case <-newCtx.Done():
				done = true
			}
		}
		newCtxCancel()
	}()
	return newCtx, newCtxCancel
}

// DelayContext2 is similar with DelayContext except not returns cancel function.
func DelayContext2(ctx context.Context, delay time.Duration) context.Context {
	newCtx, _ := DelayContext(ctx, delay)
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
// It would be cancelled when all of the sub contexts are done.
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
func Or(ctxs ...context.Context) (context.Context, context.CancelFunc) {
	return MultiContext(context.Background(), ctxs...)
}

// Or2 is similar with Or except not returns cancel function.
func Or2(ctxs ...context.Context) context.Context {
	return MultiContext2(context.Background(), ctxs...)
}

// And creates a new cancel context to cancel when all of the contexts are done.
func And(ctxs ...context.Context) (context.Context, context.CancelFunc) {
	return WaitContext(context.Background(), ctxs...)
}

// And2 is similar with And except not returns cancel function.
func And2(ctxs ...context.Context) context.Context {
	return WaitContext2(context.Background(), ctxs...)
}

// AutoCancel cancels underlying context specified with cancel function automatically.
func AutoCancel(ctx context.Context, cancel context.CancelFunc) context.Context {
	go func() {
		<-ctx.Done()
		cancel()
	}()
	return ctx
}

// WithDeadline2 is similar with context.WithDeadline, except that it doesn't need to cancel context.
func WithDeadline2(parent context.Context, d time.Time) context.Context {
	return AutoCancel(context.WithDeadline(parent, d))
}

// WithTimeout2 is similar with context.WithTimeout, except that it doesn't need to cancel context.
func WithTimeout2(parent context.Context, timeout time.Duration) context.Context {
	return AutoCancel(context.WithTimeout(parent, timeout))
}
