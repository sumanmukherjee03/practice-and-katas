## kube-controller-manager

The kube controller manager manages all the controllers.
A controller is a process that watches the status of various functionalities/components of the cluster
and remidiates situations by working towards bringing the cluster component to a desired working state.

For example the node controller manages nodes.
  - Node monitor period - every 5s (healthcheck)
  - Node monitor grace period - 40s (waits for healthcheck to turn healthy again after a failure)
  - Pod eviction timeout - 5m (reshuffles pods onto healthy nodes)

Some common controllers
  - node-controller
  - replication-controller
  - deployment-controller
  - namespace-controller
  - job-controller
  - cronjob-controller
  - endpoint-controller
  - service-account-controller
  - stateful-set-controller

All these controllers are packaged as a single process called controller manager.

Options for installing the kube-controller-manager are available here :
https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/

If we are not using kubeadm, kube-controller-manager is available as a binary to download
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kube-controller-manager
```

The kube-controller-manager runs as a linux service.
`cat /etc/systemd/system/kube-controller-manager.service`
OR
look at the running process
`ps -aux | grep kube-controller-manager`

Some options are :
```
--node-monitor-period=5s
--node-monitor-grace-period=40s
--pod-eviction-timeout=5m0s
--controllers *
```
The `--controllers` option tells us which controllers are enabled. All are enabled by default.


If deployed via kubeadm the kube-controller-manager is available as a pod in the kube-system namespace as `kube-controller-manager-master`
You can view the pod definition file located in
`cat /etc/kubernetes/manifests/kube-controller-manager.yaml`.
