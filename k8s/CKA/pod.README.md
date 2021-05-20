## pod

A pod wraps around a container.
You can have multiple containers of different kinds in a single pod.
Since containers in the same pod share the same network namespace, they can talk to each other via `localhost`.

To run a single standalone pod for example
`kubectl run nginx --image nginx --restart=Never`

To run a single standalone pod in a specific node for inspecting things
```
kubectl run ubuntu --image ubuntu --overrides='{"apiVersion": "v1", "spec": {"template": {"spec": {"nodeSelector": {"kubernetes.io/hostname": "node01"}}}}}' --command sleep 300
```

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
kubectl get pod nginx --watch
kubectl get pods --no-headers --selector env=prod,bu=finance,tier=frontend
```

```
kubectl describe pod nginx-pod | grep -i image
kubectl describe pods -o wide
kubectl delete pod nginx-pod
kubectl run redis --image=redis --dry-run=client -o yaml > redis-pod-definition.yaml
kubectl exec app -c webapp -- tail -f /log/webapp.log
```

Imperative command to start a standalone pod
```
kubectl run redis --image=redis:alpine --labels=tier=db
```

To create a pod and expose it's pod via cluster ip service in 1 single command
```
kubectl run httpd --image=httpd:alpine --port 80 --expose
```

You can override commands/entrypoints or provide args of docker containers in pod definition templates
The `command` field is analogous to the `ENTRYPOINT` directive in Dockerfile.
The `args` field is analogous to the `CMD` directive in Dockerfile.
In Dockerfile the combination of `ENTRYPOINT` and `CMD` is
  - use the `ENTRYPOINT` script as the executable with the values of the `CMD` as default arguments.
That's what the `command` and `args` combination does in pod spec.
cat `pod-definition.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu-wait-pod
spec:
  containers:
    - name: ubuntu-wait-container
      image: ubuntu-wait
      command: ["sleep"]
      args: ["10"]
      env:
        - name: COLOR
          value: green
```

There are 4 ways to pass env vars into kubernetes pod spec - key/value pair, config maps, secret keys, pre populated values from fields.
Here we are using two of the simplest way to pass env vars into a kubernetes pod definition.
```
apiVersion: v1
kind: Pod
metadata:
  name: ubuntu-wait-pod
spec:
  containers:
    - name: ubuntu-wait-container
      .....
      env:
        - name: COLOR
          value: green
        - name: HOSTNAME
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: spec.nodeName
```

Below is a crude example of a pod-definition.yaml with init containers

```
apiVersion: v1
kind: Pod
metadata:
  name: webapp-pod
  labels:
    app: webapp
spec:
  containers:
  - name: webapp-container
    image: busybox:1.28
    command: ['sh', '-c', 'echo Starting app && sleep 3600']
  initContainers:
  - name: init-redis
    image: busybox:1.28
    command: ['sh', '-c', 'until nslookup redis-service; do echo waiting for redis to be up and running; sleep 3; done;']
  - name: init-db
    image: busybox:1.28
    command: ['sh', '-c', 'until nslookup db-service; do echo waiting for db to be up and running; sleep 3; done;']
```


-----------------------------------------------------------

### security context
Kubernetes allows you to set security context on pods or individual containers.
Setting on a pod is applied to all containers and setting on a container overrides the setting on a pod.

`cat pod-definition.yaml`

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
  securityContext:
    runAsUser: 1000
```

OR to set the security context at the container level

`cat pod-definition.yaml`

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
      securityContext:
        runAsUser: 1000
        capabilities:
          add: ["MAC_ADMIN", "SYS_TIME"]
```
