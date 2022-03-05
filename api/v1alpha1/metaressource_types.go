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

// MetaRessourceSpec defines the desired state of MetaRessource
type MetaRessourceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Fields
	// -* from:
	// -- -- apiVersion
	// -- -- kind
	// -- -- name (Todo: Alternative label/annotation matchers?)
	// -- -- fieldspec (path to field)
	// -- to:
	// -- -- fieldspec (path to field)

	// Target resource definition
	// +kubebuilder:validation:EmbeddedResource
	// +kubebuilder:pruning:PreserveUnknownFields
	// +kubebuilder:validation:Required
	Target unstructured.Unstructured `json:"target"`
	//Target []unstructured.Unstructured `json:"target"`
}

// MetaRessourceStatus defines the observed state of MetaRessource
type MetaRessourceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// MetaRessource is the Schema for the metaressources API
type MetaRessource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetaRessourceSpec   `json:"spec,omitempty"`
	Status MetaRessourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MetaRessourceList contains a list of MetaRessource
type MetaRessourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetaRessource `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetaRessource{}, &MetaRessourceList{})
}
