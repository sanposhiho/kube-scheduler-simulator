package waitermanager

import (
	"context"
	"errors"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	"golang.org/x/xerrors"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"k8s.io/apimachinery/pkg/util/wait"
)

type Manager struct {
	waiters []ControllerWaiter
	stopCh  chan struct{}
	mu      sync.Mutex
}

type ControllerWaiter interface {
	Name() string
	WaitConditionFunc(ctx context.Context) (wait.ConditionFunc, error)
}

func New(waiters ...ControllerWaiter) *Manager {
	return &Manager{
		waiters: waiters,
		stopCh:  make(chan struct{}),
	}
}

var interval = 5 * time.Second

func (m *Manager) Run(ctx context.Context, controllers sets.String) (bool, error) {
	m.mu.Lock()
	// stop previous wait goroutines.
	m.stopCh <- struct{}{}
	m.stopCh = make(chan struct{})
	m.mu.Unlock()

	if err := m.wait(ctx, controllers); err != nil {
		if errors.Is(err, wait.ErrWaitTimeout) {
			// wait.ErrWaitTimeout is returned when stopCh is called.
			return false, nil
		}
		return false, xerrors.Errorf("wait waiters: %w", err)
	}
	return true, nil
}

func (m *Manager) wait(ctx context.Context, controllers sets.String) error {
	eg := utils.NewErrGroupWithSemaphore(ctx)
	for _, waiter := range m.waiters {
		if !controllers.Has(waiter.Name()) {
			continue
		}
		waitFn, err := waiter.WaitConditionFunc(ctx)
		if err != nil {
			return xerrors.Errorf("wait for waiter %s: %w", waiter.Name(), err)
		}

		if err := eg.Go(func() error { return wait.PollUntil(interval, waitFn, m.stopCh) }); err != nil {
			return xerrors.Errorf("start an error group: %w", err)
		}
	}

	return nil
}
