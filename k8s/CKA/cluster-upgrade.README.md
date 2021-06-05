## cluster upgrade

Follow the section below if you want to put your node into maintenance to perform some OS upgrade or patches

### node maintenance

If a node goes offline and then came back online almost immediately, then the kubelet starts and restarts the pods on that node.

However, if the node was offline for 5 mins or more then the pods are terminated from that node.
If the pods on the node above were part of a replica set or a deployment then they are recreated or scheduled on other nodes.

The time it takes for the controller manager to deem a pod as dead is known as the pod eviction timeout and is set on the kube-controller-manager.
`kube-controller-manager --pod-eviction-timeout=5m0s ...`

If the node comes up online after the pod eviction timeout was over, it comes up with no pods on it because the pods got rescheduled elsewhere.

The better way to upgrade a node is to start by draining the node.
`kubectl drain node01 --ignore-daemonsets`.
This evicts pods on a node and are recreated on other nodes.
Along with that, the node is also marked as cordoned or unschedulable.

If there are pods on this node that are not controlled by a replica set or deployment, then the node needs to be drained by force.
But remember, that standalone pod running in that node will be lost forever.
`kubectl drain node01 --ignore-daemonsets --force`.

Once you perform the upgrade on that cordoned node and restart it, you need to uncordon it so that pods can be scheduled on it again.
`kubectl uncordon node01`

If you want to debug things on a node, then you can cordon it.
`kubectl cordon node01`.
This doesnt evict the existing pods on that node but simply does not schedule any new pods on it.

### backup

You can backup most resources (not all) in a massive monolithic yaml file like so
`kubectl get all --all-namespaces -o yaml > everything.yaml`
You can also backup your cluster using velero from HeptIO which can backup the cluster from the kubernetes api.

OR

You can backup the etcd cluster itself
When you configure etcd in the master, you choose a data directory `--data-dir=/var/lib/etcd` .
That's the dir you can choose to backup etcd with your backup tool.

OR

You can take a snapshot of the etcd database.

To interact with etcd server via etcdctl you can setup an alias to make life easier
The paths for the certs and keys here are just examples. You can find them from describing the etcd pod
or by looking at the etcd systemd file which is used to run the etcd service.

`alias etcdctl='ETCDCTL_API=3 etcdctl --cert /etc/kubernetes/pki/etcd/server.crt --cacert /etc/kubernetes/pki/etcd/ca.crt --key /etc/kubernetes/pki/etcd/server.key --endpoints="https://127.0.0.1:2379"'`

To take an etcd cluster snapshot use the commands below
```
ETCDCTL_API=3 etcdctl snapshot save snapshot.db --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --endpoints=https://127.0.0.1:2379
ETCDCTL_API=3 etcdctl snapshot status snapshot.db --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --endpoints=https://127.0.0.1:2379
```

### cluster upgrade

Permissible versions of components during a cluster upgrade. The version is controlled by the kube-apiserver.
Here if X is 1.16, then X-1 will be 1.15 and X+1 will be 1.17
```
kube-apiserver : X
  controller-manager : permissible till X-1
  kube-scheduler : permissible till X-1
    kubelet : permissible till X-2
    kube-proxy : permissible till X-2
kubectl : permissible from X-1 to X+1
```

To check if kubeadm was used to setup the cluster
`kubeadm token list`

Upgrade can be planned and performed with kubeadm
```
kubeadm upgrade plan
kubeadm upgrade apply v1.20.0
```

Remember, kubeadm does not install or upgrade kubelets. So, that needs to be done manually.
kubeadm will upgrade all the controlplane components (master components) for us.

kubeadm tool follows the same versioning as kubernetes. So, it needs to be upgraded before upgrading the cluster.

NECESSARY STEPS FOR UPGRADES :
____________________________________

To get the latest patch version of kubeadm corresponding to the minor version use the commands below.
This is a prerequisite for the rest of the upgrade process.
```
apt update
apt-cache madison kubeadm
```

1. First upgrade kubeadm and then upgrade the controlplane
```
apt-get -y -qq update
kubectl version --short
apt install -y kubeadm=1.20.0-00
kubeadm version
kubeadm upgrade plan
kubeadm upgrade apply v1.20.0
kubectl get nodes
```

2. If you used kubeadm to setup the cluster then the controlplane would have the kubelet running as a process
because that's what is used to deploy the rest of the controlplane components as pods in the master nodes.
So, you have to manually upgrade the kubelet on the master. But to upgrade kubelet you must drain the node.
```
kubectl drain controlplane --ignore-daemon-sets
apt install -y kubelet=1.20.0-00
apt install -y kubectl=1.20.0-00
systemctl daemon-reload
systemctl restart kubelet
kubectl get nodes
kubectl uncordon controlplane
```
After this, the master node will show up as upgraded to the new version.

-------------------------
Once, the master is done, it's time to upgrade the worker nodes.

1. Drain worker node. Remember this command needs to be run either from master or from outside the cluster.
`kubectl drain node01 --ignore-daemonsets`

2. Upgrade worker node - kubeadm, kubelet versions, upgrade node config so that the kubelet can pick up the new config and then restart kubelet
```
apt-get -y -qq update
apt install -y kubeadm=1.20.0-00
apt install -y kubelet=1.20.0-00
kubeadm upgrade node
systemctl daemon-reload
systemctl kubelet restart
```

3. Make node schedulable again. Remember this command needs to be run either from master or from outside the cluster.
`kubectl uncordon node01`

### restore

To restore the backup from etcd server snapshot, stop the kube-apiserver first because otherwise objects might be getting created while you are restoring
The path of the --data-dir while restoring backup of etcd is the path of the dir where the snapshot.db is getting restored.
This is analogous to creating a new etcd cluster completely, so that new nodes dont try and join old cluster.
```
service kube-apiserver stop
ETCDCTL_API=3 etcdctl snapshot restore snapshot.db --data-dir /var/lib/etcd-from-backup --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt --key /etc/kubernetes/pki/etcd/server.key --endpoints=https://127.0.0.1:2379
```

After this, etcd needs to be configured to use the new data dir for the `etcd.service`.

If etcd is running as a linux service then :
  Edit `vim /etc/systemd/system/etcd.service`
  And change the data dir to be `--data-dir=/var/lib/etcd-from-backup` - the dir where the cluster was restored to.

  Then reload the systemd daemon to load the changes to the systemd module configs and restart services.
  ```
  systemctl daemon-reload
  systemctl etcd restart
  systemctl kube-apiserver restart
  ```

If etcd is running as a static pod in the controlplane, then update the static pod definition
  Edit `vim /etc/kubernetes/manifests/etcd.yaml`
  And change the host path of the volume mount for data dir to be `/var/lib/etcd-from-backup` - the dir where the cluster was restored to.

  This will recreate the etcd pod


### certs

The certificates are important to be preserved so that the cluster can be fully restored if we loose the master completely
If you dont keep a backup of the certs, you might need to regenerate the entire thing.

These are the sample locations of the kubernetes apiserver certs, cacert and server key.
The etcd cert, cacert and server key could be the same as kube apiserver ones.

```
/var/lib/kubernetes/ca.pem
/var/lib/kubernetes/kubernetes.pem
/var/lib/kubernetes/kubernetes-key.pem
```

OR

```
/etc/kubernetes/pki/etcd/server.cert
/etc/kubernetes/pki/etcd/ca-cert.key
/etc/kubernetes/pki/etcd/server.key
```

### valero

Check out this video to learn more about how valero can backup and restore your workloads in kubernetes.
https://www.youtube.com/watch?v=zybLTQER0yY
