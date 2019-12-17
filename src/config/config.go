package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type ConfigKubernetesNamespace struct {
	Whitelist []string `yaml:"whitelist"`
	Blacklist []string `yaml:"blacklist"`
}

type ConfigKubernetesJobDefaults struct {
}

type ConfigKubernetes struct {
	Namespace   ConfigKubernetesNamespace   `yaml:"namespace"`
	JobDefaults ConfigKubernetesJobDefaults `yaml:"jobDefaults"`
}

type ConfigMetronomeJobDefaults struct {
	Memory float32 `yaml:"memory"`
	Disk   float32 `yaml:"disk"`
	Cpus   float32 `yaml:"cpus"`
}

type ConfigMetronome struct {
	JobDefaults ConfigMetronomeJobDefaults `yaml:"jobDefaults"`
}

type Config struct {
	Kubernetes ConfigKubernetes `yaml:"kubernetes"`
	Metronome  ConfigMetronome  `yaml:"metronome"`
}

// Default config
var config = Config{
	Kubernetes: ConfigKubernetes{
		Namespace: ConfigKubernetesNamespace{
			Blacklist: []string{"kube-system"},
			Whitelist: []string{},
		},
		JobDefaults: ConfigKubernetesJobDefaults{},
	},
	Metronome: ConfigMetronome{
		JobDefaults: ConfigMetronomeJobDefaults{
			Memory: 256,
			Disk:   10,
			Cpus:   1.0,
		},
	},
}

// Load config from YAML file
func LoadConfig(path string) error {
	if f, err := os.Open(path); err != nil {
		return err
	} else {
		dec := yaml.NewDecoder(f)
		dec.SetStrict(true)
		if err = dec.Decode(&config); err != nil {
			return err
		}
	}
	return nil
}

// Return config singleton
func GetConfig() *Config {
	return &config
}
