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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type label struct {
}

// ChaosmonkeySpec defines the desired state of Chaosmonkey
type ChaosmonkeySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Chaosmonkey. Edit chaosmonkey_types.go to remove/update
	PodName string `json:"podName,omitempty"`

	Namespace string `json:"namespace,omitempty"`

	Period metav1.Duration `json:"period,omitempty"`

	WithLabels map[string]string `json:"withLabels,omitempty" protobuf:"bytes,11,rep,name=withLabels"`

	WithFields map[string]string `json:"withFields,omitempty" protobuf:"bytes,11,rep,name=withLabels"`
}

// ChaosmonkeyStatus defines the observed state of Chaosmonkey
type ChaosmonkeyStatus struct {
	State      string      `json:"state,omitempty"`
	LastDelete metav1.Time `json:"lastDelete,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Chaosmonkey is the Schema for the chaosmonkeys API
type Chaosmonkey struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`

	Spec   ChaosmonkeySpec   `json:"spec"`
	Status ChaosmonkeyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ChaosmonkeyList contains a list of Chaosmonkey
type ChaosmonkeyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Chaosmonkey `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Chaosmonkey{}, &ChaosmonkeyList{})
}
