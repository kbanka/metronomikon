package kube

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// Delete specific CronJob
// Returns nil and error when job doesn't exist
// Return job and error when deleting job failed
// Return job and nil when deletion succeeded
func DeleteCronJob(namespace string, name string) (*batchv1beta1.CronJob, error) {
	job, err := GetCronJob(namespace, name)
	if err != nil {
		return nil, err
	}
	err = client.BatchV1beta1().CronJobs(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return job, err
	}
	return job, nil
}

// Return all CronJobs in a namespace
func GetCronJobs(namespace string) ([]batchv1beta1.CronJob, error) {
	jobs, err := client.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list CronJobs: %s", err)
	}
	return jobs.Items, nil
}

// Return specific CronJob
func GetCronJob(namespace string, name string) (*batchv1beta1.CronJob, error) {
	job, err := client.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get CronJob: %s", err)
	}
	return job, nil
}

// Create Job from CronJob
func CreateJobFromCronjob(cronJob *batchv1beta1.CronJob) (*batchv1.Job, error) {
	// This duplicates the logic used by kubectl to create a Job from a CronJob
	annotations := make(map[string]string)
	annotations["cronjob.kubernetes.io/instantiate"] = "manual"
	for k, v := range cronJob.Spec.JobTemplate.Annotations {
		annotations[k] = v
	}

	jobDef := &batchv1.Job{
		TypeMeta: metav1.TypeMeta{APIVersion: batchv1.SchemeGroupVersion.String(), Kind: "Job"},
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%s-%d", cronJob.ObjectMeta.Name, time.Now().Unix()),
			Annotations: annotations,
			Labels:      cronJob.Spec.JobTemplate.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(cronJob, appsv1.SchemeGroupVersion.WithKind("CronJob")),
			},
		},
		Spec: cronJob.Spec.JobTemplate.Spec,
	}

	if job, err := client.BatchV1().Jobs(cronJob.ObjectMeta.Namespace).Create(jobDef); err != nil {
		return nil, err
	} else {
		return job, nil
	}
}

// Return specific Job
func GetJob(namespace string, name string) (*batchv1.Job, error) {
	job, err := client.BatchV1().Jobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get Job: %s", err)
	}
	return job, nil
}
