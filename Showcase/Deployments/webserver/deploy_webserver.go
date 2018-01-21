package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	// Check security config
	prompt("to check security config")
	if config.Insecure {
		fmt.Printf("Connecting to Config.Host %q is insecure!\n", config.Host)
	} else {
		fmt.Printf("Connecting to Config.Host %q is secure!\n", config.Host)
	}

	// Create Clientset
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
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	// Watch the client for deployments
	deploymentWatch, err := deploymentsClient.Watch(metav1.ListOptions{})
	watchCh := deploymentWatch.ResultChan()

	// Define Deployment
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "webserver-deployment",
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(2),
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
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	// Define Service
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

	// Create Deployment
	fmt.Println("Creating deployment...")
	resultDeployment, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}

	fmt.Println("######################################################")
	for evt := range watchCh {
		fmt.Printf("evt.Type: %v\n", evt.Type)
		if evt.Type == "MODIFIED" {
			pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
				LabelSelector: "app=webserver",
			})
			if err != nil {
				panic(err)
			}

			i := 2
			for _, pod := range pods.Items {
				if pod.Status.ContainerStatuses[0].Ready {
					fmt.Printf("Pod %q is running\n", pod.ObjectMeta.Name)

					i--
				}
			}
			if i == 0 {
				break
			}
		}
	}
	fmt.Println("######################################################")

	// Get list of pods
	podList, err := clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	// Get pods' description in json
	for n, pod := range podList.Items {
		b, err := json.MarshalIndent(&pod, fmt.Sprintf("%d:", n), "  ")
		if err != nil {
			fmt.Printf(" !!! Pod %d !!! : Error: %s", n, err)
			continue
		}
		fmt.Printf(" *** Pod %d *** \n%s\n\n", n, string(b))
	}

	// Get pods' customized description
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

		for _, condition := range pod.Status.Conditions {
			if condition.Status == "True" {
				fmt.Printf("Pod[%v].Status.Condition: True Status: %s\n", n, condition.Type)
			}
		}

		fmt.Printf("Pod[%v].Status.HostIP: %v\n", n, pod.Status.HostIP)
		fmt.Printf("Pod[%v].Status.HostIP: %v\n", n, pod.Status.PodIP)
		fmt.Printf("Pod[%v].Status.StartTime: %v\n", n, pod.Status.StartTime)

		for m, containerStatus := range pod.Status.ContainerStatuses {
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].Name: %v\n", n, m, containerStatus.Name)
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].Image: %v\n", n, m, containerStatus.Image)
			fmt.Printf("Pod[%v].Status.ContainerStatuses[%d].State: %v\n", n, m, containerStatus.State)
		}
		fmt.Println()
	}

	// Wait for running container - NOT RECOMMENDED - Use Watch() instead
	i := 0
	for i < 2 {
		podList, err = clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for n, pod := range podList.Items {
			if pod.Status.ContainerStatuses[0].Ready {
				i++
			}
			containerStatus := pod.Status.ContainerStatuses[0]
			if containerStatus.State.Waiting != nil {
				fmt.Printf("%d: %s: %s: %s\n", n, pod.ClusterName, pod.Spec.Containers[0].Name, containerStatus.State.Waiting.Reason)
			}
			if containerStatus.State.Running != nil {
				fmt.Printf("%d: %s: %s: Container is Running\n", n, pod.ClusterName, pod.Spec.Containers[0].Name)
			}
			if containerStatus.State.Terminated != nil {
				fmt.Printf("%d: %s: %s: %s\n", n, pod.ClusterName, pod.Spec.Containers[0].Name, containerStatus.State.Terminated.Reason)
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Printf("Created deployment %q.\n", resultDeployment.GetObjectMeta().GetName())

	// Create Service
	prompt("to create the service")
	fmt.Println("Creating service...")
	resultService, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created service %q.\n", resultService.GetObjectMeta().GetName())
	fmt.Printf("Please verify: http://%s:%v\n", pod.Status.HostIP, resultService.Spec.Ports[0].NodePort)

	fmt.Printf("pod.Status.Message: %+v\n", pod.Status.Message)
	fmt.Printf("pod.Status.Reason: %+v\n", pod.Status.Reason)

	// Delete Deployment
	prompt("to delete the deployment")
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("webserver-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}

	// Wait for deleted container - NOT RECOMMENDED - Use Watch() instead
	i = 0
	for i < 2 {
		podList, err = clientset.CoreV1().Pods(apiv1.NamespaceDefault).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		for n, pod := range podList.Items {
			if pod.Status.ContainerStatuses[0].Ready {
				i++
			}
			containerStatus := pod.Status.ContainerStatuses[0]
			if containerStatus.State.Waiting != nil {
				fmt.Printf("%d: %s: %s: %s\n", n, pod.ClusterName, pod.Spec.Containers[0].Name, containerStatus.State.Waiting.Reason)
			}
			if containerStatus.State.Running != nil {
				fmt.Printf("%d: %s: %s: Container is Running\n", n, pod.ClusterName, pod.Spec.Containers[0].Name)
			}
			if containerStatus.State.Terminated != nil {
				fmt.Printf("%d: %s: %s: %s\n", n, pod.ClusterName, pod.Spec.Containers[0].Name, containerStatus.State.Terminated.Reason)
			}
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Deleted deployment.")

	// Delete Deployment
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
	fmt.Printf("-> Press Return key %s", str)
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
