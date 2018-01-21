Go: Build and Run

    go build -o ./hello-world .
    ./hello-world

Docker: Build and Run

    docker build -t stefanhans/hello-world .
    docker run --rm stefanhans/hello-world:latest

Docker: Push to Docker Hub

    docker push stefanhans/hello-world:latest

Kubernetes: Deploy, Expose and Test

    kubectl create -f DeployHelloWorld.yaml
    kubectl get pods -l app=hello-world # wait until all is up and running
    kubectl logs
    
Investigate using the last two steps repeatedly!
    
Kubernetes: Cleanup

    kubectl delete all -l app=hello-world

