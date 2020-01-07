package kube

import (
	"fmt"
	v1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Delete specific CronJob
// Returns nil and error when job doesn't exist
// Return job and error when deleting job failed
// Return job and nil when deletion succeeded
func DeleteCronJob(namespace string, name string) (*v1beta1.CronJob, error) {
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
func GetCronJobs(namespace string) ([]v1beta1.CronJob, error) {
	jobs, err := client.BatchV1beta1().CronJobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not list CronJobs: %s", err)
	}
	return jobs.Items, nil
}

// Return specific CronJob
func GetCronJob(namespace string, name string) (*v1beta1.CronJob, error) {
	job, err := client.BatchV1beta1().CronJobs(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get CronJob: %s", err)
	}
	return job, nil
}
