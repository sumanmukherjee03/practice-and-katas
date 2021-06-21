## deployment

A deployment is a kubernetes object which is at a level above the pods and replicasets.
This object controls how you want to rollout a new image for a container or rollback etc.

A deployment definition looks the same as a ReplicaSet except for the kind which is different.

deployment-definition.yaml
```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app-replicaset
  labels:
    app: frontend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: frontend
  template:
    metadata:
      name: nginx-pod
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx-container
          image: nginx
```

Create the new deployment with
`kubectl create -f deployment-definition.yaml`

Some helpful commands for deployments, replicaset created by it and inturn pods.
```
kubectl get deployments
kubectl get replicaset
kubectl get pods
```

To get all the objects created at once run `kubectl get all`.

To create a sample yaml for deployment, run the command below. This gets us going with a template.
```
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml
```

Another easy way to quickly create and scale a deployment is provided below.
```
kubectl create deployment nginx --image=nginx --replicas=2
kubectl scale deployment nginx --replicas=3
```

To change the image in a deployment and record the change in annotations of the deployment
```
kubectl set image deployment/nginx nginx=nginx:1.18 --record
```

To expose a port from a running container via a service with a cluster IP
```
kubectl expose deployment nginx --port 80
```
