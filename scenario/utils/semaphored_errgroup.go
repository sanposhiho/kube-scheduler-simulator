package utils

import (
	"context"
	"runtime"

	"golang.org/x/xerrors"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

type ErrGroupWithSemaphore struct {
	g *errgroup.Group
	s *semaphore.Weighted
}

func NewErrGroupWithSemaphore(ctx context.Context) *ErrGroupWithSemaphore {
	g, _ := errgroup.WithContext(ctx)
	sem := semaphore.NewWeighted(int64(runtime.GOMAXPROCS(0)))
	return &ErrGroupWithSemaphore{
		g: g,
		s: sem,
	}
}

func (e *ErrGroupWithSemaphore) Go(fn func() error) error {
	ctx := context.Background()
	if err := e.s.Acquire(ctx, 1); err != nil {
		return xerrors.Errorf("acquire semaphore: %w", err)
	}

	e.g.Go(func() error {
		defer e.s.Release(1)
		return fn()
	})
	return nil
}

func (e *ErrGroupWithSemaphore) Wait() error {
	return e.g.Wait()
}
