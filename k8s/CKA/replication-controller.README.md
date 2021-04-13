## replication-controller

The replication controller helps us run multiple instances of a single pod for HA.
Even with a single pod a replication controller can help by bringing up a new pod when an existing one fails.
Also, having a replica set helps share the load and scale our application.

Replication controller is the older technology which is being replaced by Replica set to set up replication.

rc-definition.yaml
```
apiVersion: v1
kind: ReplicationController
metadata:
  app: app-rc
  labels:
    app: frontend
spec:
  replicas: 3
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

Create the replication controller
`kubectl create -f rc-definition.yaml`

Some helpful commands for replication controllers
```
kubectl get replicationcontroller
```

-----------------------------------------------------------------------

ReplicaSet is similar to replicationcontroller
There are a couple of differences between a ReplicaSet and replicationcontroller :
  - The apiVersion is different in this case `apps/v1`, and not `v1`.
  - The replicaset needs a selection to match pods
    - ReplicaSet can also manage pods that were not created as part of the replicaset spec definition.


replicaset-definition.yaml
```
apiVersion: apps/v1
kind: ReplicaSet
metadata:
  app: app-replicaset
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

Create the replicaset with
`kubectl create -f replicaset-definition.yaml`

Some helpful commands for replicasets
```
kubectl get replicaset
```

The selector > matchLabels is a requirement for the ReplicaSet to filter pods to monitor.
You could have 3 standalone pods created already with a pod definition template and start monitoring them by creating a replicaset.
However, in the replicaset spec you still need the template for the pod definition
because this ensures the replicaset can bring up a new pod with that spec to replace a dead pod.

If you scale up the replicas to 6 from 3 you can run the kubectl replace command.
`kubectl replace -f replicaset-definition.yaml`
OR use the kubectl scale command. This does not update the file though. Just uses the file to get the name of the replicaset object.
`kubectl scale --replicas=6 -f replicaset-definition.yaml`
You could also use the scale command like so for an imperative update
`kubectl scale --replicas=6 replicaset app-replicaset`

To delete the replicaset
`kubectl delete replicaset app-replicaset`
