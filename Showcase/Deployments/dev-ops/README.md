Prerequisites: Tested [test-webserver image](../../Images/test-webserver).

Build and run Go executables

    go build -o ./blue-green-deployment blue-green-deployment.go 
    ./blue-green-deployment
    
    
View service labels

    kubectl get service -l app=webserver -o jsonpath='{.items[0].metadata.name}{": env="}{.items[0].spec.selector.env}{"\n"}{.items[1].metadata.name}{": env="}{.items[1].spec.selector.env}{"\n"}'
    
    
View container images

    kubectl get pods -l env=production -o jsonpath='{.items[*].spec.containers[0].image}{"\n"}'
    
    
View rollout history of production deployment
    
    kubectl rollout history deploy/webserver-prod-deployment

