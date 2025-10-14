package apiclient

import (
	"context"
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetNodeName() (string, string) {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("Error getting hostname: %v\n", err)
		return "", ""
	}
	fmt.Printf("GetNodeName: Hostname: %s\n", hostname)

	//inside  POD
	var config *rest.Config
	var clientset *kubernetes.Clientset
	var node *v1.Node
	var nodeName string
	config, err = rest.InClusterConfig()
	if err == nil {
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			return "", ""
		}
	} else {
		fmt.Printf("GetNodeName: InClusterConfig failed\n")
		kubeconfig := "/etc/kubernetes/admin.conf"
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			fmt.Printf("GetNodeName: BuildConfigFromFlags failed\n")
			return "", ""
		}
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			fmt.Printf("GetNodeName: NewForConfig failed\n")
			return "", ""
		}
	}
	nodeName = hostname
	node, err = clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("GetNodeName: Nodes().Get failed\n")
		return "", ""
	}
	for label, value := range node.Labels {
		fmt.Printf("GetNodeName: label %s value %s\n", label, value)
	}
	if val, ok := node.Labels["kubernetes.io/hostname"]; ok {
		hostname = val
		fmt.Printf("GetNodeName: Hostname from API: %s\n", hostname)
	}
	var instanceType string
	if val, ok := node.Labels["kubernetes.io/azuretype"]; ok {
		instanceType = val

		fmt.Printf("GetNodeName: instance type from API: %s\n", instanceType)
	}

	return hostname, instanceType
}
