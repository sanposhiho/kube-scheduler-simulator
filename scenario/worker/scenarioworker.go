package worker

import (
	"context"
	"errors"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type ScenarioWorker struct {
	scenario *simulationv1alpha1.Scenario
	steppers steppers
	stopCh   chan<- struct{}
	client   client.Client
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
	currentStep := w.scenario.Status.StepStatus.Step.Major
	for {
		w.runOneStep(ctx, currentStep)
	}
}

func (w *ScenarioWorker) runOneStep(ctx context.Context, step int32) simulationv1alpha1.ScenarioPhase {
	stepper, err := w.steppers.next(step)
	if err != nil {
		if errors.Is(err, ErrNoStepper) {
			// no more step
			return simulationv1alpha1.ScenarioPaused
		}

		// TODO: set failed message
		return simulationv1alpha1.ScenarioFailed
	}

	if err := w.changeStepPhase(ctx, simulationv1alpha1.Operating); err != nil {
		// TODO: set failed message
		return simulationv1alpha1.ScenarioFailed
	}

	finish, err := stepper.run(ctx)
	if err != nil {
		// TODO: set failed message
		return simulationv1alpha1.ScenarioFailed
	}

	if err := w.changeStepPhase(ctx, simulationv1alpha1.OperatingCompleted); err != nil {
		// TODO: set failed message
		return simulationv1alpha1.ScenarioFailed
	}

	if finish {
		return simulationv1alpha1.ScenarioSucceeded
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

func (w *ScenarioWorker) changeStepPhase(ctx context.Context, phase simulationv1alpha1.StepPhase) error {
	w.scenario.Status.StepStatus.Phase = phase
	if err := w.client.Status().Update(ctx, w.scenario); err != nil {
		return fmt.Errorf("update step phase: %w", err)
	}
	return nil
}
