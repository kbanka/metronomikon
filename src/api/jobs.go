package api

import (
	"fmt"

	"github.com/applauseoss/metronomikon/helpers"
	"github.com/applauseoss/metronomikon/kube"
	"github.com/gin-gonic/gin"
)

func handleGetJobs(c *gin.Context) {
	ret := []*helpers.MetronomeJob{}
	namespaces, err := kube.GetNamespaces()
	if err != nil {
		JsonError(c, 500, fmt.Sprintf("Failed to list namespaces: %s", err))
		return
	}
	for _, namespace := range namespaces {
		jobs, err := kube.GetCronJobs(namespace)
		if err != nil {
			JsonError(c, 500, fmt.Sprintf("Failed to list jobs: %s", err))
			return
		}
		for _, job := range jobs {
			metronomeJob := helpers.CronJobKubernetesToMetronome(&job)

			var ginErrorMessage *helpers.GinErrorMessage
			embed := c.Query("embed")
			if embed != "" {
				metronomeJob, ginErrorMessage = helpers.HandleGetJobEmbed(embed, metronomeJob)
				if ginErrorMessage != nil {
					JsonError(c, ginErrorMessage.HTTPCode, ginErrorMessage.Message)
				}
			}

			ret = append(ret, metronomeJob)
		}
	}
	c.JSON(200, ret)
}

func handleCreateJob(c *gin.Context) {
	c.String(200, "TODO")
}

func handleGetJob(c *gin.Context) {
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

	metronomeJob := helpers.CronJobKubernetesToMetronome(cronJob)
	var ginErrorMessage *helpers.GinErrorMessage
	embed := c.Query("embed")
	if embed != "" {
		metronomeJob, ginErrorMessage = helpers.HandleGetJobEmbed(embed, metronomeJob)
		if ginErrorMessage != nil {
			JsonError(c, ginErrorMessage.HTTPCode, ginErrorMessage.Message)
		}
	}

	c.JSON(200, metronomeJob)
}

func handleUpdateJob(c *gin.Context) {
	c.String(200, "TODO")
}

func handleDeleteJob(c *gin.Context) {
	jobId := c.Param("jobid")
	namespace, name, err := helpers.SplitMetronomeJobId(jobId)
	if err != nil {
		JsonError(c, 500, err.Error())
		return
	}

	job, err := kube.DeleteCronJob(namespace, name)
	if job == nil {
		c.JSON(404, gin.H{"message": fmt.Sprintf("Job '%s' does not exist", jobId)})
		return
	} else if err != nil {
		JsonError(c, 500, fmt.Sprintf("failed to delete job: %s", err))
		return
	}

	tmp_job := helpers.CronJobKubernetesToMetronome(job)
	c.JSON(200, tmp_job)
}
