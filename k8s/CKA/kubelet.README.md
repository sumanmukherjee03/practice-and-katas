## kubelet

kubelet talks to the api server in the master node. Especially to the scheduler to schedule pods on nodes.
They also report status of the containers to the master.
kubelet registers a node with the kubernetes cluster.

If you use the kubeadm tool it does not automatically deploy the kubelet.

You must always manually install the kubelet on the worker node.
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kubelet
```

The kubelet is run as a service in the worker node.

You can view the kubelet service in the worker node through
`cat /etc/systemd/system/kubelet.service`

You can see the process via
`ps -aux | grep kubelet`

The configuration for kubelet is available in
`cat /var/lib/kubelet/kubelet-config.yaml`
