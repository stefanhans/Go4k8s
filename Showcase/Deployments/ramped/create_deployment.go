package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	//
	filename, err := filepath.Abs("./test.yaml")

	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	yamlDecoder := yaml.NewDocumentDecoder(ioutil.NopCloser(reader))

	yamlDeployment := make([]byte, 2048)
	_, err = yamlDecoder.Read(yamlDeployment)
	if err != nil {
		panic(err)
	}
	TrimmedYamlDeployment := strings.TrimRight(string(yamlDeployment), string((byte(0))))
	fmt.Println("\n####### YAML #####\n\n\n", TrimmedYamlDeployment)

	yamlService := make([]byte, 2048)
	_, err = yamlDecoder.Read(yamlService)
	if err != nil {
		panic(err)
	}
	TrimmedYamlService := strings.TrimRight(string(yamlService), string((byte(0))))
	fmt.Println("\n####### YAML #####\n\n\n", TrimmedYamlService)

	decode := scheme.Codecs.UniversalDeserializer().Decode

	jsonDeployment, _, err := decode([]byte(TrimmedYamlDeployment), nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n####### Complete JSON without indent ######\n%#v\n", jsonDeployment)

	d, err := json.MarshalIndent(&jsonDeployment, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n###### Indent JSON without empty items #####\n%s\n\n", string(d))

	jsonService, _, err := decode([]byte(TrimmedYamlService), nil, nil)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\n####### Complete JSON without indent ######\n%#v\n", jsonService)

	s, err := json.MarshalIndent(&jsonService, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n###### Indent JSON without empty items #####\n%s\n\n", string(s))

	var deployment appsv1beta1.Deployment

	err = json.Unmarshal(d, &deployment)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", deployment)

	var service corev1.Service
	err = json.Unmarshal(s, &service)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", service)

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
	deploymentsClient := clientset.AppsV1beta1().Deployments(corev1.NamespaceDefault)

	// Create Client for Services
	fmt.Println("Create Client for Services")
	servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)

	// Define Deployment
	fmt.Println("Define Deployment")

	// Create Deployment
	fmt.Println("Create Deployment")
	resultDeployment, err := deploymentsClient.Create(&deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Watch out for all Pods of %q running...\n", resultDeployment.Name)

	// Create Service
	fmt.Println("Create Service")
	resultService, err := servicesClient.Create(&service)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Watch out for all Pods of %q running...\n", resultService.Name)

	// Get Pod "kube-addon-manager-minikube" of "kube-system" to retrieve 'minikube ip'
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPlease verify: http://%s:%v\n", pod.Status.HostIP, resultService.Spec.Ports[0].NodePort)

}
