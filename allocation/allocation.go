package allocation

import (
	"fmt"

	config "github.com/adivyoseph/kwi_nri/config"
)

func Init(p config.ProfileStruct) bool {
	//debug only
	fmt.Printf("profile dump\n")
	fmt.Printf("smt %d\n", p.Smt)
	fmt.Printf("reserved cpus %s\n", p.ReservedCpus)
	fmt.Printf("extra namespaces %s\n", p.ExtraReservedDNS)
	for nodeIndex := 0; nodeIndex < len(p.NumaNodes); nodeIndex++ {
		fmt.Printf("Node[%d] (%d)\n", nodeIndex, p.NumaNodes[nodeIndex].Index)
		for ccxIndex := 0; ccxIndex < len(p.NumaNodes[nodeIndex].Ccxs); ccxIndex++ {
			fmt.Printf("\tccx[%d] (%d) cpuSet %s\n", ccxIndex, p.NumaNodes[nodeIndex].Ccxs[ccxIndex].Index, p.NumaNodes[nodeIndex].Ccxs[ccxIndex].Cpus)
		}

	}

	return true

}
