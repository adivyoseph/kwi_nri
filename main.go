package main

import (
	"flag"

	policy "github.com/adivyoseph/kwi_nri/policy"
	agent "github.com/containers/nri-plugins/pkg/agent"
	logger "github.com/containers/nri-plugins/pkg/log"
	resmgr "github.com/containers/nri-plugins/pkg/resmgr/main"
)

var (
	log = logger.Default()
)

func main() {
	flag.Parse()

	agt, err := agent.New(agent.TemplateConfigInterface())
	if err != nil {
		log.Fatal("%v", err)
	}

	mgr, err := resmgr.New(agt, policy.New())
	if err != nil {
		log.Fatalf("%v", err)
	}

	if err := mgr.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
