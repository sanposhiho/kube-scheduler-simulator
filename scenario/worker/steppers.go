package worker

import (
	"errors"
	"sort"
	"sync"

	"golang.org/x/xerrors"

	"sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

type steppers struct {
	// steps is the list sorted by ascending order of major step.
	steps []int32
	// steppermap is the map keyed by step.
	steppermap map[int32]*stepper
	mu         sync.Mutex
}

func newSteppers(scenario *v1alpha1.Scenario) steppers {
	operationmap := make(map[int32][]*v1alpha1.ScenarioOperation)
	for _, o := range scenario.Spec.Operations {
		if _, ok := operationmap[o.MajorStep]; !ok {
			operationmap[o.MajorStep] = []*v1alpha1.ScenarioOperation{}
		}
		operationmap[o.MajorStep] = append(operationmap[o.MajorStep], o)
	}

	steppermap := make(map[int32]*stepper)
	stepList := make([]int32, 0, len(operationmap))
	for step, operations := range operationmap {
		steppermap[step] = &stepper{
			step:       step,
			operations: operations,
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

// next fetches stepper which step is after a given currentStep.Major.
func (s *steppers) next(currentMajorStep int32) (*stepper, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.steps) == 0 {
		return nil, ErrNoStepper
	}

	var nextstep int32
	for {
		nextstep = s.steps[0]
		s.steps = s.steps[1:]
		if currentMajorStep < nextstep {
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
