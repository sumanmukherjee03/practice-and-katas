## kube-scheduler

The scheduler only decides which pod goes where. It doesnt actually run the pod in the host. That's the job of the kubelet.

The scheduler picks a node to place the pod on depending on a 2-step process.
  - filter out nodes that do not fit the requirement
  - rank nodes based on a prioritization algorithm and pick the node with the best score
    - for resources, the scoring is based on how much resources would be left on the host if the pod is scheduled there
      - the more resources left after pod placement, the better the score of that node


If we are not using kubeadm, kube-scheduler is available as a binary to download
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kube-scheduler
```

OR view the running process with
`ps -aux | grep kube-scheduler`

When running as a service the kube-scheduler is passed an option `--config` which points to the file containing the configuration.
`cat /etc/kubernetes/config/kube-scheduler.yaml`

if you have set it up with kubeadm tool, then kube-scheduler is available as a pod in the kube-system namespace.
`cat /etc/kubernetes/manifests/kube-scheduler.yaml`
