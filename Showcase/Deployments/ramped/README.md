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

`kubectl rollout undo deploy my-app`
`kubectl rollout pause deploy my-app`
`kubectl rollout resume deploy my-app`



