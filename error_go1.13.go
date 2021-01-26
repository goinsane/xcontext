// +build go1.13

package xcontext

import (
	"context"
	"errors"
)

func init() {
	isContextError = func(err error) bool {
		return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
	}
}
