package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {

	// Read '-f <yamlfile>'
	var yamlFilename string
	flag.StringVar(&yamlFilename, "f", "deployment.yaml", "Filename of YAML configuration")

	// Read name of deployment
	var deploymentName string
	flag.StringVar(&deploymentName, "d", "deployment", "Name of the deployment")

	// Parse commandline parameter
	flag.Parse()

	// *********************************
	fmt.Printf("Read %q\n", yamlFilename)

	// Get filename
	yamlFilepath, err := filepath.Abs(yamlFilename)

	// Get reader from file opening
	reader, err := os.Open(yamlFilepath)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q\n", string(data))

	re := regexp.MustCompile(`\{\{(\..*?)\}\}`) // ` is a backtick
	variables := re.FindAllStringSubmatch(string(data), 64)

	outstring := "apiVersion: experimental\n" +
		"kind: DynamicYaml\n" +
		"name: my-app\n"

	for _, variable := range variables {
		//fmt.Println(variable[1])
		outstring = outstring + fmt.Sprintf("%s: \n", strings.TrimLeft(strings.ToLower(variable[1]), "."))
	}


	updateYamlFilename := fmt.Sprintf("update_%s.yaml", deploymentName)

	fmt.Println(updateYamlFilename)
	fmt.Println(outstring)


	err = ioutil.WriteFile(updateYamlFilename, []byte(outstring), 0666)
	if err != nil {
		panic(err)
	}
}
