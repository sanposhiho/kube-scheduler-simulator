package worker

import (
	"sort"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

// stepper will run all the events defined in the single step.
type stepper struct {
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

func (s *stepper) run() error {
	for _, event := range s.events {
	}
	return nil
}
