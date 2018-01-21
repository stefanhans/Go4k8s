Go: Build and Run

    go build -o ./go-webserver .
    ./go-webserver
    ^C

Go: Test

    curl --head http://127.0.0.1:8080

Docker: Build and Run

    docker build -t stefanhans/go-webserver .
    docker run --publish 8080:8080 --name test --rm stefanhans/go-webserver:latest

Docker: Test

    curl --head http://127.0.0.1:8080
    docker stop test
    
Docker: Push to Docker Hub

    docker push stefanhans/go-webserver:latest

Kubernetes: Deploy, Expose and Test

    kubectl create -f DeployGoWebserver.yaml
    kubectl get pods,service -l app=go-webserver # wait until all is up and running
    curl --head http://$(minikube ip):$(kubectl get svc -l app=go-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    kubectl delete all -l app=go-webserver

