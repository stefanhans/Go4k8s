### Work In Progress

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/Go4k8s/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped?status.svg)](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)](https://goreportcard.com/report/github.com/Go4k8s/tree/master/Showcase/Deployments/ramped)

Using [test-webserver](https://github.com/stefanhans/Go4k8s/tree/master/Showcase/Images/test-webserver) as template


<a href="https://asciinema.org/a/8C4FwMI74WkbPNaIeo4MUZHgi" target="_blank"><img src="https://asciinema.org/a/8C4FwMI74WkbPNaIeo4MUZHgi.png" /></a>

### Preparations:

    cat Dockerfile

Copy your files to

    "k8s.io/client-go/rest"
    config, err := rest.InClusterConfig()

instead of

    "k8s.io/client-go/tools/clientcmd"
    config, err := clientcmd.BuildConfigFromFlags("", "/home/stefan/.kube/config")

build app

    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./update
GOOS=linux go build -o ./update
edit Dockerfile

build image

    docker build -t stefanhans/ramped-test .

push image

    docker push stefanhans/ramped-test

test run

    docker run --rm -e PLUGIN_VERSION=1.0.2 stefanhans/ramped-update
    docker run --rm stefanhans/ramped-test


create Dockerfile

    cat >Dockerfile <<EOF
    FROM debian
    COPY ["./app", "./deployment.yaml", "update.yaml", "/"]
    ENTRYPOINT /app
    EOF

build image

    docker build -t update-in-cluster .

push image as needed

create wrapper for deleting inactive job

    cat >update.bash <<EOF
    kubectl create -f update_job.yaml

    while [ "$(kubectl get jobs update-job -o jsonpath='{.status.active}')" != "" ]
    do
        sleep 1
        printf "."
    done

    kubectl delete -f update_job.yaml
    EOF

    chmod +x update.bash

run wrapper

    ./update.bash