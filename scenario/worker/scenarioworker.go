package worker

import simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"

type ScenarioWorker struct {
	scenario *simulationv1alpha1.Scenario
	stopCh   chan<- struct{}
}

func New(scenario *simulationv1alpha1.Scenario) *ScenarioWorker {
	return &ScenarioWorker{
		scenario: scenario,
	}
}

func Run(stopCh chan<- struct{}) {

}

func (w *ScenarioWorker) handleUpdate(new *simulationv1alpha1.Scenario) error {
	// TODO: we need to validate the change.
	w.scenario = new
}

func (w *ScenarioWorker) handleDelete()
