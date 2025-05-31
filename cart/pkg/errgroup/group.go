package errgroup

import (
	"context"
	"sync"
)

type dummy struct{}

type Group struct {
	cancel  context.CancelFunc
	wg      sync.WaitGroup
	errOnce sync.Once
	err     error

	limiter chan dummy
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	return &Group{cancel: cancel}, ctx
}

func (g *Group) SetLimit(n int) {
	if n > 0 {
		g.limiter = make(chan dummy, n)
	}
}

func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		if g.limiter != nil {
			g.limiter <- dummy{}
			defer func() { <-g.limiter }()
		}

		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}
