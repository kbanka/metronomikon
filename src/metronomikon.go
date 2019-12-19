package main

import (
	"flag"
	"github.com/applauseoss/metronomikon/api"
	"github.com/applauseoss/metronomikon/config"
	"github.com/applauseoss/metronomikon/kube"
	"log"
)

func main() {
	var debug bool
	var configFile string
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.StringVar(&configFile, "config", "/etc/metronomikon/config.yaml", "Path to config file")
	flag.Parse()
	if err := config.LoadConfig(configFile); err != nil {
		log.Fatalf("Failed to load config file: %s", err)
	}
	if err := kube.TestClientConnection(); err != nil {
		log.Fatalf("Failed to connect to Kubernetes API: %s", err)
	} else {
		if debug {
			log.Print("Successfully initialized kubernetes client")
		}
	}
	a := api.New(debug)
	a.Start()
}
