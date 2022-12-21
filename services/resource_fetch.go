package services

import (
	"context"

	awsv1 "github.com/ManojDhanorkar/vm-scheduler-operator/apis/aws/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var rfLog = logf.Log.WithName("resource_fetch")

// FetchCronJob returns the CronJob resource with the name in the namespace
func FetchCronJob(name, namespace string, client client.Client) (*batchv1.CronJob, error) {
	rfLog.Info("Fetching CronJob ...")
	cronJob := &batchv1.CronJob{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, cronJob)
	return cronJob, err
}

// FetchAWSVMSchedulerCR fetches CR of AWSVMScheduler object
func FetchAWSVMSchedulerCR(name, namespace string, client client.Client) (*awsv1.AWSVMScheduler, error) {
	rfLog.Info("AWSVMScheduler Backup CR ...")
	awsvmscheduler := &awsv1.AWSVMScheduler{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, awsvmscheduler)
	return awsvmscheduler, err
}
