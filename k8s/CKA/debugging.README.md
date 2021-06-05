## Debugging

----------------------
### Debugging pods

curl http://webapp-service-ip:node-port
kubectl describe svc webapp-service
Compare selectors on the service to the ones on the pods

To see application logs and if necessary to follow the logs live
  kubectl logs web
  kubectl logs web -f

If an app has been failing and recreating a container multiple times, to see the logs from the previous crashed container of the pod
  kubectl logs web --previous

You can also try and debug with kubectl exec command
  kubectl exec cassandra -- cat /var/log/cassandra/system.log

To list all objects in a namespace
  kubectl -n delta get all

List pods for a specific label
  kubectl get pods -l app=webapp

Run a simple busybox pod to see what a pod sees
  kubectl run -it --rm --restart=Never busybox --image=gcr.io/google-containers/busybox sh

To test if the pods are working properly
  kubectl get pods -l app=hostnames -o jsonpath='{{range .items}}{{.status.podIP}}{{"\n"}}{{end}}'
  You can test the IPs individually like so `wget -q0- <pod-ip>`


-------------------------
### Debugging services

To check if the service is actually targeting pods properly.
If there are no endpoints being assigned it could be a problem with selectors in the service that are targeting the pods.
  kubectl get endpoints <service-name>

Check the service IP+port (not the nodeport, but the service port) from one of the nodes and then from inside a pod as well
Both of these cases should be working.
If there is an issue with service reachability from the nodes/pods then it indicates a problem with kube-proxy and in turn with iptables rules.
  curl <service-ip>:<service-port>
  nc -zv <service-ip> <service-port>
  telnet <service-ip> <service-port>

      Check the logs of kube-proxy to see what mode kube-proxy is running as iptables/ipvs etc.
        cat /var/log/kube-proxy.log | grep -i proxier
        OR
        journalctl -u kube-proxy | grep -i proxier

      If kube-proxy is running in iptables mode, then you should be able to see the iptable entries for <service-name>
        iptables-save | grep -i <service-name>
          For each port of each Service, there should be 1 rule in KUBE-SERVICES and one KUBE-SVC-<hash> chain.
          For each Pod endpoint, there should be a small number of rules in that KUBE-SVC-<hash> and
          one KUBE-SEP-<hash> chain with a small number of rules in it.

Test the endpoint IP+port from inside a pod.
An endpoint+port is basically a pod ip + port.
Pods should be able to communicate with each other without going through a NAT.
If this is not working check if there are CNI pods available, like calico/weave/flannel etc. It is probably an issue with the pod networking.
  curl <service-endpoint-ip>:<service-endpoint-port>
  nc -zv <service-endpoint-ip> <service-endpoint-port>
  telnet <service-endpoint-ip> <service-endpoint-port>

To check if a service is resolving DNS from inside a pod, try to get into a busybox pod and run this inside
  nslookup <service-name>.default.svc.cluster.local
    Depending on the kubelets `--cluster-domain` flag the root domain may be `cluster.local` or something else.

Or test the same thing above by trying this from a node, ie by passing a DNS server to resolve the in cluster DNS records.
  nslookup <service-name>.default.svc.cluster.local <kube-dns-service-ip>

If you can resolve FQDN but not the short names of the services from inside pods
  Check /etc/resolv.conf inside the simple busybox pod to see if the domain entries for `search` are there or not
  and if the nameserver is pointing to the kube-dns service or not.




------------------------
### Debugging controlplane

On the controlplane
  kubectl get pods -n kube-system
  OR
  service kube-apiserver status
  service kube-controller-manager status
  service kube-scheduler status

If the controlplane components are installed as linux services and not running as pods you can look at the logs like so
  journalctl -u kube-apiserver
  journalctl -u kube-scheduler
  OR
  tail -f /var/log/kube-apiserver.log
  tail -f /var/log/kube-scheduler.log
  tail -f /var/log/kube-controller-manager.log


ON the worker nodes check the kubelet and kube-proxy if the kube-proxy was installed as a linux service instead
  service kubelet status
  service kube-proxy status

If there are certificate issues on the kube-scheduler or kube-apiserver components
check the paths of the certs. If the components are running as pods check that the proper
directory has been mounted for the certs.

If the controlplane components are running as static pods they will have a `-master` or `-controlplane`
appended to the end of the pod names. Their config can be found in `/etc/kubernetes/manifests`.
Another way to find the configured directory of the static pod manifests is via the command line
options passed to the kubelet.
`cat /etc/systemd/system/kubelet.service.d/10-kubeadm.conf`
Look for the config file passed to the kubelet which could be `/var/lib/kubelet/config.yaml`.
Next search for the conf dir of static pods in the kubelet config file
`grep -i staticPodPath /var/lib/kubelet/config.yaml`



--------------------
### Debugging worker nodes

To get detailed info about the overall cluster
  kubectl cluster-info
  kubectl cluster-info dump

To debug the worker nodes, from the controlplane run these kubectl commands to debug the issue
  kubectl get nodes
  kubectl describe node node01

If on the describe command on a node you see something along the lines of `FailedToStartProxierHealthcheck`
then there is a good chance that the kubelet is not running on the worker node.

On the worker node itself that is having trouble check for issues
  top
  df -h
  service kubelet status

You can view kubelet logs on the worker nodes via
  journalctl -u kubelet
  OR
  tail -f /var/log/kubelet.log

If kube-proxy is running as a linux service instead of a pod, you can view kube-proxy logs on the worker nodes via
  journalctl -u kube-proxy
  OR
  tail -f /var/log/kube-proxy.log

In case there are cert issues, check the kubelet certs for things like the certificate issuer, expiry, group (system:nodes) etc
  openssl x509 -in /var/lib/kubelet/node01.crt -text




--------------------
### Debugging DNS issues

For a coredns deployment there is always a coredns configmap which contains the Corefile content
The `proxy` directive in the Corefile configuration is what tells where to forward "out of cluster" DNS queries.
So, that's an important pointer when pods are having DNS issues when resolving out of cluster DNS.

Check the kubelets `-cni-bin-dir` and `-network-plugin` runtime options if the coredns pod is in a pending state.
This usually means the kubelet is having issues starting pods.

If the coredns pods are in a crash loop
Your node could be running an older version of docker or might have SELinux enabled.
Try to upgrade docker on the node.

You could have SELinux enabled.
To disable SELinux on the node try these commands for an Ubuntu machine
```
# To get the status of SELinux run the command below
sestatus

# Also to get the status of SELinux you can run getenforce. The values can be "enforcing", "permissive" or "disabled".
# The permissive mode only logs but doesnt enforce any SELinux policy
getenforce

# This will change the SELinux mode to permissive, but this change does not persist over a reboot
setenforce 0
```

For disabling SELinux and persisting it across a reboot, make the change in `/etc/selinux/config` with `SELINUX=disabled` directive.

Another reason for coredns to enter into a crashloop could be due to the lack of enough privilege provided to the pod to work properly.
You could update the deployment to change the `allowPrivilegeEscalation` to `true` in this case.
```
kubectl -n kube-system get deployment coredns -o yaml | sed 's/allowPrivilegeEscalation: false/allowPrivilegeEscalation: true/g' | kubectl apply -f -
```

One more reason why coredns pod could be crashing is because coredns detects a loop for dns resolution
Ways to work around the issue :
  - Add `resolvConf: <path_to_real_resolvconf>"` in the kubelet config
    so that the kubelet can pass this to the pods instead of the default `/etc/resolv.conf`
      For systems using `systemd-resolved`, that path is `/run/systemd/resolve/resolv.conf`
  - Disable local dns cache and restore `/etc/resolv.conf` to the original
    This is because ubuntu by default ships with dnsmasq installed and a dns query first hits
    the dns cache. If it is not found there it is forwarded to the DNS servers configured for the domain.
    To disabe dns cache in the system you can do the following
      `vim /etc/NetworkManager/NetworkManager.conf` -> and comment out `#dns=dnsmasq`
      `service network-manager restart`
  - Another fix for coredns is replace `proxy . /etc/resolv.conf` with `proxy . 8.8.8.8` in the coredns config, ie Corefile.
    This however wont fix the issue with the normal pods out of cluster dns queries.
    Because kubelet will still be using the `/etc/resolv.conf` for the dnsPolicy of the pods by default.
      `dnsPolicy` field in pod spec can be one of - `Default`, `ClusterFirst`, `ClusterFirstWithHostNet`, `None`
      The default dnsPolicy is `ClusterFirst`. Any DNS query not matching the cluster domain is forwarded to the upstream nameserver configured on the nodes
        You can read more about it here https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy
        You can pass a custom dnsConfig to the pod spec as well if required.

If the coredns pods are working fine, check the `kube-dns` service. Check if the service has been assigned valid endpoints.
It is possible that the selectors and/or ports of the service are wrong
`kubectl -n kube-system get endpoint kube-dns`




-----------------------------
### Debugging kube-proxy

kube-proxy is responsible for maintaining network rules on the nodes to allow communication to/from the pods.
If kubeadm tool was used to bring up the cluster, then kube-proxy would be running as a daemon set on all nodes
`kubectl -n kube-system describe daemonset kube-proxy`
The command that the kube-proxy pods would be running will look somewhat like
  `/usr/local/bin/kube-proxy --config=/var/lib/kube-proxy/config.conf --hostname-override=$(NODE_NAME) ...`
The config is mounted as a volume inside the pods.
Look for mode, clusterCIDR, bindaddress, kube-config, iptables etc in that configuration.
Check the logs in the kube-proxy pod.
Check the process in the kube-proxy pod - `netstat -plan | grep kube-proxy`




-----------------------------
### Some sample iptables and ip route outputs

This is what iptables would look like for a service with Cluster IP if things are working
```
node01 $ iptables-save | grep mysql
-A KUBE-SEP-BGVORCJSLGKD27QD -s 10.32.0.3/32 -m comment --comment "application/mysql" -j KUBE-MARK-MASQ
-A KUBE-SEP-BGVORCJSLGKD27QD -p tcp -m comment --comment "application/mysql" -m tcp -j DNAT --to-destination 10.32.0.3:3306
-A KUBE-SERVICES ! -s 10.244.0.0/16 -d 10.96.223.183/32 -p tcp -m comment --comment "application/mysql cluster IP" -m tcp --dport 3306 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.96.223.183/32 -p tcp -m comment --comment "application/mysql cluster IP" -m tcp --dport 3306 -j KUBE-SVC-T6EKFQGEJQP7KYGH
-A KUBE-SVC-T6EKFQGEJQP7KYGH -m comment --comment "application/mysql" -j KUBE-SEP-BGVORCJSLGKD27QD
```

This is what iptables would look like for a service with NodePort if things are working
```
node01 $ iptables-save | grep webapp-service
-A KUBE-NODEPORTS -p tcp -m comment --comment "application/webapp-service" -m tcp --dport 30081 -j KUBE-MARK-MASQ
-A KUBE-NODEPORTS -p tcp -m comment --comment "application/webapp-service" -m tcp --dport 30081 -j KUBE-SVC-QB3GCEADICMX7L52
-A KUBE-SEP-WHCOKKDP6FNJV2JR -s 10.32.0.4/32 -m comment --comment "application/webapp-service" -j KUBE-MARK-MASQ
-A KUBE-SEP-WHCOKKDP6FNJV2JR -p tcp -m comment --comment "application/webapp-service" -m tcp -j DNAT --to-destination 10.32.0.4:8080
-A KUBE-SERVICES ! -s 10.244.0.0/16 -d 10.107.112.64/32 -p tcp -m comment --comment "application/webapp-service cluster IP" -m tcp --dport 8080 -j KUBE-MARK-MASQ
-A KUBE-SERVICES -d 10.107.112.64/32 -p tcp -m comment --comment "application/webapp-service cluster IP" -m tcp --dport 8080 -j KUBE-SVC-QB3GCEADICMX7L52
-A KUBE-SVC-QB3GCEADICMX7L52 -m comment --comment "application/webapp-service" -j KUBE-SEP-WHCOKKDP6FNJV2JR
```

This is what a routing table looks like from a host
```
node01 $ ip route
default via 172.17.0.1 dev ens3
10.32.0.0/12 dev weave proto kernel scope link src 10.32.0.1
10.244.0.0/24 via 172.17.0.10 dev ens3
172.17.0.0/16 dev ens3 proto kernel scope link src 172.17.0.16
172.18.0.0/24 dev docker0 proto kernel scope link src 172.18.0.1 linkdown
```

This is what the routing table looks like from within a pod if pod <-> pod communication is working
```
controlplane $ kubectl -n application exec -it webapp-54db464f4f-kqj7l -- /bin/sh
/opt/webapp # ip route
default via 10.32.0.1 dev eth0
10.32.0.0/12 dev eth0 scope link  src 10.32.0.4
```

This is an example of what iptables might look like for a service when kube-proxy isn't in a working state
```
node01 $ iptables-save | grep webapp-service
-A KUBE-EXTERNAL-SERVICES -p tcp -m comment --comment "application/webapp-service has no endpoints" -m addrtype --dst-type LOCAL -m tcp --dport 30081 -j REJECT --reject-with icmp-port-unreachable
-A KUBE-SERVICES -d 10.107.112.64/32 -p tcp -m comment --comment "application/webapp-service has no endpoints" -m tcp --dport 8080 -j REJECT --reject-with icmp-port-unreachable
```
