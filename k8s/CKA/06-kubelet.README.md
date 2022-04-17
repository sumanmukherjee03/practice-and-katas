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

The kubelet reads static pod definition yamls from `/etc/kubernetes/manifests` by default.

If you want to start pods standalone in a node without submitting the pod definition to the api server,
you can place the pod definition files in this dir. The kubelet has a watcher for changes to files in
this dir and uses that to start/stop/update pods in the host. So, the new pod definitions will get picked up.

The directory the kubelet reads pod definitions from can be specified to the kubelet at startup.
This is set using the kubelet startup option - `--pod-manifest-path`.
Another way to configure this is by passing the kubelet startup option for configuration file
for the kubelet - `--config`. Then in the config file you can setup the dir with the `staticPodPath` key/value pair.

Static pods are reflected in the cluster state because the kubelet provides this information afterwards to the api-server.
However the static pod definitions can be edited using kubectl outside the host where the pod
manifest was placed for the kubelet of that host to pick up. Static pods are useful to
deploy control plane components of a kubernetes cluster like controller-manager.yaml, api-server.yaml etc.
Static pods are restarted by the kubelet if they crash. This mechanism is used by the kubeadm tool.

Kubelet contains another important component called the `cAdvisor`. The `cAdvisor` is responsible
for collecting pod level metrics and exposing them to the metrics server. The metrics server is an
in-memory metrics aggregation server in kubernetes which is also used for horizontal and vertical pod autoscaling.
The metrics-server can be installed via
`kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml`
Having the metrics server enables us to run the `kubernetes top node` or `kubernetes top pod` command to view cluster metrics.
