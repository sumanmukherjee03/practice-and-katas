## Etcd

Listens on port 2379
Stores information about the state of the cluster - pods, configs, secrets, accounts, roles, bindings etc.
Every kubectl get command fetches information from etcd
So, essentially etcd is kubernetes' database. Every change ultimately persists in this.

### Installation

```
curl -L https://github.com/etcd-io/etcd/releases/download/v3.2.32/etcd-v3.2.32-linux-amd64.tar.gz -o /tmp/etcd-v3.2.32-linux-amd64.tar.gz
tar xzvf /tmp/etcd-v3.2.32-linux-amd64.tar.gz -C /tmp/etcd-download-test --strip-components=1
rm -f /tmp/etcd-v3.2.32-linux-amd64.tar.gz
mv /tmp/etcd-download-test/etcd /usr/local/bin/etcd
etcd --version
```

This is what `etcd.service` file would look like
```
ExecStart=/usr/local/bin/etcd \\
  --name ${ETCD_NAME} \\
  --cert-file=/etc/kubernetes/pki/etcd/server.crt \\
  --key-file=/etc/kubernetes/pki/etcd/server.key \\
  --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \\
  --peer-cert-file=/etc/kubernetes/pki/etcd/server.crt \\
  --peer-key-file=/etc/kubernetes/pki/etcd/server.key \\
  --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.crt \\
  --per-client-cert-auth \\
  --client-cert-auth \\
  --initial-advertise-peer-urls https://${INTERNAL_IP}:2380 \\
  --listen-peer-urls https://{INTERNAL_IP}:2380 \\
  --listen-client-urls https://{INTERNAL_IP}:2379,https://127.0.0.1:2379 \\
  --advertise-client-urls https://{INTERNAL_IP}:2379 \\
  --initial-cluster-token etcd-cluster-0 \\
  --initial-cluster controller-0=https://{CONTROLLER0_IP}:2380,controller-1=https://{CONTROLLER1_IP}:2380 \\
  --initial-cluster-state new \\
  --data-dir=/var/lib/etcd
```

When setting up etcd, an important option is `--advertise-client-urls https://{INTERNAL_IP}:2379`. This is the address to which etcd listens.
This should be configured on the kube-apiserver during configuration.

In a HA environment when starting etcd, you must pass another option to the etcd binary to let it know the neighbouring etcd nodes
`--initial-cluster controller-0=https://${CONTROLLER0_IP}:2380,controller-1=https://${CONTROLLER1_IP}:2380`
This is essentially the address of the other 2 nodes required for etcd clustering.


### etcdctl

The kubernetes etcd certs are available in /etc/kubernetes/pki/etcd inside `etcd-master` pod if setup via kubeadm.
So, make sure to use these options in `etcdctl` below for authentication with the etcd server via etcdctl in etcd-master
```
--cacert /etc/kubernetes/pki/etcd/ca.crt
--cert /etc/kubernetes/pki/etcd/server.crt
--key /etc/kubernetes/pki/etcd/server.key
```

* Commands *

```
etcd --version
etcdctl --version
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 put foo bar # To store a key/value in the database
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 get foo # To retrieve a key/value from the database
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 cluster-health
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 endpoint health
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 get / --prefix --keys-only --limit 10 # To list all available keys in the database
ETCDCTL_API=3 etcdctl --endpoints=localhost:2379 snapshot save # Takes a backup. This is only for version 3 of the api. It's a different command in version 2.
```
In etcd api version 2, which is the default version, the backup is `etcdctl --endpoints=localhost:2379 backup`

If kubernetes is deployed with kubeadm and not from scratch, etcd is deployed as a pod in the kube-system namespace - `etcd-master`.
```
kubectl get pods -n kube-system
```

To get a list of all the keys in the etcd cluster if running as a kubernetes pod
```
kubectl exec etcd-master -n kube-system -- sh -c "ETCDCTL_API=3 etcdctl get / --prefix --keys-only --limit=10 --cacert /etc/kubernetes/pki/etcd/ca.crt --cert /etc/kubernetes/pki/etcd/server.crt  --key /etc/kubernetes/pki/etcd/server.key"
```
Output from the above command for example : `/registry/apiregistration.k8s.io/apiservices/v1.authorization.k8s.io`.
This is a sample key stored in etcd.
