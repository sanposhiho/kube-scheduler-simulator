package utils

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/errors"
)

func FetchRunningScenario(ctx context.Context, c client.Client) (*v1alpha1.Scenario, error) {
	list := v1alpha1.ScenarioList{}
	if err := c.List(ctx, &list, &client.ListOptions{
		FieldSelector: fields.OneTermEqualSelector("status.phase", string(v1alpha1.ScenarioPhaseRunning)),
	}); err != nil {
		return nil, err
	}

	if len(list.Items) > 1 {
		return nil, errors.ErrTooManyRunningScenario
	}
	if len(list.Items) == 0 {
		return nil, errors.ErrNoRunningScenario
	}

	return &list.Items[0], nil
}
