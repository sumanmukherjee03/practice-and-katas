## affinity and anti-affinity

By default pods can get scheduled on any worker node.
To make certain pods get scheduled only on certain nodes or not get scheduled on certain nodes you can use affinity and anti-affinity.

For example one way to make data-processing pods to get scheduled to a node with `labels` - `size: large`,
we can make use of the `nodeSelector` property in the pod definition yaml.
```
apiVersion: v1
kind: Pod
metadata:
  name: data-processing-pod
spec:
  containers:
    - name: data-processing-container
      image: data-processor
  nodeSelector:
    size: large
```

To label nodes
`kubectl label nodes node01 size=large`

The node obviously needs to be labelled before the creation of the pod.

The node selectors are useful for simple cases, but for more complex scheduling/placements we require node affinity or anti-affinity.
For example :
```
apiVersion: v1
kind: Pod
metadata:
  name: data-processing-pod
spec:
  containers:
    - name: data-processing-container
      image: data-processor
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
          - matchExpressions:
            - key: size
              operator: In
              values:
                - large
```
This essentially does the same thing as the nodeSelector did above. Place a new pod of data processing on the large node.

You can have multiple types of operators. For example to place the pod in nodes labelled either large or medium.
```
....
        nodeSelectorTerms:
          - matchExpressions:
            - key: size
              operator: In
              values:
                - large
                - medium
```

OR if you want to place the pod on any node that is not labelled small

```
....
        nodeSelectorTerms:
          - matchExpressions:
            - key: size
              operator: NotIn
              values:
                - small
```

OR if you want to place the pod on any node that has any value for the label size

```
....
        nodeSelectorTerms:
          - matchExpressions:
            - key: size
              operator: Exists
```

The type of node affinity can be of 2 kinds :
  - requiredDuringSchedulingIgnoredDuringExecution
  - preferredDuringSchedulingIgnoredDuringExecution
These are like hard and soft requirements. If `requiredDuringSchedulingIgnoredDuringExecution` is used
and a node is not found to place the pod on, then the pod will not be placed at all. If there are already running pods
before the nodes were labelled or affinity rules were applied, those will not be evicted.

Another node affinity type to come in the future is `requiredDuringSchedulingRequiredDuringExecution`
which will evict pods from nodes that don't satisfy the affinity requirements.

You can learn more about the complex rules of node affinity here : https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/
