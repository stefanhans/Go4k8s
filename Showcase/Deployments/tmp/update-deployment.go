package tmp

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

	// TODO: Prepare '-f <yamlfile>
	yamlFilename := "app-v2.yaml"


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
	fmt.Printf("type Deployment struct: %#v\n", deployment)


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

	/*
	// Get deployment name
	runningDeployment, err := deploymentsClient.Get("my-app", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	*/




	/*
	if jsonDeployment.GetObjectKind().GroupVersionKind().Kind != "Deployment" {
		panic(fmt.Sprintf("\n####### N0 \"Kind: Deployment\" ######\n%#v\n", jsonDeployment))
	}
	*/




	// *********************************
	fmt.Println("Update Deployment")

	updatedDeployment, err := deploymentsClient.Update(&deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Watch out for all Pods of %q running...\n", updatedDeployment.Name)

	// Get Pod "kube-addon-manager-minikube" of "kube-system" to retrieve 'minikube ip'
	pod, err := clientset.CoreV1().Pods("kube-system").Get("kube-addon-manager-minikube", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}


	// *********************************
	fmt.Println("Get running service")

	servicesClient := clientset.CoreV1().Services(corev1.NamespaceDefault)
	createdService, err := servicesClient.Get("my-app", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nPlease verify: http://%s:%v\n\n", pod.Status.HostIP, createdService.Spec.Ports[0].NodePort)

}
