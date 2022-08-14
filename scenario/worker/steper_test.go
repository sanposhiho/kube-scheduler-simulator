package worker

import (
	"reflect"
	"testing"

	simulationv1alpha1 "sigs.k8s.io/kube-scheduler-simulator/scenario/api/v1alpha1"
)

func Test_buildSteppersMap(t *testing.T) {
	type args struct {
		scenario *simulationv1alpha1.Scenario
	}
	tests := []struct {
		name  string
		args  args
		want  map[simulationv1alpha1.ScenarioStep]*stepper
		want1 []simulationv1alpha1.ScenarioStep
	}{
		{
			name: "happy",
			args: args{
				scenario: &simulationv1alpha1.Scenario{
					Spec: simulationv1alpha1.ScenarioSpec{
						Operations: []*simulationv1alpha1.ScenarioOperation{
							{
								ID:   "2",
								Step: 1,
							},
							{
								ID:   "3",
								Step: 1,
							},
							{
								ID:   "4",
								Step: 2,
							},
							{
								ID:   "1",
								Step: 0,
							},
						},
					},
				},
			},
			want: map[simulationv1alpha1.ScenarioStep]*stepper{
				0: {
					step: 0,
					operations: []*simulationv1alpha1.ScenarioOperation{
						{
							ID:   "1",
							Step: 0,
						},
					},
				},
				1: {
					step: 1,
					operations: []*simulationv1alpha1.ScenarioOperation{
						{
							ID:   "2",
							Step: 1,
						},
						{
							ID:   "3",
							Step: 1,
						},
					},
				},
				2: {
					step: 2,
					operations: []*simulationv1alpha1.ScenarioOperation{
						{
							ID:   "4",
							Step: 2,
						},
					},
				},
			},
			want1: []simulationv1alpha1.ScenarioStep{0, 1, 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := buildSteppersMap(tt.args.scenario)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildSteppersMap() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("buildSteppersMap() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
