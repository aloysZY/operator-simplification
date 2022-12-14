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

package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type AloysResources struct {
	/*	// +kubebuilder:default=20m*/
	Cpu resource.Quantity `json:"cpu,omitempty,omitempty"`
	/*	// +kubebuilder:default=64Mi*/
	Memory resource.Quantity `json:"memory,omitempty"`
}

type AloysContainers struct {
	Name  string `json:"name"`
	Image string `json:"image"`
	// crd 字段规范设置
	/*	+kubebuilder:validation:Maximum=65536
		+kubebuilder:validation:Minimum=1024*/
	Port    int32          `json:"port"`
	Limits  AloysResources `json:"limits,omitempty"`
	Request AloysResources `json:"request,omitempty"`
	// 暂时先让他可以为空
	/*	// +kubebuilder:validation:Pattern:=^/*/
	MountPath []string `json:"mountPath,omitempty"`
}

type AloysDeployment struct {
	//  +kubebuilder:default=1
	Replicas   int32              `json:"replicas"`
	Containers []*AloysContainers `json:"containers"`
}

type AloysIngress struct {
	// +kubebuilder:validation:Enum:={true,false}
	Enable bool `json:"enable"`
	// 添加,omitempty 可以不设置，因为Enable是 false 是没必要设置了
	Host string `json:"host,omitempty"`
	// +kubebuilder:validation:Pattern:=^/
	Path string `json:"path,omitempty"`
}

type AloysService struct {
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Enum:={true,false}
	Enable bool `json:"enable"`
}

type AloysCM struct {
	CmKey  string `json:"cmKey,omitempty"`
	CmDate string `json:"cmDate,omitempty"`
}

// AloysSpec defines the desired state of Aloys
type AloysSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Aloys. Edit aloys_types.go to remove/update
	// Foo string `json:"foo,omitempty"`
	ConfigMap  []*AloysCM      `json:"configMap,omitempty"`
	Deployment AloysDeployment `json:"deployment"`
	Service    AloysService    `json:"service"`
	Ingress    AloysIngress    `json:"ingress"`
}

// AloysStatus defines the observed state of Aloys
type AloysStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// crd 相关权限设置,containers是一个切片了，需要指定显示的container下标
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:categories="all",path="aloys",shortName="zy",singular="zy"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".spec.deployment.replicas"
// +kubebuilder:printcolumn:name="Port",type="integer",JSONPath=".spec.deployment.containers[0].port"
// +kubebuilder:printcolumn:name="Image",type="string",JSONPath=".spec.deployment.containers[0].image"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Aloys is the Schema for the aloys API
type Aloys struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AloysSpec   `json:"spec,omitempty"`
	Status AloysStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AloysList contains a list of Aloys
type AloysList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Aloys `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Aloys{}, &AloysList{})
}
