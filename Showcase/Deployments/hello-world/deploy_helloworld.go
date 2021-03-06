package main

import (
	"bufio"
	"fmt"
	"os"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// Build Config
	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	// Create Clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Create Client for Deployments
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	// Define Deployment
	deploymentsDefinition := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello-world-deployment",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "hello-world",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "hello-world",
							Image: "stefanhans/hello-world",
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Printf("\nCreating deployment...\n")
	deployment, err := deploymentsClient.Create(deploymentsDefinition)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nCreated deployment %q.\n", deployment.GetObjectMeta().GetName())

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("hello-world-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

func prompt() {
	fmt.Printf("\n-> Press Return key to stop the deployment.\n")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func int32Ptr(i int32) *int32 { return &i }
