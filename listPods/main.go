package main

import (
	"context"
	"flag"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/devjpt23/.kube/config", "location to your kubeconfig file")
	namespace := flag.String("namespace", "default", "Kubernetes namespace")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {

	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {

	}
	// listiing Pods 
	pods, err := clientSet.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {

	}
	fmt.Printf("These are the Pods from %s namespace\n...", *namespace)
	fmt.Println("...###...")
	for _, pod:= range pods.Items{
		fmt.Println(pod.Name)
	}
	// Listing Nodes
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {

	}
	fmt.Println("These are the nodes...")
	fmt.Println("...###...")

	for _, node:= range nodes.Items{
		fmt.Println(node.Name)
	}
}
