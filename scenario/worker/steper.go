package worker

import simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"

type stepper struct {
	step   simulationv1alpha1.ScenarioStep
	events []*simulationv1alpha1.ScenarioEvent
}

func buildSteppersMap(scenario *simulationv1alpha1.Scenario) map[simulationv1alpha1.ScenarioStep]*stepper {
	eventmap := make(map[simulationv1alpha1.ScenarioStep][]*simulationv1alpha1.ScenarioEvent)
	for _, event := range scenario.Spec.Events {
		if _, ok := eventmap[event.Step]; !ok {
			eventmap[event.Step] = []*simulationv1alpha1.ScenarioEvent{}
		}
		eventmap[event.Step] = append(eventmap[event.Step], event)
	}

	steppers := make(map[simulationv1alpha1.ScenarioStep]*stepper)
	for step, events := range eventmap {
		steppers[step] = &stepper{
			step:   step,
			events: events,
		}
	}

	return steppers
}

func (s *stepper) run() error {
	//	for _, event := range s.events {
	//		event = event
	//	}
	return nil
}
