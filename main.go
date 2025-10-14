package main

import (
	"flag"
	"fmt"

	allocation "github.com/adivyoseph/kwi_nri/allocation"
	api "github.com/adivyoseph/kwi_nri/apiclient"
	config "github.com/adivyoseph/kwi_nri/config"
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

	nodeName, instanceType := api.GetNodeName()

	fmt.Printf("NodeName %s, instance type %s\n", nodeName, instanceType)

	if instanceType != "" {
		profile := config.ReadConfig(instanceType)

		if len(profile.NumaNodes) == 0 {
			log.Fatal("%v", "profile not found")
		}

		ok := allocation.Init(profile)
		if ok == false {
			log.Fatal("%v", "allocation.Init failed")
		}

	}

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
