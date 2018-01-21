Prerequisites: Tested [webserver image](../../Images/webserver).

Build and run Go executables

    go build -o ./deploy_webserver deploy_webserver.go
    ./deploy_webserver
    
    # Go on interactively
    
Investigating...

    kubectl get pods -l app=webserver
    kubectl logs -l app=webserver