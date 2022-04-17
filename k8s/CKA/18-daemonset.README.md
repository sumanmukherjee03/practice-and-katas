## daemonset

DaemonSet is like a ReplicaSet except for the fact that it runs 1 copy of the pod in each node.
Whenever a new node is added to the cluster a replica of the pod is automatically added to the new node.
For example `kube-proxy` is deployed as a daemonset in the kubernetes cluster.

`cat daemon-set-definition.yaml`
```
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: monitoring-daemon
spec:
  selector:
    matchLabels:
      app: monitoring-agent
  template:
    metadata:
      labels:
        app: monitoring-agent
    spec:
      containers:
        - name: monitoring-agent
          image: monitoring-agent
```
Ensure the labels in the selector matches the ones in the pod template section.

To easily get a template for daemonset create one for a deployment with the dry run and then edit it.

`kubectl create -f daemon-set-definition.yaml`

Some useful commands for daemonsets
```
kubectl get daemonsets
kubectl describe daemonsets monitoring-daemon
```
