Prerequisites: Tested [webserver image](../../Images/webserver), and pushed to GitHub.


Docker: Build Image and Container

    docker build -t stefanhans/dev-ops-webserver .
    docker create --publish 8080:8080 --name dev-ops-webserver-container stefanhans/dev-ops-webserver
    
Docker: Run Container and Verify

    docker start dev-ops-webserver-container
    http://localhost:8080
    
Docker: Stop Container

    docker stop dev-ops-webserver-container
    
Docker: Push Container to Docker Hub

    docker push stefanhans/dev-ops-webserver

Kubernetes: Deploy and Test

Having a running environment, e.g. `minikube start`

    kubectl create -f DeployProdWebserver.yaml
    
    kubectl get pods,service -l app=prod-webserver
    kubectl logs -l app=prod-webserver
    
    curl --head http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')"
    
Kubernetes: Cleanup

Do not cleanup before 'Next Step'!
    
    kubectl delete all -l app=prod-webserver
    
Next Step: Try out the [DevOps deployment programmed in Go](../../Deployments/dev-ops)
    