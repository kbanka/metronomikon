package kube

import (
	"fmt"
	v1beta1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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
