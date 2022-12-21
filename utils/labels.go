package utils

import (
	awsv1 "github.com/ManojDhanorkar/vm-scheduler-operator/apis/aws/v1"
)

func AWSVMSchedulerLabels(v *awsv1.AWSVMScheduler, tier string) map[string]string {
	return map[string]string{
		"app":               "AWSVMScheduler",
		"AWSVMScheduler_cr": v.Name,
		"tier":              tier,
	}
}
