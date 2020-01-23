package api

import (
	"fmt"

	"github.com/applauseoss/metronomikon/helpers"
	"github.com/applauseoss/metronomikon/kube"
	"github.com/gin-gonic/gin"
)

func handleGetJobRuns(c *gin.Context) {
	jobId := c.Param("jobid")
	namespace, cronJobName, err := helpers.SplitMetronomeJobId(jobId)
	if err != nil {
		JsonError(c, 500, err.Error())
		return
	}
	jobs, err := kube.GetJobsFromCronJob(namespace, cronJobName)
	if err != nil {
		JsonError(c, 500, err.Error())
		return
	}
	res := []helpers.MetronomeJobRun{}
	for _, job := range jobs {
		res = append(res, *helpers.JobKubernetesToMetronome(&job))
	}
	c.JSON(200, res)
}

func handleTriggerJobRun(c *gin.Context) {
	jobId := c.Param("jobid")
	namespace, name, err := helpers.SplitMetronomeJobId(jobId)
	if err != nil {
		JsonError(c, 500, err.Error())
		return
	}
	cronJob, err := kube.GetCronJob(namespace, name)
	if err != nil {
		JsonError(c, 404, fmt.Sprintf("cannot retrieve job: %s", err))
		return
	}
	job, err := kube.CreateJobFromCronjob(cronJob)
	if err != nil {
		JsonError(c, 500, fmt.Sprintf("cannot run job: %s", err))
		return
	}
	c.JSON(201, helpers.JobKubernetesToMetronome(job))
}

func handleGetJobRun(c *gin.Context) {
	// This API endpoint also takes a 'jobid' parameter, but we don't need it
	// because each job run has a unique name already in kubernetes
	runId := c.Param("runid")
	namespace, name, err := helpers.SplitMetronomeJobId(runId)
	if err != nil {
		JsonError(c, 500, err.Error())
		return
	}
	job, err := kube.GetJob(namespace, name)
	if err != nil {
		JsonError(c, 404, fmt.Sprintf("cannot get job: %s", err))
		return
	}
	c.JSON(200, helpers.JobKubernetesToMetronome(job))
}

func handleStopJobRun(c *gin.Context) {
	c.String(200, "TODO")
}
