//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CreateOperation) DeepCopyInto(out *CreateOperation) {
	*out = *in
	if in.Object != nil {
		in, out := &in.Object, &out.Object
		*out = (*in).DeepCopy()
	}
	in.CreateOptions.DeepCopyInto(&out.CreateOptions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreateOperation.
func (in *CreateOperation) DeepCopy() *CreateOperation {
	if in == nil {
		return nil
	}
	out := new(CreateOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CreateOperationResult) DeepCopyInto(out *CreateOperationResult) {
	*out = *in
	in.Operation.DeepCopyInto(&out.Operation)
	in.Result.DeepCopyInto(&out.Result)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CreateOperationResult.
func (in *CreateOperationResult) DeepCopy() *CreateOperationResult {
	if in == nil {
		return nil
	}
	out := new(CreateOperationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeleteOperation) DeepCopyInto(out *DeleteOperation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.DeleteOptions != nil {
		in, out := &in.DeleteOptions, &out.DeleteOptions
		*out = new(v1.DeleteOptions)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteOperation.
func (in *DeleteOperation) DeepCopy() *DeleteOperation {
	if in == nil {
		return nil
	}
	out := new(DeleteOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeleteOperationResult) DeepCopyInto(out *DeleteOperationResult) {
	*out = *in
	in.Operation.DeepCopyInto(&out.Operation)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeleteOperationResult.
func (in *DeleteOperationResult) DeepCopy() *DeleteOperationResult {
	if in == nil {
		return nil
	}
	out := new(DeleteOperationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoneOperation) DeepCopyInto(out *DoneOperation) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoneOperation.
func (in *DoneOperation) DeepCopy() *DoneOperation {
	if in == nil {
		return nil
	}
	out := new(DoneOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DoneOperationResult) DeepCopyInto(out *DoneOperationResult) {
	*out = *in
	out.Operation = in.Operation
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DoneOperationResult.
func (in *DoneOperationResult) DeepCopy() *DoneOperationResult {
	if in == nil {
		return nil
	}
	out := new(DoneOperationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatchOperation) DeepCopyInto(out *PatchOperation) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	if in.PatchOptions != nil {
		in, out := &in.PatchOptions, &out.PatchOptions
		*out = new(v1.PatchOptions)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatchOperation.
func (in *PatchOperation) DeepCopy() *PatchOperation {
	if in == nil {
		return nil
	}
	out := new(PatchOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatchOperationResult) DeepCopyInto(out *PatchOperationResult) {
	*out = *in
	in.Operation.DeepCopyInto(&out.Operation)
	in.Result.DeepCopyInto(&out.Result)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatchOperationResult.
func (in *PatchOperationResult) DeepCopy() *PatchOperationResult {
	if in == nil {
		return nil
	}
	out := new(PatchOperationResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodResult) DeepCopyInto(out *PodResult) {
	*out = *in
	in.Pod.DeepCopyInto(&out.Pod)
	if in.BoundTo != nil {
		in, out := &in.BoundTo, &out.BoundTo
		*out = new(string)
		**out = **in
	}
	if in.PreemptedBy != nil {
		in, out := &in.PreemptedBy, &out.PreemptedBy
		*out = new(string)
		**out = **in
	}
	if in.BoundAt != nil {
		in, out := &in.BoundAt, &out.BoundAt
		*out = new(ScenarioStep)
		**out = **in
	}
	if in.PreemptedAt != nil {
		in, out := &in.PreemptedAt, &out.PreemptedAt
		*out = new(ScenarioStep)
		**out = **in
	}
	if in.ScheduleResult != nil {
		in, out := &in.ScheduleResult, &out.ScheduleResult
		*out = make([]ScenarioPodScheduleResult, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodResult.
func (in *PodResult) DeepCopy() *PodResult {
	if in == nil {
		return nil
	}
	out := new(PodResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Scenario) DeepCopyInto(out *Scenario) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Scenario.
func (in *Scenario) DeepCopy() *Scenario {
	if in == nil {
		return nil
	}
	out := new(Scenario)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Scenario) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioList) DeepCopyInto(out *ScenarioList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Scenario, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioList.
func (in *ScenarioList) DeepCopy() *ScenarioList {
	if in == nil {
		return nil
	}
	out := new(ScenarioList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ScenarioList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioOperation) DeepCopyInto(out *ScenarioOperation) {
	*out = *in
	if in.Create != nil {
		in, out := &in.Create, &out.Create
		*out = new(CreateOperation)
		(*in).DeepCopyInto(*out)
	}
	if in.Patch != nil {
		in, out := &in.Patch, &out.Patch
		*out = new(PatchOperation)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(DeleteOperation)
		(*in).DeepCopyInto(*out)
	}
	if in.Done != nil {
		in, out := &in.Done, &out.Done
		*out = new(DoneOperation)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioOperation.
func (in *ScenarioOperation) DeepCopy() *ScenarioOperation {
	if in == nil {
		return nil
	}
	out := new(ScenarioOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioPluginsResults) DeepCopyInto(out *ScenarioPluginsResults) {
	*out = *in
	if in.Filter != nil {
		in, out := &in.Filter, &out.Filter
		*out = make(map[NodeName]map[PluginName]string, len(*in))
		for key, val := range *in {
			var outVal map[PluginName]string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(map[PluginName]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
	if in.Score != nil {
		in, out := &in.Score, &out.Score
		*out = make(map[NodeName]map[PluginName]ScenarioPluginsScoreResult, len(*in))
		for key, val := range *in {
			var outVal map[PluginName]ScenarioPluginsScoreResult
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make(map[PluginName]ScenarioPluginsScoreResult, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioPluginsResults.
func (in *ScenarioPluginsResults) DeepCopy() *ScenarioPluginsResults {
	if in == nil {
		return nil
	}
	out := new(ScenarioPluginsResults)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioPluginsScoreResult) DeepCopyInto(out *ScenarioPluginsScoreResult) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioPluginsScoreResult.
func (in *ScenarioPluginsScoreResult) DeepCopy() *ScenarioPluginsScoreResult {
	if in == nil {
		return nil
	}
	out := new(ScenarioPluginsScoreResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioPodScheduleResult) DeepCopyInto(out *ScenarioPodScheduleResult) {
	*out = *in
	if in.Step != nil {
		in, out := &in.Step, &out.Step
		*out = new(ScenarioStep)
		**out = **in
	}
	if in.AllCandidateNodes != nil {
		in, out := &in.AllCandidateNodes, &out.AllCandidateNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.AllFilteredNodes != nil {
		in, out := &in.AllFilteredNodes, &out.AllFilteredNodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.PluginResults.DeepCopyInto(&out.PluginResults)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioPodScheduleResult.
func (in *ScenarioPodScheduleResult) DeepCopy() *ScenarioPodScheduleResult {
	if in == nil {
		return nil
	}
	out := new(ScenarioPodScheduleResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioResult) DeepCopyInto(out *ScenarioResult) {
	*out = *in
	if in.Timeline != nil {
		in, out := &in.Timeline, &out.Timeline
		*out = make(map[ScenarioStep][]ScenarioTimelineEvent, len(*in))
		for key, val := range *in {
			var outVal []ScenarioTimelineEvent
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]ScenarioTimelineEvent, len(*in))
				for i := range *in {
					(*in)[i].DeepCopyInto(&(*out)[i])
				}
			}
			(*out)[key] = outVal
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioResult.
func (in *ScenarioResult) DeepCopy() *ScenarioResult {
	if in == nil {
		return nil
	}
	out := new(ScenarioResult)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioSpec) DeepCopyInto(out *ScenarioSpec) {
	*out = *in
	if in.Operations != nil {
		in, out := &in.Operations, &out.Operations
		*out = make([]*ScenarioOperation, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ScenarioOperation)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioSpec.
func (in *ScenarioSpec) DeepCopy() *ScenarioSpec {
	if in == nil {
		return nil
	}
	out := new(ScenarioSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioStatus) DeepCopyInto(out *ScenarioStatus) {
	*out = *in
	if in.Message != nil {
		in, out := &in.Message, &out.Message
		*out = new(string)
		**out = **in
	}
	out.StepStatus = in.StepStatus
	in.ScenarioResult.DeepCopyInto(&out.ScenarioResult)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioStatus.
func (in *ScenarioStatus) DeepCopy() *ScenarioStatus {
	if in == nil {
		return nil
	}
	out := new(ScenarioStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioStepStatus) DeepCopyInto(out *ScenarioStepStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioStepStatus.
func (in *ScenarioStepStatus) DeepCopy() *ScenarioStepStatus {
	if in == nil {
		return nil
	}
	out := new(ScenarioStepStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ScenarioTimelineEvent) DeepCopyInto(out *ScenarioTimelineEvent) {
	*out = *in
	if in.Create != nil {
		in, out := &in.Create, &out.Create
		*out = new(CreateOperationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.Patch != nil {
		in, out := &in.Patch, &out.Patch
		*out = new(PatchOperationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(DeleteOperationResult)
		(*in).DeepCopyInto(*out)
	}
	if in.Done != nil {
		in, out := &in.Done, &out.Done
		*out = new(DoneOperationResult)
		**out = **in
	}
	if in.PodScheduled != nil {
		in, out := &in.PodScheduled, &out.PodScheduled
		*out = new(PodResult)
		(*in).DeepCopyInto(*out)
	}
	if in.PodUnscheduled != nil {
		in, out := &in.PodUnscheduled, &out.PodUnscheduled
		*out = new(PodResult)
		(*in).DeepCopyInto(*out)
	}
	if in.PodPreempted != nil {
		in, out := &in.PodPreempted, &out.PodPreempted
		*out = new(PodResult)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ScenarioTimelineEvent.
func (in *ScenarioTimelineEvent) DeepCopy() *ScenarioTimelineEvent {
	if in == nil {
		return nil
	}
	out := new(ScenarioTimelineEvent)
	in.DeepCopyInto(out)
	return out
}
