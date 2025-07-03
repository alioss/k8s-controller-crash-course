/*
Copyright 2024.

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

// AtlasAppSpec defines the desired state of AtlasApp
type AtlasAppSpec struct {
	// Environment specifies the deployment environment (dev, stage, prod)
	Environment string `json:"environment"`

	// Version specifies the application version to deploy
	Version string `json:"version"`

	// MigrationId specifies the database migration version
	MigrationId int `json:"migrationId"`

	// Replicas specifies the number of replicas to deploy
	Replicas int32 `json:"replicas,omitempty"`

	// AutoPromote enables automatic promotion to next environment
	AutoPromote bool `json:"autoPromote,omitempty"`

	// NextEnvironment specifies the next environment for promotion
	NextEnvironment string `json:"nextEnvironment,omitempty"`

	// RequireApproval requires manual approval for deployment
	RequireApproval bool `json:"requireApproval,omitempty"`

	// HealthCheckPath specifies the health check endpoint
	HealthCheckPath string `json:"healthCheckPath,omitempty"`
}

// AtlasAppStatus defines the observed state of AtlasApp
type AtlasAppStatus struct {
	// Phase represents the current phase of the application
	Phase string `json:"phase,omitempty"`

	// Ready indicates if the application is ready and healthy
	Ready bool `json:"ready,omitempty"`

	// ReadyReplicas indicates the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// TotalReplicas indicates the total number of replicas
	TotalReplicas int32 `json:"totalReplicas,omitempty"`

	// LastUpdate indicates when the deployment was last updated
	LastUpdate *metav1.Time `json:"lastUpdate,omitempty"`

	// ApprovalRequired indicates if manual approval is needed
	ApprovalRequired bool `json:"approvalRequired,omitempty"`

	// PromotionPending indicates if promotion to next env is pending
	PromotionPending bool `json:"promotionPending,omitempty"`

	// Conditions represents the current conditions of the application
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Message provides additional information about the current state
	Message string `json:"message,omitempty"`
}

// AtlasApp defines an Atlas application deployment
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Environment",type="string",JSONPath=".spec.environment"
//+kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version"
//+kubebuilder:printcolumn:name="Migration",type="integer",JSONPath=".spec.migrationId"
//+kubebuilder:printcolumn:name="Phase",type="string",JSONPath=".status.phase"
//+kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
//+kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".status.readyReplicas"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type AtlasApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AtlasAppSpec   `json:"spec,omitempty"`
	Status AtlasAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AtlasAppList contains a list of AtlasApp
type AtlasAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AtlasApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AtlasApp{}, &AtlasAppList{})
}
