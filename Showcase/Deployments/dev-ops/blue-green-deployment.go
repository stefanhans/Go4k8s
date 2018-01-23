package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"time"
)

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter version for new deployment: ")
	version, _ := reader.ReadString('\n')

	// Number of replicas
	replicas := int32(2)
	fmt.Printf("Number of replica configured: %d", replicas)

	fmt.Println()

	// Build Config
	fmt.Println("Build Config")
	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	// Create Clientset
	fmt.Println("Create Clientset")
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Println()

	// Create Client for Deployments
	fmt.Println("Create Client for Deployments")
	deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

	// Define Deployment
	fmt.Println("Define Deployment")
	deployment := &appsv1beta1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "webserver-staging-deployment",
			Labels: map[string]string{
				"app": "webserver",
				"env": "staging",
			},
		},
		Spec: appsv1beta1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "webserver",
						"env": "staging",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "webserver",
							Image: fmt.Sprint("stefanhans/webserver:", strings.TrimSpace(version)),
						},
					},
				},
			},
		},
	}

	fmt.Println()

	// Create Client for Services
	fmt.Println("Create Client for Services")
	servicesClient := clientset.CoreV1().Services(apiv1.NamespaceDefault)

	// Define Service
	fmt.Println("Define Service")
	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: "staging-webserver-service",
			Labels: map[string]string{
				"app": "webserver",
				"env": "staging",
			},
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": "webserver",
				"env": "staging",
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: "TCP",
					Port:     int32(8080),
					NodePort: int32(30002), // The range of valid ports is 30000-32767
				},
			},
			Type: apiv1.ServiceTypeLoadBalancer,
		},
	}

	fmt.Println()

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

	time.Sleep(time.Second)

	// Watch out for all Pods running...
	fmt.Printf("Watch out for all Pods of %q running...\n", resultDeployment.Name)
	for evt := range watchCh {
		if evt.Type == "MODIFIED" {
			pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{
				LabelSelector: "env=staging",
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

	// Create Service
	prompt("to create the service")
	fmt.Println("Create Service")
	stagingService, err := servicesClient.Create(service)
	if err != nil {
		panic(err)
	}

	// Deployment and services finished
	fmt.Printf("All Pods of %q are running\n", resultDeployment.Name)
	fmt.Printf("Created service %q.\n", stagingService.Name)

	// Get Pod "kube-addon-manager-minikube" of "kube-system" to retrieve 'minikube ip'
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nPlease verify staging: http://%s:%v\n", pod.Status.HostIP, stagingService.Spec.Ports[0].NodePort)

	// Switch production loadbalancer to staging deployment
	prompt("to switch production loadbalancer to staging deployment")
	fmt.Println("Switch production loadbalancer to staging deployment...")

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := servicesClient.Get("webserver-service", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get service 'webserver-service': %v", getErr))
		}
		result.Spec.Selector = map[string]string{
			"app": "webserver",
			"env": "staging",
		}
		_, updateErr := servicesClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Switched production loadbalancer to staging deployment")

	productionService, err := servicesClient.Get("webserver-service", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPlease verify production: http://%s:%v\n", pod.Status.HostIP, productionService.Spec.Ports[0].NodePort)

	// Switch production image to new version
	prompt("to switch production image to new version")
	fmt.Println("Switch production image to new version...")

	retryErr = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get("webserver-prod-deployment", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get deployment 'webserver-prod-deployment': %v", getErr))
		}
		result.Spec.Template.Spec.Containers[0].Image = fmt.Sprint("stefanhans/webserver:", strings.TrimSpace(version))
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Printf("Switched production image to new version %q\n", strings.TrimSuffix(version, "\n"))


	// Switch production loadbalancer back to production deployment
	prompt("to switch production loadbalancer back to production deployment")
	fmt.Println("Switch production loadbalancer back to production deployment...")

	retryErr = retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := servicesClient.Get("webserver-service", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get service 'webserver-service': %v", getErr))
		}
		result.Spec.Selector = map[string]string{
			"app": "webserver",
			"env": "production",
		}
		_, updateErr := servicesClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Switched production loadbalancer back to production deployment")

	fmt.Printf("\nPlease verify production: http://%s:%v\n", pod.Status.HostIP, productionService.Spec.Ports[0].NodePort)

	// Delete staging deployment and service
	prompt("to delete staging deployment and service")
	fmt.Println("Delete staging deployment and service...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete("webserver-staging-deployment", &metav1.DeleteOptions{
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
	fmt.Println("Delete service...")
	deletePolicy = metav1.DeletePropagationForeground
	if err := servicesClient.Delete("staging-webserver-service", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted")
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
