package manager

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"k8s.io/apimachinery/pkg/util/wait"
)

type Manager struct {
	waiters []ControllerWaiter
	stopCh  <-chan struct{}
}

type ControllerWaiter interface {
	Name() string
	WaitConditionFunc(ctx context.Context) (wait.ConditionFunc, error)
}

func New(waiters []ControllerWaiter) *Manager {
	return &Manager{
		waiters: waiters,
		stopCh:  make(<-chan struct{}),
	}
}

var interval = 5 * time.Second

func (m *Manager) Run(ctx context.Context) error {

	return nil
}

func (m *Manager) wait(ctx context.Context) error {
	eg := utils.NewErrGroupWithSemaphore(ctx)
	for _, waiter := range m.waiters {
		waitFn, err := waiter.WaitConditionFunc(ctx)
		if err != nil {
			return fmt.Errorf("wait for waiter %s: %w", waiter.Name(), err)
		}

		eg.Grp.Go()
		if err := wait.PollUntil(interval, waitFn, m.stopCh); err != nil {
			// stopCh is called
			if errors.Is(err, wait.ErrWaitTimeout) {
				return nil
			}
			return err
		}
	}
}
