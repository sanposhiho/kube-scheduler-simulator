package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/go-cmp/cmp"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/manager"

	"sigs.k8s.io/controller-runtime/pkg/client"
	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type ScenarioWorker struct {
	sync.RWMutex
	scenario *simulationv1alpha1.Scenario
	steppers steppers
	stopCh   chan<- struct{}
	client   client.Client

	manager manager.Manager
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
		finish, err := w.runOneStep(ctx, currentStep)
		if err != nil {
			// TODO: log
			// set failed message
		}
		if finish {
			break
		}
	}
}

func (w *ScenarioWorker) runOneStep(ctx context.Context, step int32) (bool, error) {
	stepper, err := w.steppers.next(step)
	if err != nil {
		if errors.Is(err, ErrNoStepper) {
			// no more step
			// simulationv1alpha1.ScenarioPhasePaused
			return false, nil
		}

		// simulationv1alpha1.ScenarioPhaseFailed
		return true, err
	}

	if err := w.changeStepPhase(ctx, simulationv1alpha1.StepPhaseOperating); err != nil {
		// simulationv1alpha1.ScenarioPhaseFailed
		return true, err
	}

	finish, err := stepper.run(ctx)
	if err != nil {
		// simulationv1alpha1.ScenarioPhaseFailed
		return true, err
	}

	if err := w.changeStepPhase(ctx, simulationv1alpha1.StepPhaseOperatingCompleted); err != nil {
		// simulationv1alpha1.ScenarioPhaseFailed
		return true, err
	}

	if finish {
		// simulationv1alpha1.ScenarioPhaseSucceeded
		return true, err
	}

	return false, nil
}

func (w *ScenarioWorker) HandleUpdate(new *simulationv1alpha1.Scenario) error {
	if !w.areOperationsChanged(new) {
		return nil
	}
	w.Lock()
	defer w.Unlock()
	w.scenario = new
	w.steppers = newSteppers(new)
	return nil
}

func (w *ScenarioWorker) stop() {
	w.stopCh <- struct{}{}
}

func (w *ScenarioWorker) HandleDelete() {
	w.stop()
}

func (w *ScenarioWorker) changeStepPhase(ctx context.Context, phase simulationv1alpha1.StepPhase) error {
	w.scenario.Status.StepStatus.Phase = phase
	if err := w.client.Status().Update(ctx, w.scenario); err != nil {
		return fmt.Errorf("update step phase: %w", err)
	}
	return nil
}

func (w *ScenarioWorker) areOperationsChanged(new *simulationv1alpha1.Scenario) bool {
	return cmp.Equal(new.Spec.Operations, w.scenario.Spec.Operations)
}
