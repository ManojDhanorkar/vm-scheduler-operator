package utils

import (
	awsv1 "github.com/ManojDhanorkar/vm-scheduler-operator/apis/aws/v1"
)

const (
	image         = "quay.io/manoj_dhanorkar/awsvmscheduler:v1.0"
	instanceIds   = "i-00db611adab6f713d,i-00db611adab6f713e"
	startSchedule = "0 9 * * 1-5"
	stopSchedule  = "0 18 * * 1-5"
)

type DefaultAWSVMSchedulerConfig struct {
	Image         string `json:"image"`
	InstanceIds   string `json:"instanceIds"`
	StartSchedule string `json:"startSchedule"`
	StopSchedule  string `json:"stopSchedule"`
}

var defaultAWSVMSchedulerConfig = NewDefaultAWSVMSchedulerConfig()

// AddAWSVMSchedulerMandatorySpecs will add the specs which are mandatory for AWSVMScheduler CR in the case them
// not be applied
func AddAWSVMSchedulerMandatorySpecs(awsVMScheduler *awsv1.AWSVMScheduler) {

	/*
	 AWSVMScheduler Container
	*/

	if awsVMScheduler.Spec.StartSchedule == "" {
		awsVMScheduler.Spec.StartSchedule = defaultAWSVMSchedulerConfig.StartSchedule
	}

	if awsVMScheduler.Spec.StopSchedule == "" {
		awsVMScheduler.Spec.StopSchedule = defaultAWSVMSchedulerConfig.StopSchedule
	}

}

func NewDefaultAWSVMSchedulerConfig() *DefaultAWSVMSchedulerConfig {
	return &DefaultAWSVMSchedulerConfig{
		Image:         image,
		InstanceIds:   instanceIds,
		StartSchedule: startSchedule,
		StopSchedule:  stopSchedule,
	}
}
