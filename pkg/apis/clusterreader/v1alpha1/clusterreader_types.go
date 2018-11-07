package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ClusterReaderSpec defines the desired state of ClusterReader
type ClusterReaderSpec struct {
	Readers []string `json:"readers"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// ClusterReaderStatus defines the observed state of ClusterReader
type ClusterReaderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterReader is the Schema for the clusterreaders API
// +k8s:openapi-gen=true
type ClusterReader struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterReaderSpec   `json:"spec,omitempty"`
	Status ClusterReaderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterReaderList contains a list of ClusterReader
type ClusterReaderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterReader `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterReader{}, &ClusterReaderList{})
}
