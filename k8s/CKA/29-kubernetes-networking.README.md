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


### Controlplane

If you are on a controlplane node and run netstat you will find an output like so
- `root@controlplane:~# netstat -tuple` OR `root@controlplane:~# netstat -ntuple`

```
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       User       Inode      PID/Program name
tcp        0      0 localhost:10248         0.0.0.0:*               LISTEN      root       9544937    4586/kubelet
tcp        0      0 localhost:10249         0.0.0.0:*               LISTEN      root       9560353    5707/kube-proxy
tcp        0      0 localhost:39657         0.0.0.0:*               LISTEN      root       9524110    4586/kubelet
tcp        0      0 localhost:2379          0.0.0.0:*               LISTEN      root       9514815    3853/etcd
tcp        0      0 controlplane:2379       0.0.0.0:*               LISTEN      root       9514814    3853/etcd
tcp        0      0 127.0.0.11:33867        0.0.0.0:*               LISTEN      nobody     9489915    -
tcp        0      0 controlplane:2380       0.0.0.0:*               LISTEN      root       9514813    3853/etcd
tcp        0      0 localhost:2381          0.0.0.0:*               LISTEN      root       9526843    3853/etcd
tcp        0      0 0.0.0.0:http-alt        0.0.0.0:*               LISTEN      root       9492874    733/ttyd
tcp        0      0 localhost:10257         0.0.0.0:*               LISTEN      root       9543713    3845/kube-controlle
tcp        0      0 localhost:10259         0.0.0.0:*               LISTEN      root       9501689    3382/kube-scheduler
tcp        0      0 127.0.0.53:domain       0.0.0.0:*               LISTEN      systemd-resolve 9508960    492/systemd-resolve
tcp        0      0 0.0.0.0:ssh             0.0.0.0:*               LISTEN      root       9477047    736/sshd
tcp6       0      0 [::]:10250              [::]:*                  LISTEN      root       9537017    4586/kubelet
tcp6       0      0 [::]:6443               [::]:*                  LISTEN      root       9517667    3833/kube-apiserver
tcp6       0      0 [::]:10256              [::]:*                  LISTEN      root       9564417    5707/kube-proxy
tcp6       0      0 [::]:ssh                [::]:*                  LISTEN      root       9477049    736/sshd
tcp6       0      0 [::]:8888               [::]:*                  LISTEN      root       9559050    4898/kubectl
udp        0      0 127.0.0.53:domain       0.0.0.0:*                           systemd-resolve 9508959    492/systemd-resolve
udp        0      0 0.0.0.0:8472            0.0.0.0:*                           root       9538545    -
udp        0      0 127.0.0.11:36257        0.0.0.0:*                           nobody     9489914    -
```

If you wanna find out which device is getting used for node to node communication
`ip route show to match <node_ip>/32` OR `ip addr | grep -C3 <node_ip>`

To see the state of the bridge network on a node
`ip link show docker0`

To see the default gateway, as in how packets from a node leave the system and go out
`ip route show default`

To find the port a component is listening to in master
`netstat -ntuple | grep -i kube-scheduler`

etcd listens to 2379/2380 (both ports). To find the number of open connections for etcd on port 2379
`netstat -anp | grep -i etcd | grep 2379 | wc -l`
In etcd 2380 is for peer-to-peer connectivity.

How to check the network partition of the nodes
```
ip addr | grep eth0
apt-get -y -qq update; apt-get install -y ipcalc
ipcalc 10.3.122.6/24
```

Quite often the kube-system controlplane pods scheduled on the controlplane nodes get the same IP as the node itself.
Make sure to check for this. But you might see some other pods on the controlplane not have the same IP as the controlplane node.
For example the coredns pod. Check the IP of these pods. If the IP seems to be on a different CIDR than the node ip,
then there is a overlay network in play here. To get an idea of the CIDR range of that network partition try this command.
This will give you an idea of what device is used + the ip range as well.
`ip route show to match <pod_ip>/32`

However, if weave has been deployed as the networking solution, then specifically you could be looking for something like this
`kubectl logs -n kube-system <weave-net-pod> -c weave | grep -i range`

Use these to get information like the ip range of the k8s services, the cluster CIDR range, what mode kube-proxy is running as etc.
```
kubectl config set-context --current --namespace=kube-system
kubectl describe pod kube-apiserver-controlplane
kubectl logs kube-apiserver-controlplane -c kube-apiserver
kubectl logs <kube-proxy-pod-name> -c kube-proxy
kubectl get configmaps kube-proxy -o yaml
kubectl get daemonset
```

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
`ip addr add 10.244.2.1/24 dev v-net-0` -> on node02
`ip addr add 10.244.3.1/24 dev v-net-0` -> on node03

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

For IP Address Management (IPAM) CNI comes with 2 builtin plugins - DHCP and host-local
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





### Kubernetes services and kube-proxy
Kubernetes services are not any specific resource, but just IPs and rules allocated to forward traffic to pods
whether from other pods inside the cluster or from some external source outside the cluster.

The component in each node that is responsible for creating the IP addresses and forwarding rules is `kube-proxy`.
It is not just an IP. It is the IP and port combination whether that's ClusterIP for in-cluster communication
or NodePort for external communication. Iptable rules get created and deleted whenever a service gets created or deleted.

When configuring the `kube-proxy` service the `--proxy-mode` can be set to one of userspace|ipvs|iptables.
`iptables` is the default mode.
You can find the `Proxier` mode in the logs of the kube-proxy pod.
Or inspect the pod in yaml format and get it's config and cat the config from inside the pod to see if the mode has been overriden.

Lets say there's a pod called backend-app with an ip `10.244.3.3` and we created a service which got an ip `10.97.111.21`.
As you can see the IPs of the pod seem to be in a different range as compared to the service IP range.

The service CIDR range need to be specified when starting up the `kube-apiserver` component.
  - via the component `--service-cluster-ip-range <CIDR range|default - 10.0.0.0/24>`

This IP range provided for the services must not overlap with the pod networking CIDR range provided in the CNI plugin.
If it is a weave network, you can find that CIDR range by checking the weave pod logs
`kubectl logs <pod-name> -n kube-system -c weave | grep ipalloc`

To see the rules created in the iptables NAT table output
```
iptables -L -t nat | grep backend-app-service

KUBE-SVC-ABCDEFGH  tcp  ---  anywhere  10.97.111.21    /* default/backend-app-service: cluster IP */   tcp  dpt:8080
DNAT               tcp  ---  anywhere  anywhere        /* default/backend-app-service: */              tcp  to:10.244.3.3:8080
KUBE-SEP-PQRSTXYV  tcp  ---  anywhere  anywhere        /* default/backend-app-service: */              tcp  dpt:8080
```
These rules above comments on them with the name of the service.
The first 2 rules mean any traffic coming to the service IP 10.97.111.21 on port 8080 should get forwarded
to the pod with IP 10.244.3.3 on port 8080.
Of course this is a cluster IP rule. The rules will be slightly different for NodePort where any traffic coming to
the node IP on a certain port will get forwarded to the pod IP on a specified port.

And lastly the logs for kube-proxy are generally visible in `/var/log/kube-proxy.log`





### Kubernetes DNS
There is DNS records of the nodes themselves which is managed externally either through your own DNS server
or via a cloud resource like route53.
There is also DNS entries for pods and finally DNS entries of services.

Kubernetes deploys a DNS server for in-cluster DNS resolution - `coredns`. Pre-1.12 kubernetes used to use kube-dns instead of coredns.

If a pod is trying to reach a service in the same namespace then the DNS is simply the name of the service,
ie for example `curl http://backend-service`.
If a pod is trying to reach a service in a separate namespace then the DNS is the <service-name>.<namespace-name>.
ie for example `curl http://backend-service.prod`. So services in a namespace are grouped in a subdomain with the namespaces' name.

All services are in turn grouped under another subdomain called `svc`.
ie `curl http://backend-service.prod.svc`.

Finally all services and pods are grouped under a root domain called `cluster.local`.
So the FQDN of a service will be `curl http://backend-service.prod.svc.cluster.local`.

DNS entries are not created for pods by default but they can be enabled.
If enabled, a pod with an ip of 10.244.20.11 will have a DNS record of
`curl http://10-244-20-11.prod.pod.cluster.local`.



The PODS in the cluster have their `/etc/resolv.conf` files pointing to the coredns server so that
other pod and service DNS records can be resolved. The DNS server stores service names to pod IP mapping
for service DNS entries and dashed ip name to ip mapping for pod DNS entries.

The coredns pods run an executable called `./Coredns` with configuration contained in a file called Corefile.
An example Corefile would look like this
```
.:53 {
  errors
  health
  kubernetes cluster.local in-addr.arpa ip6.arpa {
    pods insecure
    upstream
    fallthrough in-addr.arpa ip6.arpa
  }
  prometheus :9153
  proxy . /etc/resolv.conf
  cache 30
  reload
}
```
The directives like `errors`, `health`, `kubernetes` etc are all plugins of coredns.
The plugin that makes coredns work with kubernetes is `kubernetes` and as you can see that's where the top level
domain name for the cluster is being set as arguments passed to that plugin.
In the options passed to the kubernetes plugin, the `pods` option is what tells coredns kubernetes plugin to create DNS records for pods.
The Corefile of coredns is passed as a configmap to the coredns deployment.
You can get the configmap via - `kubectl get configmap coredns -n kube-system -o yaml`.
Everytime a new pod or service is created or deleted the DNS entries are created or deleted simultaneously by coredns.

When coredns is deployed in the cluster it creates a `service` by the name of `kube-dns` by default.
It is the IP of this service which is configured as the `nameserver` in `/etc/resolv.conf` on the PODs.
This is done by the kubelet when it creates a pod in the node. The config file on the kubelet has the required options set.
`cat /var/lib/kubelet/config.yaml`
```
...
clusterDNS:
  - 10.96.0.10
clusterDomain: cluster.local
...
```

The `/etc/resolv.conf` file in the pod looks similar to this
```
nameserver 10.96.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
```

So, from within a pod if you run `host backend-service` it will return the FQDN of the service.
Because of the `search` directive set in the pods resolv.conf file.
However the search entry is only for service. The DNS for a pod cant be resolved via a shortname.
For that you need to enter the full FQDN like `host 10-244-20-11.prod.pod.cluster.local`



Some helpful commands to get the coredns config
```
kubectl get deployment coredns -n kube-system -o yaml
kubectl get configmap coredns -n kube-system -o yaml
```
