## pod

A pod wraps around a container.
You can have multiple containers of different kinds in a single pod.
Since containers in the same pod share the same network namespace, they can talk to each other via `localhost`.

To run a single standalone pod for example
`kubectl run nginx --image nginx --restart=Never`

`kubectl get pods`

The status of the pod changes from `ContainerCreating` -> `Running`

Structure of pod-definition.yaml

A kubernetes definition file always has 4 top level definition fields.
```
apiVersion:
kind:
metadata:
spec:
```

For versions :
POD -> v1
Service -> v1
ReplicaSet -> apps/v1
Deployment -> apps/v1

For kind :
Pod, ReplicaSet, Deployment etc

Metadata is in the form of a dictionary :
```
name: nginx-pod
labels:
  app: nginx
```


pod-definition.yaml
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: nginx-container
      image: nginx
```

```
kubectl create -f pod-definition.yaml
kubectl get pods
kubectl describe pod nginx-pod
```

To filter pods based on selectors
```
kubectl get pods --show-labels
kubectl get pods --selector app=nginx
kubectl get pods -l app=nginx
kubectl get pods --no-headers --selector env=prod,bu=finance,tier=frontend
```

```
kubectl describe pod nginx-pod | grep -i image
kubectl describe pods -o wide
kubectl delete pod nginx-pod
kubectl run redis --image=redis --dry-run=client -o yaml > redis-pod-definition.yaml
```

Imperative command to start a standalone pod
```
kubectl run redis --image=redis:alpine --labels=tier=db
```

To create a pod and expose it's pod via cluster ip service in 1 single command
```
kubectl run httpd --image=httpd:alpine --port 80 --expose
```
