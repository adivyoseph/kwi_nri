package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type configCpusStruct struct {
	Cpus string `yaml:"cpus"`
}
type configNumaStruct struct {
	Ccxs map[string]configCpusStruct `yaml:"ccxs"`
}

type configProfileStruct struct {
	Smt             string                      `yaml:"smt"`
	Class           string                      `yaml:"class"`
	Schedule        string                      `yaml:"schedule"`
	Exclusive       bool                        `yaml:"exclusive"`
	ReservedCpus    string                      `yaml:"reservedcpus"`
	ExtraReservedNS string                      `yaml:"etrareservedns"`
	Sidecars        string                      `yaml:"sidecars"`
	FillOrder       string                      `yaml:"fillorder"`
	Annotations     string                      `yaml:"annotations"`
	Numas           map[string]configNumaStruct `yaml:"numas"`
}

type amdConfigStruct struct {
	Kind        string `yaml:"kind"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
	Profiles    map[string]configProfileStruct
}

type ProfileCcxStruct struct {
	Index int
	Cpus  string
}

type ProfileNodeStruct struct {
	Index int
	Ccxs  []ProfileCcxStruct
}

type ProfileStruct struct {
	Smt              int
	ReservedCpus     string
	ExtraReservedDNS string
	NumaNodes        []ProfileNodeStruct
}

func removeAlphaPrefix(s string) int {
	startIndex := 0
	for i, r := range s {
		if r >= '0' && r <= '9' {
			startIndex = i
			break
		}
		// If all characters are letters, return an empty string
		if i == len(s)-1 {
			return -1
		}
	}
	index, err := strconv.Atoi(s[startIndex:])
	if err != nil {
		fmt.Printf("removeAlphaPrefix %s not valid", s)
		log.Fatal(err)
	}

	return index
}

func ReadConfig(profile string) ProfileStruct {

	file, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	c := amdConfigStruct{}

	// Unmarshal our input YAML file into empty Car (var c)
	if err := yaml.Unmarshal(file, &c); err != nil {
		log.Fatal(err)
	}

	result := ProfileStruct{}
	for name, content := range c.Profiles {
		if name == profile {

			if content.Smt == "on" {
				result.Smt = 2
			}
			result.ReservedCpus = content.ReservedCpus
			result.ExtraReservedDNS = content.ExtraReservedNS
			for numaName, node := range content.Numas {
				newNode := ProfileNodeStruct{}
				newNode.Index = removeAlphaPrefix(numaName)
				fmt.Printf("nodename %s index %d\n", numaName, newNode.Index)
				for ccxName, ccx := range node.Ccxs {
					newCcx := ProfileCcxStruct{}
					newCcx.Index = removeAlphaPrefix(ccxName)
					fmt.Printf("ccxname %s index %d\n", ccxName, newCcx.Index)
					newCcx.Cpus = ccx.Cpus
					newNode.Ccxs = append(newNode.Ccxs, newCcx)
				}
				//sort
				result.NumaNodes = append(result.NumaNodes, newNode)
			}
			//sort
			for i := 0; i < len(result.NumaNodes)-1; i++ {
				for j := 0; j < len(result.NumaNodes)-1; j++ {
					if result.NumaNodes[j].Index > result.NumaNodes[j+1].Index {
						temp := result.NumaNodes[j]
						result.NumaNodes[j] = result.NumaNodes[j+1]
						result.NumaNodes[j+1] = temp
					}
				}
			}
			for n := 0; n < len(result.NumaNodes); n++ {
				for i := 0; i < len(result.NumaNodes[n].Ccxs)-1; i++ {
					for j := 0; j < len(result.NumaNodes[n].Ccxs)-1; j++ {
						if result.NumaNodes[n].Ccxs[j].Index > result.NumaNodes[n].Ccxs[j+1].Index {
							temp := result.NumaNodes[n].Ccxs[j]
							result.NumaNodes[n].Ccxs[j] = result.NumaNodes[n].Ccxs[j+1]
							result.NumaNodes[n].Ccxs[j+1] = temp
						}
					}
				}
			}
			break
		}
	}

	return result
}
