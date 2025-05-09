package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func int32Ptr(i int32) *int32 {
	return &i
}

func createDeployment(clientset *kubernetes.Clientset) {
	var name, image string
	var replicas int32

	fmt.Print("Enter deployment name: ")
	fmt.Scan(&name)
	fmt.Print("Enter container image (e.g., nginx:1.12): ")
	fmt.Scan(&image)
	fmt.Print("Enter number of replicas: ")
	fmt.Scan(&replicas)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Printf("Failed to create deployment: %v\n", err)
		return
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
}

func deleteDeployment(clientset *kubernetes.Clientset) {
	var name string
	fmt.Print("Enter the name of the deployment to delete: ")
	fmt.Scan(&name)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete(context.TODO(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		fmt.Printf("Failed to delete deployment: %v\n", err)
		return
	}
	fmt.Println("Deleted deployment.")
}

func updateDeployment(clientset *kubernetes.Clientset) {
	var name, image string
	var replicas int32

	fmt.Print("Enter the name of the deployment to update: ")
	fmt.Scan(&name)
	fmt.Print("Enter new container image (e.g., nginx:1.13): ")
	fmt.Scan(&image)
	fmt.Print("Enter new number of replicas: ")
	fmt.Scan(&replicas)

	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	fmt.Println("Updating deployment...")
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get(context.TODO(), name, metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("failed to get latest version of deployment: %v", getErr)
		}

		result.Spec.Replicas = int32Ptr(replicas)
		result.Spec.Template.Spec.Containers[0].Image = image
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		fmt.Printf("Update failed: %v\n", retryErr)
		return
	}
	fmt.Println("Deployment has been updated.")
}

func listDeployments(clientset *kubernetes.Clientset) {
	deployments, err := clientset.AppsV1().Deployments(apiv1.NamespaceDefault).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Failed to list deployments: %v\n", err)
		return
	}

	fmt.Println("Listing deployments in default namespace:")
	for _, deployment := range deployments.Items {
		fmt.Println(deployment.Name)
	}
}

func handleDeploymentActions(clientset *kubernetes.Clientset) {
	for {
		fmt.Println("\nChoose an action:")
		fmt.Println("1) List deployments")
		fmt.Println("2) Create deployment")
		fmt.Println("3) Update deployment")
		fmt.Println("4) Delete deployment")
		fmt.Println("5) Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			listDeployments(clientset)
		case 2:
			createDeployment(clientset)
		case 3:
			updateDeployment(clientset)
		case 4:
			deleteDeployment(clientset)
		case 5:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 5.")
		}
	}
}

func main() {
	kubeconfig := flag.String("kubeconfig", "", "Absolute file path of kubeconfig file")
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(fmt.Errorf("failed to build kubeconfig: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create Kubernetes client: %v", err))
	}

	handleDeploymentActions(clientset)
}
