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

  journalctl -u kube-apiserver


ON the worker nodes check
  service kubelet status
  service kube-proxy status
