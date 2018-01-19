Docker: Build Image and Container

    docker build -t stefanhans/webserver .
    docker create --publish 8080:8080 --name my-go-webserver stefanhans/webserver
    
Docker: Run Container and Verify

    docker start my-go-webserver
    http://localhost:8080
    
Docker: Stop Container

    docker stop my-go-webserver
    
Docker: Push Container to Docker Hub

    docker push stefanhans/webserver
    