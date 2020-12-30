// Package xcontext provides extended utilities for context.
package xcontext

import (
	"context"
	"time"
)

// DelayContext creates a new cancel context. It would be cancel when ctx is done except before delay.
// Calling the cancel function is not necessary.
func DelayContext(ctx context.Context, delay time.Duration) (context.Context, context.CancelFunc) {
	newCtx, newCtxCancel := context.WithCancel(context.Background())
	go func() {
		delayCh := time.After(delay)
		for done := false; !done; {
			select {
			case <-ctx.Done():
				select {
				case <-delayCh:
					done = true
				default:
					<-time.After(250 * time.Millisecond)
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
// It would be cancelled when at least one of the sub context is done.
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
