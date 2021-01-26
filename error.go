package xcontext

import (
	"context"
	"errors"
)

var (
	ErrAlreadyLocked = errors.New("already locked")
)

var (
	isContextError = func(err error) bool {
		return err == context.Canceled || err == context.DeadlineExceeded
	}
)

// IsContextError returns whether err is context error.
func IsContextError(err error) bool {
	return isContextError(err)
}
