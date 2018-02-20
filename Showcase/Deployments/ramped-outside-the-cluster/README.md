Not using a container from outside so far - skip this

Build app

    GOOS=linux go build -o ./app .

create Dockerfile

    cat >Dockerfile <<EOF
    FROM debian
    COPY ["./app", "./deployment.yaml", "update.yaml", "/home/stefan/.kube/config", "/"]
    ENTRYPOINT /app
    EOF

build image

    docker build -t update-out-cluster .

push image as needed

run docker container

    kubectl run --rm -i demo --image=stefanhans/update-out-cluster:1.0.0 --image-pull-policy=Never


    kubectl exec -it demo -- /bin/bash
