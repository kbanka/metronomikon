package helpers

import (
	"fmt"
	"github.com/applauseoss/metronomikon/config"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	"strings"
)

// Struct for handling parsing/generation of Metronome job JSON schema
type MetronomeJob struct {
	Id          string            `json:"id"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Run         struct {
		Args      []string `json:"args,omitempty"`
		Artifacts []struct {
			Uri        string `json:"uri"`
			Executable bool   `json:"executable"`
			Extract    bool   `json:"extract"`
			Cache      bool   `json:"cache"`
		} `json:"artifacts,omitempty"`
		Cmd    string  `json:"cmd,omitempty"`
		Cpus   float32 `json:"cpus"`
		Disk   float32 `json:"disk"`
		Docker struct {
			Image string `json:"image"`
		} `json:"docker"`
		Env            map[string]string `json:"env,omitempty"`
		MaxLaunchDelay int               `json:"maxLaunchDelay,omitempty"`
		Mem            float32           `json:"mem"`
		Placement      struct {
			Constraints []struct {
				Attribute string `json:"attribute,omitempty"`
				Operator  string `json:"operator,omitempty"`
				Value     string `json:"value,omitempty"`
			} `json:"constraints,omitempty"`
		} `json:"placement,omitempty"`
		User    string `json:"user,omitempty"`
		Restart struct {
			Policy                string `json:"policy,omitempty"`
			ActiveDeadlineSeconds int    `json:"activeDeadlineSeconds,omitempty"`
		} `json:"restart,omitempty"`
		Volumes []struct {
			ContainerPath string `json:"containerPath,omitempty"`
			HostPath      string `json:"hostPath,omitempty"`
			Mode          string `json:"mode,omitempty"`
		} `json:"volumes,omitempty"`
	} `json:"run"`
}

type MetronomeJobRun struct {
	CompletedAt *string  `json:"completedAt"` // we use a pointer so that we can get a null in the JSON if not populated
	CreatedAt   string   `json:"createdAt"`
	Id          string   `json:"id"`
	JobId       string   `json:"jobId"`
	Status      string   `json:"status"`
	Tasks       []string `json:"tasks"`
}

// Convert Kubernetes CronJob to Metronome format
func CronJobKubernetesToMetronome(cronJob *batchv1beta1.CronJob) *MetronomeJob {
	ret := &MetronomeJob{}
	cfg := config.GetConfig()
	ret.Id = fmt.Sprintf("%s.%s", cronJob.ObjectMeta.Namespace, cronJob.ObjectMeta.Name)
	// Metronome only supports a single container, so we grab the first one
	// XXX: make this configurable?
	container := cronJob.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
	ret.Run.Docker.Image = container.Image
	ret.Run.Args = container.Args
	if len(container.Command) > 0 {
		ret.Run.Cmd = strings.Join(container.Command, " ")
	}
	// TODO: pull values from CronJob 'resources'
	ret.Run.Mem = cfg.Metronome.JobDefaults.Memory
	ret.Run.Disk = cfg.Metronome.JobDefaults.Disk
	ret.Run.Cpus = cfg.Metronome.JobDefaults.Cpus
	return ret
}

// Convert Kubernetes Job to Metronome job run format
func JobKubernetesToMetronome(job *batchv1.Job) *MetronomeJobRun {
	ret := &MetronomeJobRun{}
	ret.Id = fmt.Sprintf("%s.%s", job.ObjectMeta.Namespace, job.ObjectMeta.Name)
	ret.JobId = fmt.Sprintf("%s.%s", job.ObjectMeta.Namespace, job.ObjectMeta.OwnerReferences[0].Name)
	ret.CreatedAt = job.ObjectMeta.CreationTimestamp.String()
	if job.Status.StartTime == nil {
		ret.Status = "STARTING"
	} else if job.Status.CompletionTime == nil {
		ret.Status = "RUNNING"
	} else if job.Status.Failed > 0 {
		ret.Status = "FAILED"
	} else {
		ret.Status = "COMPLETED"
		// We need a temp var to be able to use the address of it in the assignment below
		completionTime := job.Status.CompletionTime.String()
		ret.CompletedAt = &completionTime
	}
	ret.Tasks = make([]string, 0)
	return ret
}

// Split Metronome job ID into Kubernetes namespace and job name
func SplitMetronomeJobId(jobId string) (string, string, error) {
	parts := strings.SplitN(jobId, ".", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("failed parsing job ID '%s': job ID should consist of namespace and job name separated by a period (namespace.job_name)", jobId)
	}
	return parts[0], parts[1], nil
}
