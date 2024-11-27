package dupcall

import (
	"context"
	"sync"
)

type Call struct {
	mutex   sync.Mutex
	running bool
	result  interface{}
	err     error
	waiters []chan struct{}
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (interface{}, error) {
	o.mutex.Lock()
	if o.running {
		waiter := make(chan struct{})
		o.waiters = append(o.waiters, waiter)
		o.mutex.Unlock()

		select {
		case <-waiter:
			o.mutex.Lock()
			defer o.mutex.Unlock()
			return o.result, o.err
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	o.running = true
	o.mutex.Unlock()

	done := make(chan struct{})
	go func() {
		defer close(done)

		result, err := cb(ctx)

		o.mutex.Lock()
		o.result, o.err = result, err
		o.running = false

		for _, waiter := range o.waiters {
			close(waiter)
		}
		o.waiters = nil
		o.mutex.Unlock()
	}()

	select {
	case <-done:
		o.mutex.Lock()
		defer o.mutex.Unlock()
		return o.result, o.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
