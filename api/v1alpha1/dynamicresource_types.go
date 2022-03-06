/*
Copyright 2022 Tilman Eggers.

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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DynamicResourceSpec defines the desired state of DynamicResource
type DynamicResourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Target resource definition
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Required
	Target unstructured.Unstructured `json:"target"`

	// +kubebuilder:validation:Optional
	Transformations []DynamicResourceTransformation `json:"transformations"`
}

type DynamicResourceTransformation struct {
	FieldFrom ExternalFieldRef `json:"fieldFrom"`

	TargetField string `json:"targetField"`
}

// ExternalFieldRef Reference to a field of any resource on the cluster
type ExternalFieldRef struct {
	metav1.TypeMeta `json:",inline"`

	// Todo: Add more advanced resource matchers, e.g. label- and field-based matching
	// Name of the target resource
	Name string `json:"name"`

	// Todo: Add more advanced field matchers
	// Selector for the field to copy the data from
	FieldSpec string `json:"fieldSpec"`
}

// DynamicResourceStatus defines the observed state of DynamicResource
type DynamicResourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DynamicResource is the Schema for the dynamicresources API
type DynamicResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DynamicResourceSpec   `json:"spec,omitempty"`
	Status DynamicResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DynamicResourceList contains a list of DynamicResource
type DynamicResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DynamicResource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DynamicResource{}, &DynamicResourceList{})
}
