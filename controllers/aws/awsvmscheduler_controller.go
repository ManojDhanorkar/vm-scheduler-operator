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

package aws

import (
	"context"
	"reflect"

	awsv1 "github.com/ManojDhanorkar/vm-scheduler-operator/apis/aws/v1"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	awsutils "github.com/ManojDhanorkar/vm-scheduler-operator/utils"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AWSVMSchedulerReconciler reconciles a AWSVMScheduler object
type AWSVMSchedulerReconciler struct {
	Client client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=aws.xyzcompany.com,resources=awsvmschedulers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AWSVMScheduler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AWSVMSchedulerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	reqLogger := r.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling AWSVMScheduler")

	// Fetch the AWSVMScheduler CR
	//awsVMScheduler, err := services.FetchAWSVMSchedulerCR(req.Name, req.Namespace)

	// Fetch the AWSVMScheduler instance
	awsVMScheduler := &awsv1.AWSVMScheduler{}
	err := r.Client.Get(ctx, req.NamespacedName, awsVMScheduler)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			reqLogger.Info("awsVMScheduler resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		reqLogger.Error(err, "Failed to get awsVMScheduler.")
		return ctrl.Result{}, err
	}

	// Add const values for mandatory specs ( if left blank)
	// log.Info("Adding awsVMScheduler mandatory specs")
	// utils.AddBackupMandatorySpecs(awsVMScheduler)

	// Check if the CronJob already exists, if not create a new one

	found := &batchv1.CronJob{}
	err = r.Client.Get(ctx, types.NamespacedName{Name: awsVMScheduler.Name, Namespace: awsVMScheduler.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new CronJob
		cron := r.cronJobForAWSVMScheduler(awsVMScheduler)
		reqLogger.Info("Creating a new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
		err = r.Client.Create(ctx, cron)
		if err != nil {
			reqLogger.Error(err, "Failed to create new CronJob", "CronJob.Namespace", cron.Namespace, "CronJob.Name", cron.Name)
			return ctrl.Result{}, err
		}
		// Cronjob created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		reqLogger.Error(err, "Failed to get Cronjob")
		return ctrl.Result{}, err
	}

	//if err := r.createResources(awsVMScheduler, req); err != nil {
	//	log.Error(err, "Failed to create and update the secondary resource required for the AWSVMScheduler CR")
	//	return ctrl.Result{}, err
	//}

	// Ensure the cron inputs are the same as the spec.
	// TODO : Add support for instanceIds and stopSchedule

	// Check for any updates for redeployment
	applyChange := false

	// Ensure image name is correct, update image if required
	instanceIds := awsVMScheduler.Spec.InstanceIds
	startSchedule := awsVMScheduler.Spec.StartSchedule
	image := awsVMScheduler.Spec.Image

	var currentImage string = ""
	var currentStartSchedule string = ""
	var currentInstanceIds string = ""

	// Check schedule
	if found.Spec.Schedule != "" {
		currentStartSchedule = found.Spec.Schedule
	}

	if startSchedule != currentStartSchedule {
		found.Spec.Schedule = currentStartSchedule
		applyChange = true
	}

	// Check image
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentImage = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image
	}

	if image != currentImage {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Image = currentImage
		applyChange = true
	}

	// Check instanceIds
	if found.Spec.JobTemplate.Spec.Template.Spec.Containers != nil {
		currentInstanceIds = found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Name
	}

	if instanceIds != currentInstanceIds {
		found.Spec.JobTemplate.Spec.Template.Spec.Containers[0].Env[0].Name = currentInstanceIds
		applyChange = true
	}

	if applyChange {
		err = r.Client.Update(ctx, found)
		if err != nil {
			reqLogger.Error(err, "Failed to update CronJob", "CronJob.Namespace", found.Namespace, "CronJob.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the AWSVMScheduler status
	// TODO: Define what needs to be added in status. Currently adding just instanceIds
	if !reflect.DeepEqual(instanceIds, awsVMScheduler.Status.VMStartStatus) ||
		!reflect.DeepEqual(instanceIds, awsVMScheduler.Status.VMStopStatus) {
		awsVMScheduler.Status.VMStartStatus = instanceIds
		awsVMScheduler.Status.VMStopStatus = instanceIds
		err := r.Client.Status().Update(ctx, awsVMScheduler)
		if err != nil {
			reqLogger.Error(err, "Failed to update awsVMScheduler status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AWSVMSchedulerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&awsv1.AWSVMScheduler{}).
		Owns(&batchv1.CronJob{}).
		Complete(r)
}

// CronJob Spec
func (r *AWSVMSchedulerReconciler) cronJobForAWSVMScheduler(awsVMScheduler *awsv1.AWSVMScheduler) *batchv1.CronJob {

	cron := &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name:      awsVMScheduler.Name,
			Namespace: awsVMScheduler.Namespace,
			Labels:    awsutils.AWSVMSchedulerLabels(awsVMScheduler, "awsVMScheduler"),
		},
		Spec: batchv1.CronJobSpec{
			Schedule: awsVMScheduler.Spec.StartSchedule,
			// TODO: Add Stop schedule
			//Schedule:  awsVMScheduler.Spec.StopSchedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{{
								Name:  awsVMScheduler.Name,
								Image: awsVMScheduler.Spec.Image,
								Env: []corev1.EnvVar{
									{
										Name:  "InstanceIds",
										Value: awsVMScheduler.Spec.InstanceIds,
									},
									{
										Name: "AWS_ACCESS_KEY_ID",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "aws-secret",
												},
												Key: "aws-access-key-id",
											},
										},
									},
									{
										Name: "AWS_SECRET_ACCESS_KEY",
										ValueFrom: &corev1.EnvVarSource{
											SecretKeyRef: &corev1.SecretKeySelector{
												LocalObjectReference: corev1.LocalObjectReference{
													Name: "aws-secret",
												},
												Key: "aws-secret-access-key",
											},
										},
									}},
							}},
						},
					},
				},
			},
		},
	}
	// Set awsVMScheduler instance as the owner and controller
	ctrl.SetControllerReference(awsVMScheduler, cron, r.Scheme)
	return cron
}

// // createResources will create and update the  resource which are required
// func (r *AWSVMSchedulerReconciler) createResources(awsVMScheduler *awsv1.AWSVMScheduler, request ctrl.Request) error {

// 	log := r.Log.WithValues("AWSVMScheduler", request.NamespacedName)
// 	log.Info("Creating   resources ...")

// 	// Check if the cronJob is created, if not create one
// 	if err := r.createCronJob(awsVMScheduler); err != nil {
// 		log.Error(err, "Failed to create the CronJob")
// 		return err
// 	}

// 	return nil
// }

// // Check if the cronJob is created, if not create one
// func (r *AWSVMSchedulerReconciler) createCronJob(awsVMScheduler *awsv1.AWSVMScheduler) error {
// 	if _, err := services.FetchCronJob(awsVMScheduler.Name, awsVMScheduler.Namespace); err != nil {
// 		if err := r.client.Create(context.TODO(), resources.NewAWSVMSchedulerCronJob(awsVMScheduler)); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
