package manager

import (
	"context"
	"errors"
	"sync"
	"time"

	"golang.org/x/xerrors"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"k8s.io/apimachinery/pkg/util/wait"
)

type Manager struct {
	waiters []ControllerWaiter
	stopCh  chan struct{}
	sync.Mutex
}

type ControllerWaiter interface {
	Name() string
	WaitConditionFunc(ctx context.Context) (wait.ConditionFunc, error)
}

func New(waiters []ControllerWaiter) *Manager {
	return &Manager{
		waiters: waiters,
		stopCh:  make(chan struct{}),
	}
}

var interval = 5 * time.Second

func (m *Manager) Run(ctx context.Context) error {
	m.Lock()
	// stop previous wait goroutines.
	m.stopCh <- struct{}{}
	m.stopCh = make(chan struct{})
	m.Unlock()

	if err := m.wait(ctx); err != nil {
		return xerrors.Errorf("wait waiters: %w", err)
	}
	return nil
}

func (m *Manager) wait(ctx context.Context) error {
	eg := utils.NewErrGroupWithSemaphore(ctx)
	for _, waiter := range m.waiters {
		waitFn, err := waiter.WaitConditionFunc(ctx)
		if err != nil {
			return xerrors.Errorf("wait for waiter %s: %w", waiter.Name(), err)
		}

		if err := eg.Go(func() error {
			// wait.ErrWaitTimeout is returned when stopCh is called.
			if err := wait.PollUntil(interval, waitFn, m.stopCh); err != nil && !errors.Is(err, wait.ErrWaitTimeout) {
				return err
			}

			return nil
		}); err != nil {
			return xerrors.Errorf("start an error group: %w", err)
		}
	}

	return nil
}
