## Kubernetes networking

kube-api server listens to port 6443. So that port needs to be opened on the master node.
kubelets communicate over port 10250. So that port needs to be opened on the worker nodes and also the master node.
kube-scheduler which runs on the master needs the port 10251 to be open.
kube-controller-manager which runs on the master needs the port 10252 to be open.
kubelets need ports 30000 - 32767 to be open on the worker nodes for external services to be exposed.
etcd which runs on the master also usually needs the port 2379 to be open.
etcd clients communicate with each other on port 2380. So in a HA master environment, you need that port to be open as well.

So, a master would have these ports open : 2379, 2380, 6443, 10250, 10251, 10252
So, a worker would have these ports open : 10250, 30000-32767

https://kubernetes.io/docs/concepts/cluster-administration/addons/
https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/high-availability/#steps-for-the-first-control-plane-node


### Pod networking

Pods must each get an unique IP
Pods must be able to communicate with each other
Pods must be able to reach each other without going through a NAT


Imagine a cluster of 3 worker nodes - 192.168.1.7, 192.168.1.8, 192.168.1.9
These nodes are all part of an external LAN 192.168.1.0/24 lets say.

On each node we create a bridge network to be able to attach the network namespaces to
```
ip link add v-net-0 type bridge
ip link set dev v-net-0 up
```

Lets give the bridge devices in each node an ip address
`ip addr add 10.244.1.1/24 dev v-net-0` -> on node01
`ip addr add 10.244.2.1/24 dev v-net-0` -> on node01
`ip addr add 10.244.3.1/24 dev v-net-0` -> on node01

Next we create a script that runs every time a container comes up on a node or gets deleted from a node- `cat container-net.sh`
```
Add)
  # Create a veth pair
  ip link add ...

  # Attach veth pair
  ip link set ...
  ip link set ...

  # Assign IP address
  ip -n <namespace> addr add ...
  ip -n <namespace> route add ...

  # Bring up the interface
  ip -n <namespace> link set ...

Delete)
  # Delete veth pair
  ip link del ...
```

Running a script like this every time a container comes up allows the network namespaces of the container to communicate with other containers inside one node.

For a container in node01 to be able to communicate with another container of node02 or node03 we need to add a route
```
ip route add 10.244.2.2 via 192.168.1.8
ip route add 10.244.3.2 via 192.168.1.9
```
We have to add similar routes for all other containers on other hosts on each node.


But if there was a router we could have added some routes like this instead of configuring routes on each node

```
Network              Gateway
----------           ---------
10.244.1.0/24        192.168.1.7
10.244.2.0/24        192.168.1.8
10.244.3.0/24        192.168.1.9
10.244.1.0/24        192.168.1.7
```


The `container-net.sh` script is somewhat similar to what a CNI plugin does.
When bringing up the kubelet a config option is passed `--cni-conf-dir=/etc/cni/net.d`.
Based on this config the bin dir is looked up `--cni-bin-dir=/etc/cni/bin`
which points the kubelet to the `./net-script.sh`
and the kubelet in turn invoked the script `./net-script.sh add <container> <namespace>`.

This script is called again when a container gets deleted.
This is pretty much a rough explanation of how the CNI plugins work.



### Configuring CNI

This is configured in the kubelet initialization options. `cat kubelet.service`
```
ExecStart=/usr/local/bin/kubelet \\
  --config=/var/lib/kubelet/kubelet-config.yaml \\
  --container-runtime=remote \\
  --container-runtime-endpoint=unix:///var/run/containerd/containerd.sock \\
  --image-pull-progress-deadline=2m \\
  --kubeconfig=/var/lib/kubelet/kubeconfig \\
  --network-plugin=cni \\
  --cni-bin-dir=/opt/cni/bin \\
  --cni-conf-dir=/etc/cni/net.d \\
  --register-node=true \\
  --v=2
```

OR

`ps -aux | rep kubelet`

Based on the configuration specified in `--cni-conf-dir` an appropriate binary is
picked from `--cni-bin-dir` to run the script for network configuration when a container comes up or dies.

If the conf dir is not listed search in `/etc/cni/net.d`. That would be the default location.

For IP Address Management or IPAM CNI comes with 2 builtin plugins - DHCP and host-local
which can handle unique ip address allocation to the pods.
The cni configuration file (`cat /etc/cni/net.d/net-script.conf`) has a section for ipam where
the plugin that does IPAM is configured under the key called `type`.




### Weave CNI Plugin

The weave CNI plugin deploys an agent on every node. These agents communicate with each other.
These agents are aware of the topology of the nodes and they maintain information of how to route information in between these nodes.
When a tcp packet from one pod in a node tries to reach another pod in another node, this packet is
first intercepted by the weave agent on the node and that packet is then wrapped in
another packet destined for that other node by weave. The receiving nodes weave agent intercepts it again
and unwraps it and passes the original packet to the pod that it is destined for.

Weave creates it's own bridge called WEAVE.
Important to remember that a container network namespace can be attached to multiple bridge networks.
Weave's agent configures the correct route for the pod.

Weave is deployed as daemons on each node in the cluster or as pods on each node via daemonset.
`kubectl apply -f "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"`

To get the weave peers
`kubectl get pods -n kube-system | grep 'weave-net'`
