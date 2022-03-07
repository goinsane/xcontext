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
// It cancels the context when it was done through parent, deadline, timeout or in any way.
// The code should call Terminate method or cancel function to release resources associated with it.
func NewTerminateContext(ctx context.Context, cancel context.CancelFunc) TerminateContext {
	result := new(terminateContext)
	result.Context, result.CancelFunc = ctx, cancel
	AutoCancel(ctx, cancel)
	return result
}

// WithTerminate creates a new cancel context as TerminateContext.
// It cancels the context when it was done through parent.
// The code should call Terminate method to release resources associated with it, as cancel function.
func WithTerminate(parent context.Context) TerminateContext {
	return NewTerminateContext(context.WithCancel(parent))
}
