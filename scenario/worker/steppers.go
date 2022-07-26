package worker

import (
	"errors"
	"sort"
	"sync"

	"golang.org/x/xerrors"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type steppers struct {
	// steps is the list sorted by ascending order of step.
	steps []v1alpha1.ScenarioStep
	// steppermap is the map keyed by step.
	steppermap map[v1alpha1.ScenarioStep]*stepper
	mu         sync.Mutex
}

func newSteppers(scenario *v1alpha1.Scenario) steppers {
	eventmap := make(map[v1alpha1.ScenarioStep][]*v1alpha1.ScenarioOperation)
	for _, event := range scenario.Spec.Operations {
		if _, ok := eventmap[event.Step]; !ok {
			eventmap[event.Step] = []*v1alpha1.ScenarioOperation{}
		}
		eventmap[event.Step] = append(eventmap[event.Step], event)
	}

	steppermap := make(map[v1alpha1.ScenarioStep]*stepper)
	stepList := make([]v1alpha1.ScenarioStep, 0, len(eventmap))
	for step, events := range eventmap {
		steppermap[step] = &stepper{
			step:   step,
			events: events,
		}
		stepList = append(stepList, step)
	}

	sort.Slice(stepList, func(i, j int) bool {
		return stepList[i] < stepList[j]
	})

	return steppers{
		steps:      stepList,
		steppermap: steppermap,
	}
}

func (s *steppers) add(newstepper stepper) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.steppermap[newstepper.step]
	if ok {
		return xerrors.New("stepper is already exist")
	}
	s.steppermap[newstepper.step] = &newstepper

	s.steps = append(s.steps, newstepper.step)

	sort.Slice(s.steps, func(i, j int) bool {
		return s.steps[i] < s.steps[j]
	})
}

var ErrNoStepper = errors.New("steppers doesn't have stepper")

// next fetches stepper which step is after a given currentStep.
func (s *steppers) next(currentStep v1alpha1.ScenarioStep) (*stepper, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.steps) == 0 {
		return nil, ErrNoStepper
	}

	var nextstep v1alpha1.ScenarioStep
	for {
		nextstep = s.steps[0]
		s.steps = s.steps[1:]
		if currentStep < nextstep {
			break
		}
	}

	stepper, ok := s.steppermap[nextstep]
	if !ok {
		return nil, xerrors.New("stepper is not found")
	}

	delete(s.steppermap, nextstep)

	return stepper, nil
}
