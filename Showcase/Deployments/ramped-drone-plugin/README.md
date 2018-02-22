### Work In Progress

<a href="https://asciinema.org/a/8C4FwMI74WkbPNaIeo4MUZHgi" target="_blank"><img src="https://asciinema.org/a/8C4FwMI74WkbPNaIeo4MUZHgi.png" /></a>

This scenario is using or can use, respectively, the following images from docker hub:

stefanhans/webserver:1.0.0 (green)

stefanhans/webserver:1.0.1 (blue)

stefanhans/webserver:1.0.2 (yellow)

stefanhans/webserver:1.0.3 (red)

stefanhans/webserver:1.0.4 (gray)

All images are simply presenting one webpage and the versions change mainly the background color.

### Preparations:

- Prepare a cluster environment, e.g. [minikube](https://github.com/kubernetes/minikube)

- Copy all files according to the Dockerfile, i.e. config, ca.crt, client.crt, client.key

-

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


Using [test-webserver](https://github.com/stefanhans/Go4k8s/tree/master/Showcase/Images/test-webserver) as template
