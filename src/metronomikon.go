package main

import (
	"flag"
	"fmt"
	"github.com/applauseoss/metronomikon/api"
	"github.com/applauseoss/metronomikon/config"
	"github.com/applauseoss/metronomikon/kube"
	"log"
	"os"
)

var defaultConfigFile = "/etc/metronomikon/config.yaml"

func main() {
	var debug bool
	var configFile string
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.StringVar(&configFile, "config", "", fmt.Sprintf("Path to config file (defaults to %s, if it exists)", defaultConfigFile))
	flag.Parse()
	// Use default config file path if none was provided and it exists
	if configFile == "" {
		if _, err := os.Stat(defaultConfigFile); err == nil {
			configFile = defaultConfigFile
		}
	}
	// Load config file if specified
	if configFile != "" {
		if err := config.LoadConfig(configFile); err != nil {
			log.Fatalf("Failed to load config file: %s", err)
		}
	}
	if err := kube.TestClientConnection(); err != nil {
		log.Fatalf("Failed to connect to Kubernetes API: %s", err)
	} else {
		if debug {
			log.Print("Successfully initialized kubernetes client")
		}
	}
	a := api.New(debug)
	_ = a.Start()
}
