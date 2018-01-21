Go: Build and Run

    go build -o ./hello-world .
    chmod +x hello-world
    ./hello-world

Docker: Build and Run

    docker build -t stefanhans/hello-world .
    docker run --rm stefanhans/hello-world:latest

Docker: Push to Docker Hub

    docker push stefanhans/hello-world:latest

Kubernetes: Deploy and Test

Having a running environment, e.g. `minikube start`

    kubectl create -f DeployHelloWorld.yaml
    
    kubectl get pods -l app=hello-world 
    kubectl logs -l app=hello-world
    
Investigate using the last two steps repeatedly!
    
Kubernetes: Cleanup

    kubectl delete all -l app=hello-world
    
Next Step: Try out the [deployment programmed in Go](../../Deployments/hello-world)

