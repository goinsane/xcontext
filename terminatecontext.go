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
// The code should call TerminateContext.Terminate method or cancel function to release resources associated with it.
func NewTerminateContext(ctx context.Context, cancel context.CancelFunc) (TerminateContext, context.CancelFunc) {
	result := new(terminateContext)
	result.Context, result.CancelFunc = ctx, cancel
	return result, cancel
}

// WithTerminate creates a new cancel context as TerminateContext.
// The code should call TerminateContext.Terminate method or cancel function to release resources associated with it.
func WithTerminate(parent context.Context) (TerminateContext, context.CancelFunc) {
	return NewTerminateContext(context.WithCancel(parent))
}

// WithTerminate2 is similar with WithTerminate.
// It returns only a new context inherited from parent.
func WithTerminate2(parent context.Context) TerminateContext {
	ctx, _ := WithTerminate(parent)
	return ctx
}

// WithTerminateAutoCancel is similar with WithTerminate.
// But it cancels the context when it was done through parent.
func WithTerminateAutoCancel(parent context.Context) (TerminateContext, context.CancelFunc) {
	return NewTerminateContext(AutoCancel(context.WithCancel(parent)))
}

// WithTerminateAutoCancel2 is similar with WithTerminateAutoCancel.
// It returns only a new context inherited from parent.
func WithTerminateAutoCancel2(parent context.Context) TerminateContext {
	ctx, _ := WithTerminateAutoCancel(parent)
	return ctx
}
