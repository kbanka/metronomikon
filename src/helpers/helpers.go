package helpers

import (
	"fmt"
	"github.com/applauseoss/metronomikon/config"
	v1beta1 "k8s.io/api/batch/v1beta1"
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

// Convert Kubernetes CronJob to Metronome format
func JobKubernetesToMetronome(job *v1beta1.CronJob) *MetronomeJob {
	ret := &MetronomeJob{}
	cfg := config.GetConfig()
	ret.Id = fmt.Sprintf("%s.%s", job.ObjectMeta.Namespace, job.ObjectMeta.Name)
	// Metronome only supports a single container, so we grab the first one
	// XXX: make this configurable?
	container := job.Spec.JobTemplate.Spec.Template.Spec.Containers[0]
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

// Split Metronome job ID into Kubernetes namespace and job name
func SplitMetronomeJobId(jobId string) (string, string, error) {
	parts := strings.SplitN(jobId, ".", 2)
	if len(parts) < 2 {
		return "", "", fmt.Errorf("failed parsing job ID '%s': job ID should consist of namespace and job name separated by a period (namespace.job_name)", jobId)
	}
	return parts[0], parts[1], nil
}
