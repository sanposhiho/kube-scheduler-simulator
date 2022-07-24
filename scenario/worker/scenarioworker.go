package worker

import (
	"context"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type ScenarioWorker struct {
	scenario *simulationv1alpha1.Scenario
	// steps is the list sorted by ascending order of step.
	steps []simulationv1alpha1.ScenarioStep
	// steppers is the map keyed by step.
	steppers map[simulationv1alpha1.ScenarioStep]*stepper
	stopCh   chan<- struct{}
}

func New(scenario *simulationv1alpha1.Scenario) *ScenarioWorker {
	steppers, steps := buildSteppersMap(scenario)
	return &ScenarioWorker{
		scenario: scenario,
		steppers: steppers,
		steps:    steps,
		stopCh:   make(chan<- struct{}),
	}
}

func (w *ScenarioWorker) Run(stopCh chan<- struct{}) {
	// TODO(sanposhiho): configure timeout from somewhere.
	ctx := context.Background()
	for _, s := range w.steppers {
		finish, err := s.run(ctx)
		if err != nil {
			// TODO: change Scenario state to fail.
			return
		}
		if finish {
			// TODO: change Scenario state to success.
			return
		}
	}
}

func (w *ScenarioWorker) handleUpdate(new *simulationv1alpha1.Scenario) error {
	// TODO: we need to validate the change.
	w.scenario = new
	return nil
}

func (w *ScenarioWorker) Stop() {
	w.stopCh <- struct{}{}
}

func (w *ScenarioWorker) handleDelete() {
	w.Stop()
}
