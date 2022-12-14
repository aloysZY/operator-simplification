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

package v1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Aloys) DeepCopyInto(out *Aloys) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Aloys.
func (in *Aloys) DeepCopy() *Aloys {
	if in == nil {
		return nil
	}
	out := new(Aloys)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Aloys) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysCM) DeepCopyInto(out *AloysCM) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysCM.
func (in *AloysCM) DeepCopy() *AloysCM {
	if in == nil {
		return nil
	}
	out := new(AloysCM)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysContainers) DeepCopyInto(out *AloysContainers) {
	*out = *in
	in.Limits.DeepCopyInto(&out.Limits)
	in.Request.DeepCopyInto(&out.Request)
	if in.MountPath != nil {
		in, out := &in.MountPath, &out.MountPath
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysContainers.
func (in *AloysContainers) DeepCopy() *AloysContainers {
	if in == nil {
		return nil
	}
	out := new(AloysContainers)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysDeployment) DeepCopyInto(out *AloysDeployment) {
	*out = *in
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]*AloysContainers, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(AloysContainers)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysDeployment.
func (in *AloysDeployment) DeepCopy() *AloysDeployment {
	if in == nil {
		return nil
	}
	out := new(AloysDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysIngress) DeepCopyInto(out *AloysIngress) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysIngress.
func (in *AloysIngress) DeepCopy() *AloysIngress {
	if in == nil {
		return nil
	}
	out := new(AloysIngress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysList) DeepCopyInto(out *AloysList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Aloys, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysList.
func (in *AloysList) DeepCopy() *AloysList {
	if in == nil {
		return nil
	}
	out := new(AloysList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *AloysList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysResources) DeepCopyInto(out *AloysResources) {
	*out = *in
	out.Cpu = in.Cpu.DeepCopy()
	out.Memory = in.Memory.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysResources.
func (in *AloysResources) DeepCopy() *AloysResources {
	if in == nil {
		return nil
	}
	out := new(AloysResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysService) DeepCopyInto(out *AloysService) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysService.
func (in *AloysService) DeepCopy() *AloysService {
	if in == nil {
		return nil
	}
	out := new(AloysService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysSpec) DeepCopyInto(out *AloysSpec) {
	*out = *in
	if in.ConfigMap != nil {
		in, out := &in.ConfigMap, &out.ConfigMap
		*out = make([]*AloysCM, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(AloysCM)
				**out = **in
			}
		}
	}
	in.Deployment.DeepCopyInto(&out.Deployment)
	out.Service = in.Service
	out.Ingress = in.Ingress
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysSpec.
func (in *AloysSpec) DeepCopy() *AloysSpec {
	if in == nil {
		return nil
	}
	out := new(AloysSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AloysStatus) DeepCopyInto(out *AloysStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AloysStatus.
func (in *AloysStatus) DeepCopy() *AloysStatus {
	if in == nil {
		return nil
	}
	out := new(AloysStatus)
	in.DeepCopyInto(out)
	return out
}
