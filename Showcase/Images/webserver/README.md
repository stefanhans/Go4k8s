Go: Build and Run

    go build -o ./webserver .
    chmod +x webserver
    ./webserver
    # Stop later via Ctrl-C

Go: Test

    curl --head http://127.0.0.1:8080

Docker: Build and Run

    docker build -t stefanhans/webserver .
    docker run --publish 8080:8080 --name test --rm stefanhans/webserver:latest

Docker: Test

    curl --head http://127.0.0.1:8080
    docker stop test
    
Docker: Push to Docker Hub

    docker push stefanhans/webserver:latest

Kubernetes: Deploy, Expose and Test

    kubectl create -f DeployGoWebserver.yaml
    kubectl get pods,service -l app=webserver # wait until all is up and running
    curl --head http://$(minikube ip):$(kubectl get svc -l app=webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    kubectl delete all -l app=webserver

