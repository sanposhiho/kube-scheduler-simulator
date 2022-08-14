package scheduler

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/apimachinery/pkg/fields"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clientset "k8s.io/client-go/kubernetes"
)

type schedulerWaiter struct {
	client clientset.Interface
}

func New(client clientset.Interface) *schedulerWaiter {
	return &schedulerWaiter{client: client}
}

func (s *schedulerWaiter) Name() string {

}

// これ以上何もできないかどうかだけを見る
func (s *schedulerWaiter) WaitConditionFunc(ctx context.Context) (wait.ConditionFunc, error) {
	unscheduledPods, err := s.client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("spec.nodename", "").String(),
	})
	if err != nil {
		return nil, err
	}

	lastCheckTime := time.Now()
	waitFn := func() (bool, error) {
		if len(unscheduledPods.Items) == 0 {
			// これ以上何もできない。
			return true, nil
		}

		currentUnscheduledPods, err := s.client.CoreV1().Pods(metav1.NamespaceAll).List(ctx, metav1.ListOptions{
			FieldSelector: fields.OneTermEqualSelector("spec.nodename", "").String(),
		})
		if err != nil {
			return false, err
		}
		if len(currentUnscheduledPods.Items) != len(unscheduledPods.Items) {
			unscheduledPods = currentUnscheduledPods
			lastCheckTime = time.Now()
			return false, nil
		}

		for _, p := range unscheduledPods.Items {
			eventOpt := metav1.ListOptions{
				FieldSelector: fields.AndSelectors(
					fields.OneTermEqualSelector("reason", "FailedScheduling"),
					fields.OneTermEqualSelector("involvedObject.name", p.Name),
					fields.OneTermEqualSelector("involvedObject.kind", p.Kind),
				).String(),
			}

			failedSchedulingEvents, err := s.client.CoreV1().Events(p.Namespace).List(context.Background(), eventOpt)
			if err != nil {
				return false, fmt.Errorf("list events: %w", err)
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
		// 全部がunscheduledだった
		return true, nil
	}
	return waitFn, nil
}
