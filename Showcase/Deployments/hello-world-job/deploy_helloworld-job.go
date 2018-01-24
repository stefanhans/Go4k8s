package main

import (
	"bufio"
	"fmt"
	"os"

	apibatchv1 "k8s.io/api/batch/v1"
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

	// Create Client for Jobs
	jobsClient := clientset.BatchV1().Jobs(apiv1.NamespaceDefault)


	// Define Job
	job := &apibatchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: "hello-world-job",
		},
		Spec: apibatchv1.JobSpec{
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
					RestartPolicy: apiv1.RestartPolicyNever,
				},
			},
			BackoffLimit: int32Ptr(4),
		},
	}

	// Create Job
	fmt.Printf("\nCreating job...\n")
	result, err := jobsClient.Create(job)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nCreated job %q.\n", result.GetObjectMeta().GetName())

	// Delete Job
	prompt()
	fmt.Println("Deleting job...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := jobsClient.Delete("hello-world-job", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted job.")
}

func prompt() {
	fmt.Printf("\n-> Press Return key to stop the job.\n")
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
