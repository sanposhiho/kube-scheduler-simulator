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
	// Also, you can add a new operation while the scenario runs.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Operations []*ScenarioOperation `json:"operations"`

	// Controllers have the configuration for controllers working with simulation.
	Controllers *Controllers `json:"controllers"`
}

type Controllers struct {
	// PreparingControllers is a list of controllers that should be run before SimulatedControllers.
	// They will run in parallel.
	//
	// It's an optional field.
	// All controllers registered in the simulator will be enabled automatically. (except controllers set in Simulate.)
	// So, you need to configure it only when you want to disable some controllers enabled by default.
	//
	// +optional
	PreparingControllers *ControllerSet `json:"preparingControllers"`
	// SimulatedControllers is a list of controllers that are the target of this simulation.
	// These are run one by one in the same order specified in Enabled field.
	//
	// It's a required field; no controllers will be enabled automatically.
	SimulatedControllers *ControllerSet `json:"simulatedControllers"`
}

type ControllerSet struct {
	// Enabled specifies controllers that should be enabled.
	// +listType=atomic
	Enabled []Controller `json:"enabled"`
	// Disabled specifies controllers that should be disabled.
	// When all controllers need to be disabled, an array containing only one "*" should be provided.
	// +listType=map
	// +listMapKey=name
	Disabled []Controller `json:"disabled"`
}

type Controller struct {
	Name string `json:"name"`
}

type ScenarioOperation struct {
	// ID for this operation. Normally, the system sets this field for you.
	ID string `json:"id"`
	// MajorStep indicates when the operation should be done.
	MajorStep int32 `json:"step"`

	// One of the following four fields must be specified.
	// If more than one is set or all are empty, the operation is invalid, and the scenario will fail.

	// Create is the operation to create a new resource.
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
// Please see the following defined OperationType; all operation types not listed below are invalid.
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
	PatchType types.PatchType `json:"patchType"`

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

// ScenarioStep is the time represented by a set of numbers, MajorStep and MinorStep,
// which are like hours and minutes in clocks in the real world.
// ScenarioStep.Major is moved to the next ScenarioStep.Major when the simulated controller can no longer do anything with the current cluster state.
// Scenario.Minor is moved to the next Scenario.Minor when any resources operations(create/edit/delete) happens.
type ScenarioStep struct {
	Major int32 `json:"major"`
	Minor int32 `json:"minor"`
}

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
	StepStatus ScenarioStepStatus `json:"stepStatus"`
	// ScenarioResult has the result of the simulation.
	// Just before Step advances, this result is updated based on all occurrences at that step.
	//
	// +optional
	ScenarioResult ScenarioResult `json:"scenarioResult,omitempty"`
}

type ScenarioStepStatus struct {
	// Step indicates the current ScenarioStep.
	//
	// +optional
	Step ScenarioStep `json:"step,omitempty"`
	// Phase indicates the current phase in a single step.
	//
	// +optional
	Phase StepPhase `json:"phase,omitempty"`
	// RunningSimulatedController indicates one of the simulated controllers that is currently running/paused/completed.
	RunningSimulatedController string `json:"runningSimulatedController"`
}

type StepPhase string

const (
	// StepPhaseOperating means controller is currently operating operation defined for the step.
	StepPhaseOperating StepPhase = "Operating"
	// StepPhaseOperatingCompleted means the preparing controllers have finished operating operation defined for the step.
	StepPhaseOperatingCompleted StepPhase = "OperatingCompleted"
	// StepPhaseControllerRunning means the simulated controller is working.
	StepPhaseControllerRunning StepPhase = "ControllerRunning"
	// StepPhaseControllerPaused means the simulated controller is paused(or will be paused).
	StepPhaseControllerPaused StepPhase = "ControllerPaused"
	// StepPhaseControllerCompleted means the current running simulated controller no longer do anything with the current cluster state.
	StepPhaseControllerCompleted StepPhase = "ControllerCompleted"
	// StepPhaseCompleted means the controller is preparing to move to the next step.
	StepPhaseCompleted StepPhase = "Finished"
)

type ScenarioPhase string

const (
	// ScenarioPhasePending phase indicates the scenario isn't started yet.
	// e.g., waiting for another scenario to finish running.
	ScenarioPhasePending ScenarioPhase = "Pending"
	// ScenarioPhaseRunning phase indicates the scenario is running.
	ScenarioPhaseRunning ScenarioPhase = "Running"
	// ScenarioPhasePaused phase indicates all ScenarioSpec.Operations
	// has been finished but not marked as done by ScenarioDone ScenarioOperations.
	ScenarioPhasePaused ScenarioPhase = "Paused"
	// ScenarioPhaseSucceeded phase describes Scenario is fully completed
	// by ScenarioDone ScenarioOperations. User
	// can’t add any ScenarioOperations once
	// Scenario reached this phase.
	ScenarioPhaseSucceeded ScenarioPhase = "Succeeded"
	// ScenarioPhaseFailed phase indicates something wrong happened while running the scenario.
	// For example:
	// - the controller cannot create a resource for some reason.
	// - users change the scheduler configuration via simulator API.
	ScenarioPhaseFailed  ScenarioPhase = "Failed"
	ScenarioPhaseUnknown ScenarioPhase = "Unknown"
)

type ScenarioResult struct {
	// SimulatorVersion represents the version of the simulator that runs this scenario.
	SimulatorVersion string `json:"simulatorVersion"`
	// Timeline is a map of operations keyed with ScenarioStep.Major(string).
	// This may have many of the same operations as .spec.operations but has additional PodScheduled and Delete operations for Pods
	// to represent a Pod is scheduled or preempted by the scheduler.
	//
	// +patchMergeKey=ID
	// +patchStrategy=merge
	Timeline map[string][]ScenarioTimelineEvent `json:"timeline"`
}

type ScenarioTimelineEvent struct {
	// The ID will be the same as spec.ScenarioOperations.ID if it is from the defined operation.
	// Otherwise, it'll be newly generated.
	ID string `json:"id"`
	// Step indicates the ScenarioStep at which the operation has been done.
	Step ScenarioStep `json:"step"`

	// Only one of the following fields must be non-empty.

	// Create is the result of ScenarioSpec.Operations.CreateOperation.
	Create *CreateOperationResult `json:"create"`
	// Patch is the result of ScenarioSpec.Operations.PatchOperation.
	Patch *PatchOperationResult `json:"patch"`
	// Delete is the result of ScenarioSpec.Operations.DeleteOperation.
	Delete *DeleteOperationResult `json:"delete"`
	// Done is the result of ScenarioSpec.Operations.DoneOperation.
	Done *DoneOperationResult `json:"done"`
	// PodScheduled represents the Pod is scheduled to a Node.
	PodScheduled *PodResult `json:"podScheduled"`
	// PodUnscheduled represents "the scheduler tried to schedule the Pod, but cannot schedule to any Node."
	PodUnscheduled *PodResult `json:"podUnscheduled"`
}

type CreateOperationResult struct {
	// Operation is the operation that was done.
	Operation CreateOperation `json:"operation"`
	// Result is the resource after the creation.
	Result unstructured.Unstructured `json:"result"`
}

type PatchOperationResult struct {
	// Operation is the operation that was done.
	Operation PatchOperation `json:"operation"`
	// Result is the resource after the patch.
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
	// PreemptedBy indicates which Pod the scheduler deleted this Pod for.
	// This field may be nil if this Pod has not been preempted.
	PreemptedBy *string `json:"preemptedBy"`
	// CreatedAt indicates when the Pod was created.
	CreatedAt ScenarioStep `json:"createdAt"`
	// BoundAt indicates when the scheduler schedule this Pod.
	// This field may be nil if this Pod has not been scheduled.
	BoundAt *ScenarioStep `json:"boundAt"`
	// PreemptedAt indicates when the scheduler preempted this Pod.
	// This field may be nil if this Pod has not been preempted.
	PreemptedAt *ScenarioStep `json:"preemptedAt"`
	// ScheduleResult has the results of all scheduling for the Pod.
	//
	// If the scheduler working with a simulator isn't created on the scheduling framework,
	// this field will be empty.
	// TODO: add the link to the doc when it's empty.
	//
	// +patchStrategy=replace
	// +optional
	ScheduleResult []ScenarioPodScheduleResult `json:"scheduleResult"`
}

type ScenarioPodScheduleResult struct {
	// Step indicates when the scheduler performs this schedule.
	Step *ScenarioStep `json:"step"`
	// AllCandidateNodes indicates all candidate Nodes before Filter.
	AllCandidateNodes []string `json:"allCandidateNodes"`
	// AllFilteredNodes indicates all candidate Nodes after Filter.
	AllFilteredNodes []string `json:"allFilteredNodes"`
	// PluginResults has each plugin's result.
	PluginResults ScenarioPluginsResults `json:"pluginResults"`
}

type (
	NodeName   string
	PluginName string
)

type ScenarioPluginsResults struct {
	// Filter has each filter plugin's result.
	Filter map[NodeName]map[PluginName]string `json:"filter"`
	// Score has each score plugin's score.
	Score map[NodeName]map[PluginName]ScenarioPluginsScoreResult `json:"score"`
}

type ScenarioPluginsScoreResult struct {
	// RawScore has the score from the Score method of Score plugins.
	RawScore int64 `json:"rawScore"`
	// NormalizedScore has the score calculated by the NormalizeScore method of Score plugins.
	NormalizedScore int64 `json:"normalizedScore"`
	// FinalScore has the score plugin's final score calculated by NormalizedScore and the score plugin weight.
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
