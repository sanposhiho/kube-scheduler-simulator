package worker

import simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"

type ScenarioWorker struct {
	scenario *simulationv1alpha1.Scenario
	steppers map[simulationv1alpha1.ScenarioStep]*stepper
	stopCh   chan<- struct{}
}

func New(scenario *simulationv1alpha1.Scenario) *ScenarioWorker {
	return &ScenarioWorker{
		scenario: scenario,
		steppers: buildSteppersMap(scenario),
		stopCh:   make(chan<- struct{}),
	}
}

func Run(stopCh chan<- struct{}) {

}

func (w *ScenarioWorker) handleUpdate(new *simulationv1alpha1.Scenario) error {
	// TODO: we need to validate the change.
	w.scenario = new
	return nil
}

func (w *ScenarioWorker) Stop() {
	w.stopCh <- struct{}{}
}

func (w *ScenarioWorker) handleDelete() {
	w.Stop()
}
