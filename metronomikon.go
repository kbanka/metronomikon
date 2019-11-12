package main

import (
	"flag"
	"github.com/applauseoss/metronomikon/api"
	"github.com/applauseoss/metronomikon/kube"
	"log"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.Parse()
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
