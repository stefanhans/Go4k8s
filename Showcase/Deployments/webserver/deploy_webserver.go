package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	// Number of replicas
	replicas := int32(2)

	fmt.Println()

	// Build Config
	fmt.Println("Build Config")
	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	// Check security config
	if config.Insecure {
		fmt.Printf("Config.Host %q is insecure!\n", config.Host)
	} else {
		fmt.Printf("Config.Host %q is secure!\n", config.Host)
	}

	// Create Clientset
	fmt.Println("Create Clientset")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// !!!
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// Create Client for Deployments
	fmt.Println("Create Client for Deployments")
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	// Define Deployment
	fmt.Println("Define Deployment")
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "webserver-deployment",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "webserver",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "webserver",
							Image: "stefanhans/webserver",
						},
					},
				},
			},
		},
	}

	// Create Client for Services
	fmt.Println("Create Client for Services")
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	// Define Service
	fmt.Println("Define Service")
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "webserver-service",
			Labels: map[string]string{
				"app": "webserver",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "webserver",
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

	// Start watching the client for deployments
	fmt.Println("Start watching the client for deployments")
	deploymentWatch, err := deploymentsClient.Watch(metav1.ListOptions{})
	watchCh := deploymentWatch.ResultChan()

	// Create Deployment
	prompt("to create the deployment")
	fmt.Println("Create Deployment")
	resultDeployment, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}

	// Watch out for all Pods running...
	fmt.Printf("Watch out for all Pods of %q running...\n", resultDeployment.Name)
	for evt := range watchCh {
		if evt.Type == "MODIFIED" {
			pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
				LabelSelector: "app=webserver",
			})
			if err != nil {
				panic(err)
			}

			i := replicas
			for _, pod := range pods.Items {
				if pod.Status.ContainerStatuses[0].Ready {
					fmt.Printf("Pod %q is running\n", pod.ObjectMeta.Name)

					i--
				}
			}
			if i == int32(0) {
				break
			}
		}
	}
	fmt.Printf("All Pods of %q are running\n", resultDeployment.Name)

	// Get pods' description in json
	prompt("to get pods' description in json")

	// Get list of pods
	fmt.Println("Get list of pods")
	podList, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Get pods' description in json")
	for n, pod := range podList.Items {
		b, err := json.MarshalIndent(&pod, fmt.Sprintf("%d:", n), "  ")
		if err != nil {
			fmt.Printf(" !!! Pod %d !!! : Error: %s", n, err)
			continue
		}
		fmt.Printf(" *** Pod %d *** \n%s\n\n", n, string(b))
	}

	// Get pods' customized description
	prompt("to get pods' customized description")
	fmt.Println("Get pods' customized description")
	for n, pod := range podList.Items {
		fmt.Println()
		fmt.Printf(" *** Metadata *** \n")
		fmt.Printf("Pod[%v].Name: %s\n", n, pod.ObjectMeta.Name)
		fmt.Printf("Pod[%v].Namespace: %s\n", n, pod.ObjectMeta.Namespace)
		fmt.Printf("Pod[%v].CreationTimestamp: %s\n", n, pod.ObjectMeta.CreationTimestamp)
		fmt.Printf("Pod[%v].Labels[\"app\"]: %s\n", n, pod.ObjectMeta.Labels["app"])

		fmt.Println()
		fmt.Printf(" *** Spec *** \n")
		fmt.Printf("Pod[%v].Spec.Containers[0].Name: %s\n", n, pod.Spec.Containers[0].Name)
		fmt.Printf("Pod[%v].Spec.Containers[0].Image: %s\n", n, pod.Spec.Containers[0].Image)
		fmt.Printf("Pod[%v].Spec.NodeName: %s\n", n, pod.Spec.NodeName)

		fmt.Println()
		fmt.Printf(" *** Status *** \n")
		fmt.Printf("Pod[%v].Status.Phase: %v\n", n, pod.Status.Phase)
		fmt.Printf("Pod[%v].Status.StartTime: %v\n", n, pod.Status.StartTime)
		fmt.Printf("Pod[%v].Status.HostIP: %v\n", n, pod.Status.HostIP)
		fmt.Printf("Pod[%v].Status.PodIP: %v\n", n, pod.Status.PodIP)

		for _, condition := range pod.Status.Conditions {
			if condition.Status == "True" {
				fmt.Printf("Pod[%v].Status.Condition: True Status: %s\n", n, condition.Type)
			}
		}

		for m, containerStatus := range pod.Status.ContainerStatuses {
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].Name: %v\n", n, m, containerStatus.Name)
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].Image: %v\n", n, m, containerStatus.Image)
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].State: %v\n", n, m, containerStatus.State)
		}
		fmt.Println()
	}

	// Create Service
	prompt("to create the service")
	fmt.Println("Create Service")
	resultService, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}

	// Deployment and services finished
	fmt.Printf("All Pods of %q are running\n", resultDeployment.Name)
	fmt.Printf("Created service %q.\n", resultService.Name)

	fmt.Printf("\nPlease verify: http://%s:%v\n", pod.Status.HostIP, resultService.Spec.Ports[0].NodePort)

	// Delete Deployment
	prompt("to delete the deployment")
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("webserver-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}

	// Watch out for the deployment deleted...
	fmt.Printf("Watch out for the deployment %q deleted...\n", resultDeployment.Name)
	for evt := range watchCh {
		if evt.Type == "DELETED" {
			break
		}
	}
	fmt.Printf("The deployment of %q is deleted\n", resultDeployment.Name)

	// Delete Service
	prompt("to delete the service")
	fmt.Println("Deleting service...")
	deletePolicy = metav1.DeletePropagationForeground
	if err := servicesClient.Delete("webserver-service", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted service.")
}

func prompt(str string) {
	fmt.Printf("\n-> Press Return key %s", str)
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
