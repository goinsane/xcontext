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

func (c *terminateContext) Terminate() {
	c.CancelFunc()
}

// WithTerminate creates a new TerminateContext.
func WithTerminate(parent context.Context) TerminateContext {
	ctx := new(terminateContext)
	ctx.Context, ctx.CancelFunc = context.WithCancel(parent)
	return ctx
}
