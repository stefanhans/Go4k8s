Prerequisites: Tested [webserver image](../../Images/webserver), created two versions (*.yaml), and pushed to GitHub.

Having a running environment, e.g. `minikube start`

### Do a Blue-Green Deployment

Kubernetes Staging: Deploy and Test

    kubectl create -f DeployStagingWebserver.yaml
    
    kubectl get pods,service -l app=webserver,env=staging
    kubectl logs -l app=webserver,env=staging
    
    curl http://$(minikube ip):$(kubectl get svc -l app=webserver,env=staging -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=webserver,env=staging -o jsonpath='{.items[0].spec.ports[0].nodePort}')"

Kubernetes Production: Deploy and Test

    kubectl create -f DeployProdWebserver.yaml
    
    kubectl get pods,service -l app=webserver,env=production
    kubectl logs -l app=webserver,env=production
    
    curl http://$(minikube ip):$(kubectl get svc -l app=webserver,env=production -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=webserver,env=production -o jsonpath='{.items[0].spec.ports[0].nodePort}')"
    
Kubernetes: Switch Production Loadbalancer to Staging Deployment

    kubectl edit -f DeployProdWebserver.yaml
    
    Service: .spec.selector.env: production     =>    staging
    
Verify!    

Kubernetes: Switch Production Image to New Version 
  
    kubectl edit -f DeployProdWebserver.yaml
    
    Deployment: .spec.template.spec.containers.image: stefanhans/webserver:1.0.0     =>    stefanhans/webserver:1.0.1
    
Kubernetes: Switch Staging Loadbalancer to Production Deployment 

    kubectl edit -f DeployStagingWebserver.yaml
    
    Service: .spec.selector.env: staging     =>    production
    
Verify!   

Kubernetes: Switch Production Loadbalancer back to Production Deployment

    kubectl edit -f DeployProdWebserver.yaml
    
    Service: .spec.selector.env: staging     =>    production
    
Kubernetes Staging: Cleanup

    kubectl delete all -l app=webserver,env=staging
  
Leave production running for next step!

### Prepare Automated Blue-Green Deployment  

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
    
    curl http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')
    echo "http://$(minikube ip):$(kubectl get svc -l app=prod-webserver -o jsonpath='{.items[0].spec.ports[0].nodePort}')"
    
Kubernetes: Cleanup

Do not cleanup before 'Next Step'!
    
    kubectl delete all -l app=prod-webserver
    
Next Step: Try out the [DevOps deployment programmed in Go](../../Deployments/dev-ops)
    