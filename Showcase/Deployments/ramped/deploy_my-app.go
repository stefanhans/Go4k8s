package main

import (
	"encoding/json"
	"flag"
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

	// Read '-f <yamlfile>'
	var yamlFilename string
	flag.StringVar(&yamlFilename, "f", "", "Filename of YAML configuration")
	flag.Parse()

	// TODO: Prepare '-f <yamlfile>
	//yamlFilename = "app-v1.yaml"

	// *********************************
	fmt.Printf("Read %q\n", yamlFilename)

	// Get filename
	yamlFilepath, err := filepath.Abs(yamlFilename)

	// Get reader from file opening
	reader, err := os.Open(yamlFilepath)
	if err != nil {
		panic(err)
	}

	// *********************************
	fmt.Println("Prepare YAML to JSON decoding")

	// Split YAML into chunks or k8s resources, respectively
	yamlDecoder := yaml.NewDocumentDecoder(ioutil.NopCloser(reader))

	// Create decoding function used for YAML to JSON decoding
	decode := scheme.Codecs.UniversalDeserializer().Decode

	// *********************************
	fmt.Println("Decode deployment from YAML to JSON")

	// Read first resource - expecting deployment with size < 2048
	// TODO: handle size expectations programmatically
	yamlDeployment := make([]byte, 2048)
	_, err = yamlDecoder.Read(yamlDeployment)
	if err != nil {
		panic(err)
	}

	// Trim unnecessary trailing 0x0 signs which are not accepted
	TrimmedYamlDeployment := strings.TrimRight(string(yamlDeployment), string((byte(0))))
	//fmt.Println("\n####### TrimmedYamlDeployment #####\n\n\n", TrimmedYamlDeployment)

	// Decode deployment resource from YAML to JSON
	jsonDeployment, _, err := decode([]byte(TrimmedYamlDeployment), nil, nil)
	if err != nil {
		panic(err)
	}

	// Check "kind: deployment"
	if jsonDeployment.GetObjectKind().GroupVersionKind().Kind != "Deployment" {
		panic(fmt.Sprintf("\n####### N0 \"Kind: Deployment\" ######\n%#v\n", jsonDeployment))
	}

	// Marshall JSON Deployment
	d, err := json.Marshal(&jsonDeployment)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("\n###### Indent JSON without empty items #####\n%s\n\n", string(d))

	// *********************************
	fmt.Println("Define Deployment")

	// Unmarshall JSON into Deployment struct
	var deployment appsv1beta1.Deployment
	err = json.Unmarshal(d, &deployment)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("type Deployment struct: %#v\n", deployment)

	// *********************************
	fmt.Println("Decode service from YAML to JSON")

	// Read second resource - expecting service with size < 1024
	// TODO: handle size expectations programmatically
	yamlService := make([]byte, 1024)
	_, err = yamlDecoder.Read(yamlService)
	if err != nil {
		panic(err)
	}

	// Trim unnecessary trailing 0x0 signs which are not accepted
	TrimmedYamlService := strings.TrimRight(string(yamlService), string((byte(0))))
	//fmt.Println("\n####### TrimmedYamlService #####\n\n\n", TrimmedYamlService)

	// Decode service resource from YAML to JSON
	jsonService, _, err := decode([]byte(TrimmedYamlService), nil, nil)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("\n####### jsonService ######\n%#v\n", jsonService)

	// Check "kind: service"
	if jsonService.GetObjectKind().GroupVersionKind().Kind != "Service" {
		panic(fmt.Sprintf("\n####### N0 \"Kind: Service\" ######\n%#v\n", jsonService))
	}

	// Marshall JSON Service
	s, err := json.Marshal(&jsonService)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("\n###### Indent JSON without empty items #####\n%s\n\n", string(s))

	// *********************************
	fmt.Println("Define Service")

	// Unmarshall JSON into Service struct
	var service corev1.Service
	err = json.Unmarshal(s, &service)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("type Service struct: %#v\n", service)

	// *********************************
	fmt.Println("Build Config")

	config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")
	if err != nil {
		panic(err)
	}

	// *********************************
	fmt.Println("Create Clientset")

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// *********************************
	fmt.Println("Create Client for Deployments")

	deploymentsClient := clientset.AppsV1beta1().Deployments(corev1.NamespaceDefault)

	// *********************************
	fmt.Println("Create Client for Services")

	servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)

	// *********************************
	fmt.Println("Create Deployment")

	createdDeployment, err := deploymentsClient.Create(&deployment)
	if err != nil {
		fmt.Println(err)

		// *********************************
		fmt.Println("Update Deployment")

		updatedDeployment, err := deploymentsClient.Update(&deployment)
		if err != nil {
			panic(err)
		}
		fmt.Println("Deployment updated")
		fmt.Printf("Watch out for all Pods of %q running...\n", updatedDeployment.Name)
	} else {

		fmt.Println("Deployment created")
		fmt.Printf("Watch out for all Pods of %q running...\n", createdDeployment.Name)

		// *********************************
		fmt.Println("Create Service")

		createdService, err := servicesClient.Create(&service)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Watch out for the Service of %q running...\n", createdService.Name)
	}

	// *********************************
	fmt.Println("Get running service")

	// Get Pod "kube-addon-manager-minikube" of "kube-system" to retrieve 'minikube ip'
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}

	runningService, err := servicesClient.Get("my-app", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPlease verify: http://%s:%v\n\n", pod.Status.HostIP, runningService.Spec.Ports[0].NodePort)
}
