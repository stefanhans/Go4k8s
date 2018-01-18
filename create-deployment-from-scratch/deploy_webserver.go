package main

import (
	"bufio"
	//"flag"
	"fmt"
	"os"
	//"path/filepath"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/homedir"
	//"k8s.io/client-go/util/retry"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	//"k8s.io/apimachinery/pkg/util/intstr"
	//"k8s.io/client-go/kubernetes/typed/core/v1"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	fmt.Printf("config.Host: %s\n", config.Host)

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "go-webserver-deployment",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(2),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "go-webserver",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "go-webserver",
							Image: "stefanhans/go-webserver",
						},
					},
				},
			},
		},
	}

	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "go-webserver-service",
			Labels: map[string]string{
				"app": "go-webserver",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "go-webserver",
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: "TCP",
					Port:     int32(8080),
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	resultDeployment, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", resultDeployment.GetObjectMeta().GetName())

	// Create Service
	prompt()
	fmt.Println("Creating service...")
	resultService, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", resultService.GetObjectMeta().GetName())

	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("go-webserver-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")

	// Delete Deployment
	prompt()
	fmt.Println("Deleting service...")
	deletePolicy = metav1.DeletePropagationForeground
	if err := servicesClient.Delete("go-webserver-service", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service.")
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
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
