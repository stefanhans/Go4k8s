The images and deployments refers to the following showcases of the talk "[K8s API & Go Programming](http://go-talks.appspot.com/github.com/stefanhans/go-present/slides/Kubernetes/IntroductionIntoClient-Go.slide#1)".

- "Hello World" needs [./Images/hello-world](Images/hello-world) and [./Deployments/hello-world](Deployments/hello-world)


- A simple webserver needs [./Images/webserver](Images/webserver) and [./Deployments/webserver](Deployments/webserver)


- A minimalistic CI/CD approach to deploy the simple webserver needs (addition to the webserver) [./Images/dev-ops-webserver](Images/test-webserver) and [./Deployments/dev-ops](Deployments/dev-ops)

All is hardcoded to the [GitHub](https://github.com/stefanhans) and [Docker Hub](https://hub.docker.com/search/?isAutomated=0&isOfficial=0&page=1&pullCount=0&q=stefanhans&starCount=0) accounts of 'stefanhans', and the path to the '.kube/config' as well.

## After the talk...

[Deployment Strategies](https://github.com/ContainerSolutions/k8s-deployment-strategies)

### ramped: release a new version on a rolling update fashion, one after the other
