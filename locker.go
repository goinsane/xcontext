package xcontext

import (
	"context"
	"sync"
)

type Locker struct {
	mu sync.Mutex
	ch chan struct{}
}

func (l *Locker) Lock(ctx context.Context) error {
	l.init()
	if ctx == nil {
		ctx = context.Background()
	}
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case l.ch <- struct{}{}:
			return nil
		}
	}
}

func (l *Locker) Unlock() {
	select {
	case <-l.ch:
	default:
		panic("unlock of unlocked Locker")
	}
}

func (l *Locker) TryLock(ctx context.Context) error {
	l.init()
	if ctx == nil {
		ctx = context.Background()
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case l.ch <- struct{}{}:
		return nil
	default:
		return ErrAlreadyLocked
	}
}

func (l *Locker) init() {
	if l.ch != nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.ch != nil {
		return
	}
	l.ch = make(chan struct{}, 1)
}
