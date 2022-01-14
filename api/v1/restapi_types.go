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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RestAPISpec defines the desired state of RestAPI
type RestAPISpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Frontend is a struct
	Frontend FrontendSpec `json:"frontend"`

	// RedisName is where the Redis is
	RedisName string `json:"redisName,omitempty"`
}

type FrontendSpec struct {
	//+optional

	// cpu/memory resources
	Resources corev1.ResourceRequirements `json:"resources"`

	//+optional
	//+kubebuilder:default=8080
	//+kubebuilder:validation:Minimum=0

	// ServerPort is what port we serve the API on
	ServingPort int32 `json:"serverPort,omitempty"`

	//+optional
	//+kubebuilder:default=1
	//+kubebuilder:validation:Minimum=0

	// How many replicas
	Replicas *int32 `json:"replicas,omitempty"`
}

// RestAPIStatus defines the observed state of RestAPI
type RestAPIStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// URL is a string of... stuff
	URL string `json:"url,omitempty"`

	// Conditions of the RestAPI
	Conditions []StatusCondition `json:"conditions,omitempty"`
}

type ConditionStatus string

var (
	ConditionStatusHealthy   ConditionStatus = "Healthy"
	ConditionStatusUnhealthy ConditionStatus = "Unhealthy"
	ConditionStatusUnknown   ConditionStatus = "Unknown"
)

type StatusCondition struct {
	Type   string          `json:"type"`
	Status ConditionStatus `json:"status"`
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// +optional
	Reason string `json:"reason,omitempty"`
	// +optional
	Message string `json:"message,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:priority=0,name=URL,type=string,JSONPath=".status.url",description="GuestBook Frontend URL",format=""
//+kubebuilder:printcolumn:priority=0,name=Deployment,type=string,JSONPath=".status.conditions[?(@.type==\"DeploymentUpToDate\")].status",description="Is the Deployment Up-To-Date",format=""
//+kubebuilder:printcolumn:priority=0,name=Service,type=string,JSONPath=".status.conditions[?(@.type==\"ServiceUpToDate\")].status",description="Is the Service Up-To-Date",format=""

// RestAPI is the Schema for the restapis API
type RestAPI struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RestAPISpec   `json:"spec,omitempty"`
	Status RestAPIStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RestAPIList contains a list of RestAPI
type RestAPIList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RestAPI `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RestAPI{}, &RestAPIList{})
}
