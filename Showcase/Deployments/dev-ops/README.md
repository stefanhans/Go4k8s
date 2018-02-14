[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/Go4k8s/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped?status.svg)](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)](https://goreportcard.com/report/github.com/Go4k8s/tree/master/Showcase/Deployments/ramped)

Prerequisites: Tested [test-webserver image](../../Images/test-webserver).

Build and run Go executables

    go build -o ./blue-green-deployment blue-green-deployment.go 
    ./blue-green-deployment
    
    
View service labels

    kubectl get service -l app=webserver -o jsonpath='{.items[0].metadata.name}{": env="}{.items[0].spec.selector.env}{"\n"}{.items[1].metadata.name}{": env="}{.items[1].spec.selector.env}{"\n"}'
    
    
View container images

    kubectl get pods -l env=production -o jsonpath='{.items[*].spec.containers[0].image}{"\n"}'
    
    
View rollout history of production deployment
    
    kubectl rollout history deploy/webserver-prod-deployment

