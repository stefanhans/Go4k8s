### Work In Progress

[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/stefanhans/Go4k8s/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped?status.svg)](https://godoc.org/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)
[![Go Report Card](https://goreportcard.com/badge/github.com/stefanhans/Go4k8s/tree/master/Showcase/Deployments/ramped)](https://goreportcard.com/report/github.com/Go4k8s/tree/master/Showcase/Deployments/ramped)


[Ramped Deployment](https://github.com/ContainerSolutions/k8s-deployment-strategies/blob/master/ramped/README.md)

Using [test-webserver](https://github.com/stefanhans/Go4k8s/tree/master/Showcase/Images/test-webserver) as template

### 1. Do the ramped deployment manually

Initial deployment: `kubectl apply -f app-v1.yaml`

Check deployment: `kubectl get all -l app=my-app`

Open in browser `minikube service my-app --url`

Set <ip:port>: `service=$(minikube service my-app --url)`

Watch app in another terminal: `while sleep 0.5; do curl -s "$service" | grep Version; done`

Watch pods in another terminal: `watch kubectl get pods -l app=my-app`

Rolling Update: `kubectl apply -f app-v2.yaml`

Cleanup: `kubectl delete all -l app=my-app`


Optional:

Rollback: `kubectl rollout undo deploy my-app`

Pause: `kubectl rollout pause deploy my-app`

Resume: `kubectl rollout resume deploy my-app`

### 2. Encapsulate Initial Deployment by k8s/client-go in Docker container

Explore `kubectl --v=10 create -f app-v1.yaml > app-v1.out 2>&1`




