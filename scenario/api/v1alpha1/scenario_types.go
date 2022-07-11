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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ScenarioSpec defines the desired state of Scenario
type ScenarioSpec struct {
	// Events field has all operations for a scenario.
	// Also you can add new events during the scenario is running.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Events []*ScenarioEvent `json:"events"`
}

type ScenarioEvent struct {
	// ID for this event. Normally, the system sets this field for you.
	ID string `json:"id"`
	// Step indicates the step at which the event occurs.
	Step ScenarioStep `json:"step"`
	// Operation describes which operation this event wants to do.
	// Only "Create", "Patch", "Delete", "Done" are valid operations in ScenarioEvent.
	Operation OperationType `json:"operation"`

	// One of the following four fields must be specified.
	// If more than one is specified or if all are empty, the event is invalid and the scenario will fail.

	// CreateOperation is the operation to create new resource.
	// When use CreateOperation, Operation should be "Create".
	//
	// +optional
	CreateOperation *CreateOperation `json:"createOperation,omitempty"`
	// PatchOperation is the operation to patch a resource.
	// When use PatchOperation, Operation should be "Patch".
	//
	// +optional
	PatchOperation *PatchOperation `json:"patchOperation,omitempty"`
	// DeleteOperation indicates the operation to delete a resource.
	// When use DeleteOperation, Operation should be "Delete".
	//
	// +optional
	DeleteOperation *DeleteOperation `json:"deleteOperation,omitempty"`
	// DoneOperation indicates the operation to mark the scenario as DONE.
	// When use DoneOperation, Operation should be "Done".
	//
	// +optional
	DoneOperation *DoneOperation `json:"doneOperation,omitempty"`
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
	// Object is the Object to be create.
	Object unstructured.Unstructured `json:"object"`

	// +optional
	CreateOptions metav1.CreateOptions `json:"createOptions,omitempty"`
}

type PatchOperation struct {
	TypeMeta   metav1.TypeMeta   `json:"typeMeta"`
	ObjectMeta metav1.ObjectMeta `json:"objectMeta"`
	// Patch is the patch for target.
	Patch string `json:"patch"`

	// +optional
	PatchOptions metav1.PatchOptions `json:"patchOptions,omitempty"`
}

type DeleteOperation struct {
	TypeMeta   metav1.TypeMeta   `json:"typeMeta"`
	ObjectMeta metav1.ObjectMeta `json:"objectMeta"`

	// +optional
	DeleteOptions metav1.DeleteOptions `json:"deleteOptions,omitempty"`
}

type DoneOperation struct {
	Done bool `json:"done"`
}

// ScenarioStep is the step simply represented by numbers and used in the simulation.
// In ScenarioStep, step is moved to next step when it can no longer schedule any more Pods in that step.
// See [TODO: document here] for more information about ScenarioStep.
type ScenarioStep int32

// ScenarioStatus defines the observed state of Scenario
type ScenarioStatus struct {
	// The phase is a simple, high-level summary of where the Scenario is in its lifecycle.
	//
	// +optional
	Phase ScenarioPhase `json:"phase,omitempty"`
	// Current state of scheduler.
	//
	// +optional
	SchedulerStatus SchedulerStatus `json:"schedulerStatus,omitempty"`
	// A human readable message indicating details about why the scenario is in this phase.
	//
	// +optional
	Message *string `json:"message,omitempty"`
	// Step indicates the current step.
	//
	// +optional
	Step ScenarioStep `json:"step,omitempty"`
	// ScenarioResult has the result of the simulation.
	// Just before Step advances, this result is updated based on all occurrences at that step.
	//
	// +optional
	ScenarioResult ScenarioResult `json:"scenarioResult,omitempty"`
}

type SchedulerStatus string

const (
	// SchedulerWillRun indicates the scheduler is expected to start to schedule.
	// In other words, the scheduler is currently stopped,
	// and will start to schedule Pods when the state is SchedulerWillRun.
	SchedulerWillRun SchedulerStatus = "WillRun"
	// SchedulerRunning indicates the scheduler is scheduling Pods.
	SchedulerRunning SchedulerStatus = "Running"
	// SchedulerWillStop indicates the scheduler is expected to stop scheduling.
	// In other words, the scheduler is currently scheduling Pods,
	// and will stop scheduling when the state is SchedulerWillStop.
	SchedulerWillStop SchedulerStatus = "WillStop"
	// SchedulerStoped indicates the scheduler stops scheduling Pods.
	SchedulerStoped SchedulerStatus = "Stoped"
	// SchedulerUnknown indicates the scheduler's status is unknown.
	SchedulerUnknown ScenarioPhase = "Unknown"
)

type ScenarioPhase string

const (
	// ScenarioPending phase indicates the scenario isn't started yet.
	// e.g. waiting for another scenario to finish running.
	ScenarioPending ScenarioPhase = "Pending"
	// ScenarioRunning phase indicates the scenario is running.
	ScenarioRunning ScenarioPhase = "Running"
	// ScenarioPaused phase indicates all ScenarioSpec.Events
	// has been finished but has not been marked as done by ScenarioDone ScenarioEvent.
	ScenarioPaused ScenarioPhase = "Paused"
	// ScenarioSucceeded phase describes Scenario is fully completed
	// by ScenarioDone ScenarioEvent. User
	// can’t add any ScenarioEvent once
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
	// Timeline is a map of events keyed with ScenarioStep.
	// This may have many of the same events as .spec.events, but has additional PodScheduled and Delete events for Pods
	// to represent a Pod is scheduled or preempted by the scheduler.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Timeline map[ScenarioStep][]ScenarioTimelineEvent `json:"timeline"`
}

type ScenarioTimelineEvent struct {
	// The ID will be the same as spec.ScenarioEvent.ID if it is from the defined event.
	// Otherwise, it'll be newly generated.
	ID string
	// Step indicates the step at which the event occurs.
	Step ScenarioStep `json:"step"`
	// Operation describes which operation this event wants to do.
	// Only "Create", "Patch", "Delete", "Done", "PodScheduled", "PodUnscheduled", "PodPreempted" are valid operations in ScenarioTimelineEvent.
	Operation OperationType `json:"operation"`

	// Only one of the following fields must be non-empty.

	// Create is the result of ScenarioSpec.Events.CreateOperation.
	// When Create is non nil, Operation should be "Create".
	Create *CreateOperationResult `json:"create"`
	// Patch is the result of ScenarioSpec.Events.PatchOperation.
	// When Patch is non nil, Operation should be "Patch".
	Patch *PatchOperationResult `json:"patch"`
	// Delete is the result of ScenarioSpec.Events.DeleteOperation.
	// When Delete is non nil, Operation should be "Delete".
	Delete *DeleteOperationResult `json:"delete"`
	// Done is the result of ScenarioSpec.Events.DoneOperation.
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
	// +patchStrategy=replace
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
