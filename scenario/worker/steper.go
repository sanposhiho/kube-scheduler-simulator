package worker

import (
	"context"
	"sort"

	"golang.org/x/xerrors"
	"sigs.k8s.io/kube-scheduler-simulator/scenario/utils"

	"k8s.io/client-go/rest"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

// stepper will run all the events defined in the single step.
type stepper struct {
	config *rest.Config
	step   simulationv1alpha1.ScenarioStep
	events []*simulationv1alpha1.ScenarioOperation
}

// buildSteppersMap returns two data about steps;
// - one is steppers map keyed by step.
// - another is the list sorted by ascending order of step.
func buildSteppersMap(scenario *simulationv1alpha1.Scenario) (map[simulationv1alpha1.ScenarioStep]*stepper, []simulationv1alpha1.ScenarioStep) {
	eventmap := make(map[simulationv1alpha1.ScenarioStep][]*simulationv1alpha1.ScenarioOperation)
	for _, event := range scenario.Spec.Operations {
		if _, ok := eventmap[event.Step]; !ok {
			eventmap[event.Step] = []*simulationv1alpha1.ScenarioOperation{}
		}
		eventmap[event.Step] = append(eventmap[event.Step], event)
	}

	steppers := make(map[simulationv1alpha1.ScenarioStep]*stepper)
	stepList := make([]simulationv1alpha1.ScenarioStep, 0, len(eventmap))
	for step, events := range eventmap {
		steppers[step] = &stepper{
			step:   step,
			events: events,
		}
		stepList = append(stepList, step)
	}

	sort.Slice(stepList, func(i, j int) bool {
		return stepList[i] < stepList[j]
	})

	return steppers, stepList
}

// run runs all events registered in s.step.
// It returns boolean shows whether the Scenario should finish in this step.
func (s *stepper) run(ctx context.Context) (bool, error) {
	finish, err := s.operate(ctx)
	if err != nil {
		return true, xerrors.Errorf("run .spec.operation: %w", err)
	}

	return finish, nil
}

func (s *stepper) operate(ctx context.Context) (bool, error) {
	eg := utils.NewErrGroupWithSemaphore(ctx)

	// whether the Scenario should finish in this step.
	finishFlag := false
	for _, event := range s.events {
		if err := eg.Sem.Acquire(ctx, 1); err != nil {
			return false, xerrors.Errorf("acquire semaphore: %w", err)
		}
		event := event
		eg.Grp.Go(func() error {
			finish, err := event.Run(ctx, s.config)
			if err != nil {
				return xerrors.Errorf("run event: event id: %v, step %v, error: %w", event.ID, event.Step, err)
			}
			if finish {
				// update flag
				finishFlag = finish
			}
			return nil
		})
	}

	if err := eg.Grp.Wait(); err != nil {
		return true, err
	}

	return finishFlag, nil
}
