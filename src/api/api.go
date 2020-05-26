package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type Api struct {
	engine *gin.Engine
}

func New(debug bool) *Api {
	api := &Api{}
	api.init(debug)
	return api
}

func (a *Api) init(debug bool) {
	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.DisableConsoleColor()
	a.engine = gin.New()
	a.engine.Use(gin.Recovery())
	a.engine.Use(gin.LoggerWithFormatter(accessLogger))
	// Healthcheck
	a.engine.GET("/ping", handlePing)
	// Jobs
	a.engine.GET("/v1/jobs", handleGetJobs)
	a.engine.POST("/v1/jobs", handleCreateJob)
	a.engine.GET("/v1/jobs/:jobid", handleGetJob)
	a.engine.PUT("/v1/jobs/:jobid", handleUpdateJob)
	a.engine.DELETE("/v1/jobs/:jobid", handleDeleteJob)
	// Job schedules
	a.engine.GET("/v1/jobs/:jobid/schedules", handleGetJobSchedules)
	a.engine.POST("/v1/jobs/:jobid/schedules", handleCreateJobSchedule)
	a.engine.GET("/v1/jobs/:jobid/schedules/:scheduleid", handleGetJobSchedule)
	a.engine.PUT("/v1/jobs/:jobid/schedules/:scheduleid", handleUpdateJobSchedule)
	a.engine.DELETE("/v1/jobs/:jobid/schedules/:scheduleid", handleDeleteJobSchedule)
	// Job runs
	a.engine.GET("/v1/jobs/:jobid/runs", handleGetJobRuns)
	a.engine.POST("/v1/jobs/:jobid/runs", handleTriggerJobRun)
	a.engine.GET("/v1/jobs/:jobid/runs/:runid", handleGetJobRun)
	a.engine.POST("/v1/jobs/:jobid/runs/:runid/actions/stop", handleStopJobRun)
	// Metrics
	a.engine.GET("/v1/metrics", handleGetMetrics)
}

func (a *Api) Start() {
	a.engine.Run()
}

// Output HTTP 500 with JSON body containing error message
func JsonError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"message": message})
}

// Format access log in JSON format
func accessLogger(param gin.LogFormatterParams) string {
	logEntry := gin.H{
		"type":          "access",
		"client_ip":     param.ClientIP,
		"timestamp":     param.TimeStamp.Format(time.RFC1123),
		"method":        param.Method,
		"path":          param.Path,
		"proto":         param.Request.Proto,
		"status_code":   param.StatusCode,
		"latency":       param.Latency,
		"user_agent":    param.Request.UserAgent(),
		"error_message": param.ErrorMessage,
	}
	ret, _ := json.Marshal(logEntry)
	return string(ret) + "\n"
}
