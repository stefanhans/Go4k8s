Build and run Go executable

    go build -o ./go-hello-world .
    
    ./go-hello-world

Build and run docker image

    docker build -t stefanhans/go-hello-world .

    docker run --rm stefanhans/go-hello-world:latest 
    
Push image to Docker hub

    docker push stefanhans/go-hello-world:latest