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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// CreateDeploymentSpec defines the desired state of CreateDeployment
type CreateDeploymentSpec struct {
	//deployment's Namespace
	//+nullable
	Namespace string `json:"namespace"`

	//deployment's Name
	//+nullable
	Name string `json:"name"`
}

// CreateDeploymentStatus defines the observed state of CreateDeployment
type CreateDeploymentStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CreateDeployment is the Schema for the createdeployments API
type CreateDeployment struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CreateDeploymentSpec   `json:"spec,omitempty"`
	Status CreateDeploymentStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CreateDeploymentList contains a list of CreateDeployment
type CreateDeploymentList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CreateDeployment `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CreateDeployment{}, &CreateDeploymentList{})
}
