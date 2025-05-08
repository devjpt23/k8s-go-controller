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
		fmt.Printf("Error in building config: %s\n", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error in creating clientset: %s\n", err.Error())

	}
	// listiing Pods
	pods, err := clientSet.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error in listing pods: %s\n", err.Error())
	}
	fmt.Printf("These are the Pods from %s namespace\n...", *namespace)
	for _, pod:= range pods.Items{
		fmt.Println(pod.Name)
	}
	

	// Listing Nodes
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error in listing the nodes: %s\n", err.Error())
	}
	fmt.Println("These are the nodes...")
	

	for _, node:= range nodes.Items{
		fmt.Println(node.Name)
	}
	// Listing Deployments
	deploys, err := clientSet.AppsV1().Deployments(*namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil{
		fmt.Printf("There was an error in listing the deployments: %s", err)
	}
	fmt.Println("These are the deployments")


	for _, deploys := range deploys.Items{
		fmt.Println(deploys.Name)
	}
}
