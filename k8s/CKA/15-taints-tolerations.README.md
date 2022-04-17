## taints and tolerations

By default pods cant be scheduled on a tainted node unless the pod has a toleration for the taint.

To taint a node in the cluster, run a command like this :

`kubectl taint nodes node01 key=value:taint-effect`

The taint effects can be NoSchedule|PreferNoSchedule|NoExecute
  - NoSchedule : Don't schedule any new pods but don't evicts existing ones
  - PreferNoSchedule : Try not to schedule any new pods, but it is not guaranteed
  - NoExecute : Evict existing pods running on this node

For example to taint and isolate a node completely
```
kubectl taint node node01 color=blue:NoExecute
```

And when inspection is finished to remove the taint again
```
kubectl taint node node01 color=blue:NoExecute-
```

To add a toleration to a pod add it to the pod definition yaml.
Remember, all the values of the map in `tolerations` need to be encoded in double quotes.
```
apiVersion: v1
kind: Pod
metadata:
  name: frontend-app-pod
spec:
  containers:
    - name: nginx-container
      image:nginx
  tolerations:
    - key: "color"
      operator: "Equal"
      value: "blue"
      effect: "NoSchedule"
```

To find the format of the tolerations object specification
`kubectl explain pod --recursive | less`
`kubectl explain pod --recursive | grep -A5 tolerations`

Taints and tolerations are only meant for nodes from accepting or not accepting pods.
It does not mean that a pod with a toleration wont get scheduled on another node.

The master node is tainted at startup such that no pod is scheduled on this node.
`kubectl describe node kubemaster | grep -i taint`
This will show a taint set as `node-role.kubernetes.io/master:NoSchedule`

It is not advised to start scheduling pods on master, but if you remove the taint on it then new pods will start scheduling there.
`kubectl taint node controlplane node-role.kubernetes.io/master:NoSchedule-`
