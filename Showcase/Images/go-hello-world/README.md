Go: Build and Run

    go build -o ./go-hello-world .
    ./go-hello-world

Docker: Build and Run

    docker build -t stefanhans/go-hello-world .
    docker run --rm stefanhans/go-hello-world:latest

Docker: Push to Docker Hub

    docker push stefanhans/go-hello-world:latest

Kubernetes: Deploy, Expose and Test

    kubectl create -f DeployGoHelloWorld.yaml
    kubectl get pods -l app=go-hello-world # wait until all is up and running
    kubectl logs
    kubectl delete all -l app=go-hello-world

