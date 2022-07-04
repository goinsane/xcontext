package xcontext

import (
	"context"
)

// CancelableContext is a custom implementation of context.Context with the Cancel() method to cancel context.Context.
type CancelableContext interface {
	context.Context
	Cancel()
}

type cancelableContext struct {
	context.Context
	context.CancelFunc
}

// Cancel calls the cancel function of the underlying context.
func (c *cancelableContext) Cancel() {
	c.CancelFunc()
}

// NewCancelableContext returns the underlying context as CancelableContext.
// The code should call CancelableContext.Cancel method or the cancel function to release resources associated with it.
func NewCancelableContext(ctx context.Context, cancel context.CancelFunc) (CancelableContext, context.CancelFunc) {
	result := new(cancelableContext)
	result.Context, result.CancelFunc = ctx, cancel
	return result, cancel
}

// WithCancelable creates a new cancel context as CancelableContext.
// The code should call CancelableContext.Cancel method or the cancel function to release resources associated with it.
func WithCancelable(parent context.Context) (CancelableContext, context.CancelFunc) {
	return NewCancelableContext(context.WithCancel(parent))
}

// WithCancelable2 is similar with WithCancelable.
// It returns only a new context inherited from parent.
func WithCancelable2(parent context.Context) CancelableContext {
	ctx, _ := WithCancelable(parent)
	return ctx
}

// WithCancelableAutoCancel is similar with WithCancelable.
// But it cancels the context when it was done through parent.
func WithCancelableAutoCancel(parent context.Context) (CancelableContext, context.CancelFunc) {
	return NewCancelableContext(AutoCancel(context.WithCancel(parent)))
}

// WithCancelableAutoCancel2 is similar with WithCancelableAutoCancel.
// It returns only a new context inherited from parent.
func WithCancelableAutoCancel2(parent context.Context) CancelableContext {
	ctx, _ := WithCancelableAutoCancel(parent)
	return ctx
}
