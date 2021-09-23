/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PipelineModuleSpec defines the desired state of PipelineModule
type PipelineModuleSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Image   string   `json:"image"`
	Port    int32    `json:"port"`
	Command []string `json:"command,omitempty"`
}

// PipelineModuleStatus defines the observed state of PipelineModule
type PipelineModuleStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName="pm"
//+kubebuilder:subresource:status

// PipelineModule is the Schema for the pipelinemodules API
type PipelineModule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineModuleSpec   `json:"spec,omitempty"`
	Status PipelineModuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineModuleList contains a list of PipelineModule
type PipelineModuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PipelineModule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PipelineModule{}, &PipelineModuleList{})
}
