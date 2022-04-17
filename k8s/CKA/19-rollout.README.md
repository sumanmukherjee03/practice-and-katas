## rollout

Kubernetes maintains versions of objects that are being applied to make it easier to rollback to a previous version.

Here's some useful commands to see the status of the deployment and history of deployments.
```
kubectl rollout status deployment/webapp-deployment
kubectl rollout history deployment/webapp-deployment
kubectl describe deployment deployment/webapp-deployment
```

There are different of types of deployment strategies
  - RollingUpdate : This is the default deployment strategy
  - Recreate : This strategy will involve a downtime for the application since it takes down all pods and brings new ones up

For rolling out a new version of an application, you could update the deployment yaml and kubectl apply.
Another way is by setting the image in a deployment and in the process also record the change
that happened in the annotations of the deployment.
```
kubectl set image deployment/webapp-deployment nginx=nginx:1.9.1 --record
```
This also causes a rolling update but the process is dirty because this change wont be captured in the manifest file which can be version copntrolled.

During deployment kubernetes maintains 2 replicasets, one original and a new one. As pods are created in the new
replicaset, pods are also simultaneously terminated in the old replicaset. To undo a deployment if something went wrong
```
kubectl rollout undo deployment/webapp-deployment
```

Relevant section of a `deployment.yaml` with the rolling update strategy
```
spec:
  .....
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
```
