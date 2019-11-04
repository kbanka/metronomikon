package main

import (
	"flag"
	"github.com/applauseoss/metronomikon/api"
)

func main() {
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	a := api.New(*debug)
	a.Start()
}
