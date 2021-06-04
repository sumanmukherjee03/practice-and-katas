## Debugging

----------------------
### Debugging pods

curl http://web-service-ip:node-port
kubectl describe svc web-service
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

To check if a service is resolving DNS from inside a pod, try to get into a busybox pod and run this inside
  nslookup <service-name>.default.svc.cluster.local
  Depending on the kubelets --cluster-domain flag the root domain may be cluster.local or something else.

Or to test the same thing above try this from a node
  nslookup <service-name>.default.svc.cluster.local <kube-dns-service-ip>

If you can resolve FQDN but not the short names of the services
  Check /etc/resolv.conf inside the simple busybox pod to see of the entries for `search` are there or not
  and if the nameserver is pointing to the kube-dns service or not.

To check if the service is actually targeting pods properly
  kubectl get endpoints <service-name>

If all checks pass so far in debugging a service, then kube-proxy might be having some issues
  ps auxw | grep kube-proxy

Check the logs of kube-proxy to see what mode kube-proxy is running as iptables/ipvs etc.
  cat /var/log/kube-proxy.log | grep -i proxier
  OR
  journalctl -u kube-proxy | grep -i proxier

If kube-proxy is running in iptables mode, then you should be able to see the iptable entries for <service-name>
  iptables-save | grep -i <service-name>
    For each port of each Service, there should be 1 rule in KUBE-SERVICES and one KUBE-SVC-<hash> chain.
    For each Pod endpoint, there should be a small number of rules in that KUBE-SVC-<hash> and
    one KUBE-SEP-<hash> chain with a small number of rules in it.

Curl the service IP from one of the nodes
  curl <service-ip>:80



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


ON the worker nodes check
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
The proxy directive in the Corefile configuration is what tells where to forward out of cluster DNS queries.
So, that's an important pointer when pods are having DNS issues.

Check the kubelets `-cni-bin-dir` and `-network-plugin` runtime options if the coredns pod is in a pending state.
This usually means the kubelet is having issues placing pods.

If the coredns pods are in a crash loop
Your node could be running an older version of docker or might have SELinux enabled.
Try to upgrade docker on the node.

You could have SELinux enabled.
To disable SELinux on the node try these commands for an Ubuntu machine
```
# To get the status of SELinux
sestatus

# Also to get the status of SELinux - can be enforcing, permissive or disabled
# The permissive mode only logs but doesnt enforce any SELinux policy
getenforce

# This will change the SELinux mode to permissive, but this change does not persist over a reboot
setenforce 0
```

For disabling SELinux and persisting it across a reboot, amke the change in `/etc/selinux/config` with `SELINUX=disabled` directive.

```
kubectl -n kube-system get deployment coredns -o yaml | sed 's/allowPrivilegeEscalation: false/allowPrivilegeEscalation: true/g' | kubectl apply -f -
```
