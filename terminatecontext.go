package xcontext

import (
	"context"
)

// TerminateContext is a custom implementation of context.Context with Terminate() method to cancel context.Context.
type TerminateContext interface {
	context.Context
	Terminate()
}

type terminateContext struct {
	context.Context
	context.CancelFunc
}

// Terminate calls cancel function of the underlying context.
func (c *terminateContext) Terminate() {
	c.CancelFunc()
}

// NewTerminateContext returns the underlying context as TerminateContext.
// The code should call Terminate method or cancel function to release resources associated with it.
func NewTerminateContext(ctx context.Context, cancel context.CancelFunc) TerminateContext {
	result := new(terminateContext)
	result.Context, result.CancelFunc = ctx, cancel
	return result
}

// NewTerminateContext2 is similar with NewTerminateContext.
// But it cancels the context when it was done through parent, deadline, timeout or in any way.
func NewTerminateContext2(ctx context.Context, cancel context.CancelFunc) TerminateContext {
	AutoCancel(ctx, cancel)
	return NewTerminateContext(ctx, cancel)
}

// WithTerminate creates a new cancel context as TerminateContext.
// The code should call Terminate method to release resources associated with it, as cancel function.
func WithTerminate(parent context.Context) TerminateContext {
	return NewTerminateContext(context.WithCancel(parent))
}

// WithTerminate2 is similar with WithTerminate.
// But it cancels the context when it was done through parent.
func WithTerminate2(parent context.Context) TerminateContext {
	return NewTerminateContext(WithCancel(parent))
}
