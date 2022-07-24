/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScenarioSpec defines the desired state of Scenario
type ScenarioSpec struct {
	// Operations field has all operations for a scenario.
	// Also you can add a new operation during the scenario is running.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Operations []*ScenarioOperation `json:"operations"`
}

type ScenarioOperation struct {
	// ID for this operation. Normally, the system sets this field for you.
	ID string `json:"id"`
	// Step indicates the step at which the operation should be done.
	Step ScenarioStep `json:"step"`

	// One of the following four fields must be specified.
	// If more than one is specified or if all are empty, the operation is invalid and the scenario will fail.

	// Create is the operation to create new resource.
	//
	// +optional
	Create *CreateOperation `json:"createOperation,omitempty"`
	// Patch is the operation to patch a resource.
	//
	// +optional
	Patch *PatchOperation `json:"patchOperation,omitempty"`
	// Delete indicates the operation to delete a resource.
	//
	// +optional
	Delete *DeleteOperation `json:"deleteOperation,omitempty"`
	// Done indicates the operation to mark the scenario as Succeeded.
	// When finish the step DoneOperation belongs, this Scenario changes its status to Succeeded.
	//
	// +optional
	Done *DoneOperation `json:"doneOperation,omitempty"`
}

// OperationType describes Operation.
// Please see the following defined OperationType, all operation types not listed below are invalid.
type OperationType string

const (
	CreateOperationType         OperationType = "Create"
	PatchOperationType          OperationType = "Patch"
	DeleteOperationType         OperationType = "Delete"
	DoneOperationType           OperationType = "Done"
	PodScheduledOperationType   OperationType = "PodScheduled"
	PodUnscheduledOperationType OperationType = "PodUnscheduled"
	PodPreemptedOperationType   OperationType = "PodPreempted"
)

type CreateOperation struct {
	// Object is the Object to be created.
	Object *unstructured.Unstructured `json:"object"`

	// +optional
	CreateOptions metav1.CreateOptions `json:"createOptions,omitempty"`
}

type PatchOperation struct {
	TypeMeta   metav1.TypeMeta   `json:"typeMeta"`
	ObjectMeta metav1.ObjectMeta `json:"objectMeta"`
	// Patch is the patch for target.
	Patch string `json:"patch"`
	// PatchType
	PatchType types.PatchType

	// +optional
	PatchOptions metav1.PatchOptions `json:"patchOptions,omitempty"`
}

type DeleteOperation struct {
	TypeMeta   metav1.TypeMeta   `json:"typeMeta"`
	ObjectMeta metav1.ObjectMeta `json:"objectMeta"`

	// +optional
	DeleteOptions metav1.DeleteOptions `json:"deleteOptions,omitempty"`
}

type DoneOperation struct{}

// ScenarioStep is the step simply represented by numbers and defined by Scenario.spec.operations.
// ScenarioStep is moved to next step when it can no longer schedule any more Pods in that ScenarioStep.
// See also the description on ScenarioSubStep below.
type ScenarioStep int32

// ScenarioSubStep represents time smaller than ScenarioStep.
//
// Within one ScenarioStep, operations other than Scenario.spec.operations may happen like Pod Scheduling, Pod Preemption ...etc.
// And more, if you use Scenario against your original scheduler, any resource operations may happen during scheduling.
//
// It's only used in Scenario.status.scenarioResult.Timeline because ScenarioSubStep is only managed by the scenario controller,
// whereas ScenarioStep is defined in Scenario.spec.operations.
//
// Yeah, you may still confuse about the concept ScenarioStep and ScenarioSubStep.
//
// Let's say you define operation to create replicaset that creates 4 Pod in Scenario.spec.operations.
// 0. ScenarioStep is X and ScenarioSubStep is 0.
// 1. In Operating StepPhase, 4 Pods (Pod1, Pod2, Pod3 and Pod4) are created.
// 2. the scheduler starts scheduling.
// 3. the scheduler tries to schedule Pod1 and the Pod1 is scheduled to Node1.
// 4. Since Pod1 got scheduled, the Pod1.spec.NodeName will be updated.
//    And that increments the ScenarioSubStep.
// 5. the scheduler tries to schedule Pod2 next. And the Pod1 is preempted.
// 6. Since Pod1 is preempted, the Pod1
//
//
// ScenarioSubStep is moved to next step when an operation is going to be done against any resource.
type ScenarioSubStep int32

// ScenarioStatus defines the observed state of Scenario
type ScenarioStatus struct {
	// The phase is a simple, high-level summary of where the Scenario is in its lifecycle.
	//
	// +optional
	Phase ScenarioPhase `json:"phase,omitempty"`
	// A human-readable message indicating details about why the scenario is in this phase.
	//
	// +optional
	Message *string `json:"message,omitempty"`
	// StepStatus has the status related to step.
	//
	StepStatus ScenarioStepStatus
	// ScenarioResult has the result of the simulation.
	// Just before Step advances, this result is updated based on all occurrences at that step.
	//
	// +optional
	ScenarioResult ScenarioResult `json:"scenarioResult,omitempty"`
}

type ScenarioStepStatus struct {
	// Step indicates the current step.
	//
	// +optional
	Step ScenarioStep `json:"step,omitempty"`
	// Phase indicates the current phase in single step.
	//
	// Within a single step, the phase proceeds as follows:
	// 1. run all scenario.Spec.Operations defined for that step. (Operating)
	// 2. finish (1) (OperatingFinished)
	// 3. the scheduler starts scheduling. (Scheduling)
	// 4. the scheduler stops scheduling and changes scenario.Status.StepStatus.Phase to SchedulingFinished
	//    when it can no longer schedule any more Pods. (Scheduling -> SchedulingFinished)
	// 5. update status.scenarioResult and move to next step. (StepFinished)
	// +optional
	Phase StepPhase `json:"phase,omitempty"`
}

type StepPhase string

const (
	// Operating means controller is currently operating operation defined for the step.
	Operating StepPhase = "Operating"
	// OperatingFinished means controller have finished operating operation defined for the step.
	OperatingFinished StepPhase = "OperatingFinished"
	// Scheduling means scheduler is scheduling Pods.
	Scheduling StepPhase = "Scheduling"
	// SchedulingFinished means scheduler is trying to schedule Pods.
	// But, it can no longer schedule any more Pods.
	SchedulingFinished StepPhase = "SchedulingFinished"
	// StepFinished means controller is preparing to move to next step.
	StepFinished StepPhase = "Finished"
)

type ScenarioPhase string

const (
	// ScenarioPending phase indicates the scenario isn't started yet.
	// e.g. waiting for another scenario to finish running.
	ScenarioPending ScenarioPhase = "Pending"
	// ScenarioRunning phase indicates the scenario is running.
	ScenarioRunning ScenarioPhase = "Running"
	// ScenarioPaused phase indicates all ScenarioSpec.Operations
	// has been finished but has not been marked as done by ScenarioDone ScenarioOperations.
	ScenarioPaused ScenarioPhase = "Paused"
	// ScenarioSucceeded phase describes Scenario is fully completed
	// by ScenarioDone ScenarioOperations. User
	// can’t add any ScenarioOperations once
	// Scenario reached at the phase.
	ScenarioSucceeded ScenarioPhase = "Succeeded"
	// ScenarioFailed phase indicates something wrong happened during running scenario.
	// For example:
	// - the controller cannot create resource for some reason.
	// - users change the scheduler configuration via simulator API.
	ScenarioFailed  ScenarioPhase = "Failed"
	ScenarioUnknown ScenarioPhase = "Unknown"
)

type ScenarioResult struct {
	// SimulatorVersion represents the version of the simulator that runs this scenario.
	SimulatorVersion string `json:"simulatorVersion"`
	// Timeline is a map of operations keyed with ScenarioStep.
	// This may have many of the same operations as .spec.operations, but has additional PodScheduled and Delete operations for Pods
	// to represent a Pod is scheduled or preempted by the scheduler.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Timeline map[ScenarioStep][]ScenarioTimelineEvent `json:"timeline"`
}

type ScenarioTimelineEvent struct {
	// The ID will be the same as spec.ScenarioOperations.ID if it is from the defined operation.
	// Otherwise, it'll be newly generated.
	ID string
	// Step indicates the step at which the operation has been done.
	Step ScenarioStep `json:"step"`

	// Only one of the following fields must be non-empty.

	// Create is the result of ScenarioSpec.Operations.CreateOperation.
	// When Create is non nil, Operation should be "Create".
	Create *CreateOperationResult `json:"create"`
	// Patch is the result of ScenarioSpec.Operations.PatchOperation.
	// When Patch is non nil, Operation should be "Patch".
	Patch *PatchOperationResult `json:"patch"`
	// Delete is the result of ScenarioSpec.Operations.DeleteOperation.
	// When Delete is non nil, Operation should be "Delete".
	Delete *DeleteOperationResult `json:"delete"`
	// Done is the result of ScenarioSpec.Operations.DoneOperation.
	// When Done is non nil, Operation should be "Done".
	Done *DoneOperationResult `json:"done"`
	// PodScheduled represents the Pod is scheduled to a Node.
	// When PodScheduled is non nil, Operation should be "PodScheduled".
	PodScheduled *PodResult `json:"podScheduled"`
	// PodUnscheduled represents the scheduler tried to schedule the Pod, but cannot schedule to any Node.
	// When PodUnscheduled is non nil, Operation should be "PodUnscheduled".
	PodUnscheduled *PodResult `json:"podUnscheduled"`
	// PodPreempted represents the scheduler preempted the Pod.
	// When PodPreempted is non nil, Operation should be "PodPreempted".
	PodPreempted *PodResult `json:"podPreempted"`
}

type CreateOperationResult struct {
	// Operation is the operation that was done.
	Operation CreateOperation `json:"operation"`
	// Result is the resource after patch.
	Result unstructured.Unstructured `json:"result"`
}

type PatchOperationResult struct {
	// Operation is the operation that was done.
	Operation PatchOperation `json:"operation"`
	// Result is the resource after patch.
	Result unstructured.Unstructured `json:"result"`
}

type DeleteOperationResult struct {
	// Operation is the operation that was done.
	Operation DeleteOperation `json:"operation"`
}

type DoneOperationResult struct {
	// Operation is the operation that was done.
	Operation DoneOperation `json:"operation"`
}

// PodResult has the results related to the specific Pod.
// Depending on the status of the Pod, some fields may be empty.
type PodResult struct {
	Pod v1.Pod `json:"pod"`
	// BoundTo indicates to which Node the Pod was scheduled.
	BoundTo *string `json:"boundTo"`
	// PreemptedBy indicates which Pod the Pod was deleted for.
	// This field may be nil if this Pod has not been preempted.
	PreemptedBy *string `json:"preemptedBy"`
	// CreatedAt indicates when the Pod was created.
	CreatedAt ScenarioStep `json:"createdAt"`
	// BoundAt indicates when the Pod was scheduled.
	// This field may be nil if this Pod has not been scheduled.
	BoundAt *ScenarioStep `json:"boundAt"`
	// PreemptedAt indicates when the Pod was preempted.
	// This field may be nil if this Pod has not been preempted.
	PreemptedAt *ScenarioStep `json:"preemptedAt"`
	// ScheduleResult has the results of all scheduling for the Pod.
	//
	// If the scheduler working with a simulator isn't worked on scheduling framework,
	// this field will be empty.
	// TODO: add the link to doc when it's empty.
	//
	// +patchStrategy=replace
	// +optional
	ScheduleResult []ScenarioPodScheduleResult `json:"scheduleResult"`
}

type ScenarioPodScheduleResult struct {
	// Step indicates the step scheduling at which the scheduling is performed.
	Step *ScenarioStep `json:"step"`
	// AllCandidateNodes indicates all candidate Nodes before Filter.
	AllCandidateNodes []string `json:"allCandidateNodes"`
	// AllFilteredNodes indicates all candidate Nodes after Filter.
	AllFilteredNodes []string `json:"allFilteredNodes"`
	// PluginResults has each plugin’s result.
	PluginResults ScenarioPluginsResults `json:"pluginResults"`
}

type (
	NodeName   string
	PluginName string
)

type ScenarioPluginsResults struct {
	// Filter has each filter plugin’s result.
	Filter map[NodeName]map[PluginName]string `json:"filter"`
	// Score has each score plugin’s score.
	Score map[NodeName]map[PluginName]ScenarioPluginsScoreResult `json:"score"`
}

type ScenarioPluginsScoreResult struct {
	// RawScore has the score from Score method of Score plugins.
	RawScore int64 `json:"rawScore"`
	// NormalizedScore has the score calculated by NormalizeScore method of Score plugins.
	NormalizedScore int64 `json:"normalizedScore"`
	// FinalScore has score plugin’s final score calculated by normalizing with NormalizedScore and applied Score plugin weight.
	FinalScore int64 `json:"finalScore"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Scenario is the Schema for the scenarios API
type Scenario struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScenarioSpec   `json:"spec,omitempty"`
	Status ScenarioStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScenarioList contains a list of Scenario
type ScenarioList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Scenario `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Scenario{}, &ScenarioList{})
}
