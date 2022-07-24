package worker

import (
	"sort"
	"sync"

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

func (s *steppers) add(newstepper stepper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.s = append(s.s, newstepper)
}

func (s *steppers) next() stepper {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := s.s[0]
	s.s = s.s[1:]

	return next
}
