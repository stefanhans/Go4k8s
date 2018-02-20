kubectl create -f update_job.yaml

while [ "$(kubectl get jobs update-job -o jsonpath='{.status.active}')" != "" ]
do
    sleep 1
    printf "."
done
kubectl delete -f update_job.yaml
