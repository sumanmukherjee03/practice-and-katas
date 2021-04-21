## cluster upgrade

If a node goes offline and then came back online almost immediately, then the kubelet starts and restarts the pods on that node.

However, if the node was offline for 5 mins or more then the pods are terminated from that node.
If the pods on the node above were part of a replica set or a deployment then they are recreated or scheduled on other nodes.

The time it takes for the controller manager to deem a pod as dead is known as the pod eviction timeout and is set on the kube-controller-manager.
`kube-controller-manager --pod-eviction-timeout=5m0s ...`

If the node comes up online after the pod eviction timeout was over, it comes up with no pods on it because the pods got rescheduled elsewhere.

The better way to upgrade a node is to start by draining the node.
`kubectl drain node01`.
This evicts pods on a node and are recreated on other nodes.
Along with that the node is also marked as cordoned or unschedulable.
Once you perform the upgrade on that cordoned node and restart it, you need to uncordon it so that pods can be scheduled on it again.
`kubectl uncordon node01`

If you want to debug things on a node, then you can cordon it.
`kubectl cordon node01`.
This doesnt evict the pods on that node but simply does not schedule any new pods on it.

