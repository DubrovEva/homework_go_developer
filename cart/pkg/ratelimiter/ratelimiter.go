package ratelimiter

import (
	"context"
	"errors"
	"time"
)

var ErrContextCancelled = errors.New("context cancelled")

type RateLimiter struct {
	tokens chan struct{}
}

func NewRateLimiter(ctx context.Context, rps int) *RateLimiter {
	tokens := make(chan struct{}, rps)
	for i := 0; i < rps; i++ {
		tokens <- struct{}{}
	}

	rl := &RateLimiter{
		tokens: tokens,
	}

	ticker := time.NewTicker(time.Second / time.Duration(rps))

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				select {
				case rl.tokens <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return rl
}

func (rl *RateLimiter) Allow(ctx context.Context) error {
	select {
	case <-rl.tokens:
		return nil
	case <-ctx.Done():
		return ErrContextCancelled
	}
}
