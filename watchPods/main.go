
package main

import (
    "context"
    "flag"
    "log"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
)

func main() {
    
    kubeconfig := flag.String("kubeconfig", "", "absolute path to the kubeconfig file (local only)")
    namespace  := flag.String("namespace", "default", "Kubernetes namespace")
    flag.Parse()

    
    config, err := rest.InClusterConfig()
    if err != nil {
        log.Printf("InClusterConfig failed: %v", err)
        
        if *kubeconfig == "" {
            log.Fatalf("no kubeconfig provided and not running in-cluster")
        }
        config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
        if err != nil {
            log.Fatalf("BuildConfigFromFlags failed: %v", err)
        }
    }

    
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("NewForConfig failed: %v", err)
    }

    
    pods, err := clientset.CoreV1().Pods(*namespace).List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("listing pods failed: %v", err)
    }
    log.Printf("Pods in %q:", *namespace)
    for _, p := range pods.Items {
        log.Println(" •", p.Name)
    }

    
    nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("listing nodes failed: %v", err)
    }
    log.Println("Nodes in cluster:")
    for _, n := range nodes.Items {
        log.Println(" •", n.Name)
    }

    
    deps, err := clientset.AppsV1().Deployments(*namespace).List(context.Background(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("listing deployments failed: %v", err)
    }
    log.Printf("Deployments in %q:", *namespace)
    for _, d := range deps.Items {
        log.Println(" •", d.Name)
    }
}
