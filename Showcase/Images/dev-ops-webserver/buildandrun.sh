#!/bin/sh +x

ls -l main.go
echo "curl https://raw.githubusercontent.com/stefanhans/Go4k8s/master/Showcase/Images/webserver/main.go > main.go"
curl https://raw.githubusercontent.com/stefanhans/Go4k8s/master/Showcase/Images/webserver/main.go > main.go
ls -l main.go

echo "gofmt -w ."
gofmt -w .

echo "golint ./..."
golint ./...

echo "go build -o ./webserver main.go"
go build -o ./webserver main.go

echo "./webserver &"
./webserver
