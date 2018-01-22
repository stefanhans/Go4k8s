Go: Build and Run

    go build -o ./webserver .
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

Kubernetes: Deploy and Test

Having a running environment, e.g. `minikube start`

    kubectl create -f DeployGoWebserver.yaml
    
    kubectl get pods,service -l app=webserver
    kubectl logs -l app=webserver
    
    curl --head http://$(minikube ip):$(kubectl get svc -l app=webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')"
    
Kubernetes: Cleanup
    
    kubectl delete all -l app=webserver
    
Next Step: Try out the [deployment programmed in Go](../../Deployments/webserver)

