Prerequisites: Tested [hello-world image](../../Images/hello-world).

Build and run Go executable

    go build -o ./deploy_helloworld deploy_helloworld.go    
    chmod +x deploy_helloworld
    ./deploy_helloworld
    
Investigating...

    kubectl get pods -l app=hello-world 
    kubectl logs -l app=hello-world
    
Exercise: Create a "Hello World" [job which runs to completion](https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/).