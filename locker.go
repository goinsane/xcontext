package xcontext

import (
	"context"
	"sync"
)

// A Locker is same with sync.Mutex except that it implements several methods
// using context.Context as an argument; e.g., LockContext(), TryLock(), ...
//
// A Locker must not be copied after first use, same as sync.Mutex.
type Locker struct {
	mu sync.Mutex
	ch chan struct{}
}

// LockContext tries to lock l.
// If the Locker is already in use, the calling goroutine blocks until the Locker is available or ctx is done.
// If ctx is done before locking the Locker, it returns context error.
//
// If ctx is nil, it uses context.Background().
func (l *Locker) LockContext(ctx context.Context) error {
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

// TryLock tries to lock l.
// If ctx is done before calling the TryLock, it returns context error.
// If the Locker is already in use, it returns ErrAlreadyLocked immediately.
//
// If ctx is nil, it uses context.Background().
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

// Lock locks l.
// If the Locker is already in use, the calling goroutine blocks until the Locker is available.
func (l *Locker) Lock() {
	_ = l.LockContext(nil)
}

// Unlock unlocks l. It panics if l is not locked on entry to Unlock.
//
// A locked Locker is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Locker and then arrange for another goroutine to unlock it.
func (l *Locker) Unlock() {
	select {
	case <-l.ch:
	default:
		panic("unlock of unlocked Locker")
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
