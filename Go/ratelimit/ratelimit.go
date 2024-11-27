//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	maxCount int
	interval time.Duration
	ch       chan struct{}
	stopped  chan struct{}
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	limiter := &Limiter{
		maxCount: maxCount,
		interval: interval,
		ch:       make(chan struct{}, maxCount),
		stopped:  make(chan struct{}),
	}

	if interval == 0 {
		return limiter
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-limiter.stopped:
				return
			case <-ticker.C:
				for i := 0; i < maxCount; i++ {
					select {
					case <-limiter.ch:
					default:
						break
					}
				}
			}
		}
	}()

	return limiter
}

func (l *Limiter) Acquire(ctx context.Context) error {
	for {
		select {
		case <-l.stopped:
			return ErrStopped
		default:
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				if l.interval == 0 {
					return nil
				} else {
					select {
					case l.ch <- struct{}{}:
						return nil
					default:
						continue
					}
				}
			}
		}
	}
}

func (l *Limiter) Stop() {
	close(l.stopped)
}
