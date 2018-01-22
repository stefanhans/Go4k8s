Prerequisites: Tested [webserver image](../../Images/webserver), create two versions (*.yaml), and push to GitHub.

Having a running environment, e.g. `minikube start`

Kubernetes Staging: Deploy and Test

    kubectl create -f DeployStagingWebserver.yaml
    
    kubectl get pods,service -l app=webserver,env=staging
    kubectl logs -l app=webserver,env=staging
    
    curl --head http://$(minikube ip):$(kubectl get svc -l app=webserver,env=staging -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=webserver,env=staging -o jsonpath='{.items[0].spec.ports[0].nodePort}')"

Kubernetes Staging: Cleanup

    kubectl delete all -l app=webserver,env=staging

Kubernetes Production: Deploy and Test

    kubectl create -f DeployProdWebserver.yaml
    
    kubectl get pods,service -l app=webserver,env=production
    kubectl logs -l app=webserver,env=production
    
    curl --head http://$(minikube ip):$(kubectl get svc -l app=webserver,env=production -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=webserver,env=production -o jsonpath='{.items[0].spec.ports[0].nodePort}')"

Docker: Build Image and Container

    docker build -t stefanhans/dev-ops-webserver .
    docker create --publish 8080:8080 --name dev-ops-webserver-container stefanhans/dev-ops-webserver
    
Docker: Run Container and Verify

    docker start dev-ops-webserver-container
    curl --head http://localhost:8080
    
Docker: Stop Container

    docker stop dev-ops-webserver-container
    
Docker: Push Container to Docker Hub

    docker push stefanhans/dev-ops-webserver

Kubernetes: Deploy and Test

    kubectl create -f DeployProdWebserver.yaml
    
    kubectl get pods,service -l app=prod-webserver
    kubectl logs -l app=prod-webserver
    
    curl --head http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')"
    
Kubernetes: Cleanup

Do not cleanup before 'Next Step'!
    
    kubectl delete all -l app=prod-webserver
    
Next Step: Try out the [DevOps deployment programmed in Go](../../Deployments/dev-ops)
    