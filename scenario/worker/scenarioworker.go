package worker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/waitermanager"

	"sigs.k8s.io/controller-runtime/pkg/client"
	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type ScenarioWorker struct {
	sync.RWMutex

	ScenarioName string

	scenario *simulationv1alpha1.Scenario
	steppers steppers
	stopCh   chan struct{}
	client   client.Client

	manager *waitermanager.Manager
}

func New(scenario *simulationv1alpha1.Scenario, cli client.Client, manager *waitermanager.Manager) *ScenarioWorker {
	// TODO: remove all resource.
	steppers := newSteppers(scenario)
	return &ScenarioWorker{
		ScenarioName: scenario.Name,
		scenario:     scenario,
		steppers:     steppers,
		client:       cli,
		stopCh:       make(chan struct{}),
		manager:      manager,
	}
}

func (w *ScenarioWorker) Run() {
	// TODO(sanposhiho): configure timeout from somewhere.
	ctx := context.Background()
	currentStep := w.scenario.Status.StepStatus.Step.Major
	for {
		finish, err := w.operateOneStep(ctx, currentStep)
		if err != nil {
			if errors.Is(err, ErrNoStepper) {
				// no more step.
				w.scenario.Status.Phase = simulationv1alpha1.ScenarioPhasePaused
			} else {
				// TODO: error log
				w.scenario.Status.Phase = simulationv1alpha1.ScenarioPhaseFailed

				msg := fmt.Sprintf("failed to handle the scenario correctly: operating: %v", err)
				w.scenario.Status.Message = &msg
			}

			if err := w.updateScenarioStatus(ctx); err != nil {
				// TODO: error log
				return
			}

			return
		}

		if err := w.waitSimulatedController(ctx); err != nil {
			w.scenario.Status.Phase = simulationv1alpha1.ScenarioPhaseFailed

			msg := fmt.Sprintf("failed to handle the scenario correctly: waiting simulated controller: %v", err)
			w.scenario.Status.Message = &msg
			// TODO: error log
			return
		}

		if finish {
			w.scenario.Status.Phase = simulationv1alpha1.ScenarioPhaseSucceeded
			if err := w.updateScenarioStatus(ctx); err != nil {
				// TODO: error log
				return
			}
			return
		}
	}
}

var interval = 5 * time.Second

func (w *ScenarioWorker) operateOneStep(ctx context.Context, step int32) (bool, error) {
	w.RLock()
	defer w.RUnlock()
	stepper, err := w.steppers.next(step)
	if err != nil {
		return true, err
	}

	if err := w.changeStepPhase(ctx, simulationv1alpha1.StepPhaseOperating); err != nil {
		return true, err
	}

	finish, err := stepper.run(ctx)
	if err != nil {
		return true, err
	}

	return finish, nil
}

func (w *ScenarioWorker) waitSimulatedController(ctx context.Context) error {
	for {
		if err := wait.PollUntil(interval, w.waitOperatingCompleted(ctx), w.stopCh); err != nil {
			return err
		}

		controllers := sets.NewString()
		controllers.Insert(w.scenario.Status.StepStatus.RunningSimulatedController)
		done, err := w.manager.Run(ctx, controllers)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		// webhook canceled manager run ==> simulated controller perform some operations.
	}
}

func (w *ScenarioWorker) waitOperatingCompleted(ctx context.Context) wait.ConditionFunc {
	return func() (bool, error) {
		if err := w.syncScenario(ctx); err != nil {
			return false, err
		}

		if w.scenario.Status.StepStatus.Phase == simulationv1alpha1.StepPhaseOperatingCompleted {
			return true, nil
		}

		return false, nil
	}
}

func (w *ScenarioWorker) HandleUpdate(new *simulationv1alpha1.Scenario) {
	w.Lock()
	defer w.Unlock()

	if !w.areOperationsChanged(new) && w.scenario.Status.Phase == simulationv1alpha1.ScenarioPhaseRunning {
		// updates running scenario only when spec.operations get changed.
		return
	}

	w.scenario = new
	w.steppers = newSteppers(new)
}

func (w *ScenarioWorker) stop() {
	w.stopCh <- struct{}{}
}

func (w *ScenarioWorker) HandleDelete() {
	w.stop()
}

func (w *ScenarioWorker) changeStepPhase(ctx context.Context, phase simulationv1alpha1.StepPhase) error {
	w.scenario.Status.StepStatus.Phase = phase
	if err := w.updateScenarioStatus(ctx); err != nil {
		return fmt.Errorf("update step phase: %w", err)
	}
	return nil
}

func (w *ScenarioWorker) syncScenario(ctx context.Context) error {
	new := &simulationv1alpha1.Scenario{}
	if err := w.client.Get(ctx, client.ObjectKeyFromObject(w.scenario), new); err != nil {
		return fmt.Errorf("update running scenario: %w", err)
	}
	w.Lock()
	defer w.Unlock()
	w.scenario = new
	return nil
}

func (w *ScenarioWorker) updateScenarioStatus(ctx context.Context) error {
	if err := w.client.Status().Update(ctx, w.scenario); err != nil {
		return fmt.Errorf("update running scenario: %w", err)
	}
	return nil
}

func (w *ScenarioWorker) areOperationsChanged(new *simulationv1alpha1.Scenario) bool {
	return cmp.Equal(new.Spec.Operations, w.scenario.Spec.Operations)
}
