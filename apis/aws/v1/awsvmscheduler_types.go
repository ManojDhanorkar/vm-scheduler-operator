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

// AWSVMSchedulerSpec defines the desired state of AWSVMScheduler
type AWSVMSchedulerSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Comma separated list of AWS instance ids which will be scheduled by CR
	InstanceIds string `json:"instanceIds"`

	// Schedule period for the CronJob.
	// This spec allows you to setup the start schedule for VM
	StartSchedule string `json:"startSchedule"`

	// Schedule period for the CronJob.
	// This spec allows you to setup the stop schedule for VM
	StopSchedule string `json:"stopSchedule"`

	// This spec allows you to supply image name which will start/stop VM
	Image string `json:"image"`
}

// AWSVMSchedulerStatus defines the observed state of AWSVMScheduler
type AWSVMSchedulerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Schedule period for the CronJob.
	// This will show the status of stop action for the VM(s)
	VMStopStatus string `json:"vmStopStatus"`

	// Schedule period for the CronJob.
	// This will show the status of start action for the VM(s)
	VMStartStatus string `json:"vmStartStatus"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// AWSVMScheduler is the Schema for the awsvmschedulers API
type AWSVMScheduler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AWSVMSchedulerSpec   `json:"spec,omitempty"`
	Status AWSVMSchedulerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AWSVMSchedulerList contains a list of AWSVMScheduler
type AWSVMSchedulerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSVMScheduler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AWSVMScheduler{}, &AWSVMSchedulerList{})
}
