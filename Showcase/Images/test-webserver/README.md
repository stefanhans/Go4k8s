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

    kubectl create -f DeployProdWebserver.yaml --record
    
    kubectl get pods,service -l app=webserver,env=production
    kubectl logs -l app=webserver,env=production
    
    curl http://$(minikube ip)::30001
    echo "http://$(minikube ip)::30001
    
    
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


### Prepare Automated Test Container

Prerequisites: 

Build Image for test container and pushed to Docker Hub

    docker build -f Dockerfile.test -t stefanhans/test-webserver .
    docker push stefanhans/test-webserver  
    
Kubernetes Production: Deploy and Test

    kubectl create -f DeployProdWebserver.yaml --record 
    
    kubectl get pods,service -l app=webserver,env=production
    
    curl http://$(minikube ip):30001
    

Go: Push New Version of 'main.go' to GitHub

    git add main.go
    git commit -m "Test version"
    git push 
    
    
Docker: Test Image, build, and push to Docker Hub

    docker run --rm --name test-webserver-container --publish 8080:8080 stefanhans/test-webserver
    
    curl http://localhost:8080

    docker stop test-webserver-container

Choose new image tag for staging and production, respectively.

    docker build -t stefanhans/webserver:1.0.4 .
    docker push stefanhans/webserver:1.0.4

    
Next Step: Try out the [DevOps deployment programmed in Go](../../Deployments/dev-ops)
    