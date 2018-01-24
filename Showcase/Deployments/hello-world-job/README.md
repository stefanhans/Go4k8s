Prerequisites: Tested [hello-world image](../../Images/hello-world).

Take yaml file as guideline:

    apiVersion: batch/v1
    kind: Job
    metadata:
      name: hello-world-job
    spec:
      template:
        spec:
          containers:
          - name: hello-world
            image: stefanhans/hello-world
            command: ["/hello-world"]
          restartPolicy: Never
      backoffLimit: 4


Build and run Go executable

    go build -o ./deploy_helloworld-job deploy_helloworld-job.go
    ./deploy_helloworld-job
    
Investigating...

    kubectl get all -l app=hello-world-job
    kubectl logs -l app=hello-world-job
