package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes/scheme"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/client-go/kubernetes/scheme"
	//"io"
	"os"
	//"k8s.io/kubernetes/pkg/util/strings"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	//"k8s.io/client-go/util/retry"
)

func main() {

	filename, err := filepath.Abs("./test.yaml")

	reader, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	yamlDecoder := yaml.NewDocumentDecoder(ioutil.NopCloser(reader))
	d := make([]byte, 2048)
	_, err = yamlDecoder.Read(d)
	if err != nil {
		panic(err)
	}
	deployYaml := strings.TrimRight(string(d), string((byte(0))))
	fmt.Println("\n####### YAML #####\n\n\n", deployYaml)



	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(deployYaml), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	fmt.Printf("\n####### Complete JSON without indent ######\n%#v\n", obj)

	b, err := json.MarshalIndent(&obj, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("\n###### Indent JSON without empty items #####\n%s\n\n", string(b))


	var deployment appsv1beta1.Deployment
	//var service corev1.Service

	err = json.Unmarshal(b, &deployment)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Value: %#v\n", deployment)



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

	// Define Deployment
	fmt.Println("Define Deployment")

	// Create Deployment
	fmt.Println("Create Deployment")
	resultDeployment, err := deploymentsClient.Create(&deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Watch out for all Pods of %q running...\n", resultDeployment.Name)
}
