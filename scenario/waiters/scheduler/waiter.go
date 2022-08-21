package scheduler

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"golang.org/x/xerrors"

	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/apimachinery/pkg/fields"
)

type schedulerWaiter struct {
	client client.Client
}

func New(client client.Client) *schedulerWaiter {
	return &schedulerWaiter{client: client}
}

func (s *schedulerWaiter) Name() string {
	return "scheduler"
}

func (s *schedulerWaiter) WaitConditionFunc(ctx context.Context) (wait.ConditionFunc, error) {
	unscheduledPods, err := s.listUnscheduledPods(ctx)
	if err != nil {
		return nil, xerrors.Errorf("fetch unscheduled Pods: %w", err)
	}

	lastCheckTime := time.Now()
	waitFn := func() (bool, error) {
		if len(unscheduledPods.Items) == 0 {
			return true, nil
		}

		currentUnscheduledPods, err := s.listUnscheduledPods(ctx)
		if err != nil {
			return false, err
		}
		if len(currentUnscheduledPods.Items) != len(unscheduledPods.Items) {
			unscheduledPods = currentUnscheduledPods
			lastCheckTime = time.Now()
			return false, nil
		}

		for _, p := range unscheduledPods.Items {
			eventOpt := &client.ListOptions{
				FieldSelector: fields.AndSelectors(
					fields.OneTermEqualSelector("reason", "FailedScheduling"),
					fields.OneTermEqualSelector("involvedObject.name", p.Name),
					fields.OneTermEqualSelector("involvedObject.kind", p.Kind),
				),
			}

			failedSchedulingEvents := &v1.EventList{}
			err := s.client.List(ctx, failedSchedulingEvents, eventOpt)
			if err != nil {
				return false, xerrors.Errorf("list events: %w", err)
			}

			isUnscheduled := false
			for _, event := range failedSchedulingEvents.Items {
				if event.LastTimestamp.After(lastCheckTime) {
					isUnscheduled = true
					break
				}
			}
			if !isUnscheduled {
				return false, nil
			}
		}
		return true, nil
	}
	return waitFn, nil
}

func (s *schedulerWaiter) listUnscheduledPods(ctx context.Context) (*v1.PodList, error) {
	pods := &v1.PodList{}
	err := s.client.List(ctx, pods, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.nodename", ""),
	})
	if err != nil {
		return nil, err
	}

	return pods, nil
}
