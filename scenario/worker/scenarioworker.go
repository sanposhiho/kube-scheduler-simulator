package worker

import (
	"context"
	"errors"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type ScenarioWorker struct {
	scenario *simulationv1alpha1.Scenario
	steppers steppers
	stopCh   chan<- struct{}
}

func New(scenario *simulationv1alpha1.Scenario) *ScenarioWorker {
	steppers := newSteppers(scenario)
	return &ScenarioWorker{
		scenario: scenario,
		steppers: steppers,
		stopCh:   make(chan<- struct{}),
	}
}

func (w *ScenarioWorker) Run(stopCh chan<- struct{}) {
	// TODO(sanposhiho): configure timeout from somewhere.
	ctx := context.Background()
	currentStep := w.scenario.Status.StepStatus.Step
	for {
		stepper, err := w.steppers.next(currentStep)
		if err != nil {
			if errors.Is(err, ErrNoStepper) {
				// no more step

				// TODO: change Scenario state to success.
				break
			}

			// TODO: change Scenario state to fail.
			break
		}

		finish, err := stepper.run(ctx)
		if err != nil {
			// TODO: change Scenario state to fail.
			return
		}
		if finish {
			// TODO: change Scenario state to success.
			return
		}
		currentStep = stepper.step

		// update scenario
	}
}

func (w *ScenarioWorker) handleUpdate(new *simulationv1alpha1.Scenario) error {
	// TODO: we need to validate the change.
	w.scenario = new
	w.steppers = newSteppers(new)
	return nil
}

func (w *ScenarioWorker) Stop() {
	w.stopCh <- struct{}{}
}

func (w *ScenarioWorker) handleDelete() {
	w.Stop()
}
