## kube-proxy

kube-proxy is the component which is responsible for managing the network so that
every pod can talk to every other pod no matter which worker node they are running in.

Instead of ip for pod to pod communication, we expose a pod via a service.
The service object itself gets an ip too.
A pod can talk to another pod using the name of the service.

A service object however is not an actual container and does not have any network interface or actively listening process.
So, it cant join the pod network. It is a virtual component that lives in memory in the kubernetes cluster.
Then, how do pods talk to each other via a service?

That's where kube-proxy comes in. kube-proxy is a process that runs on each worker node in the kubernetes cluster.
Everytime a service object is created, the job of the kube-proxy is to create the appropriate iptable rules,
such that traffic can be forward from a pod on a node to other pods exposed via services.
It creates iptable rules to forward traffic directed towards the ip of the service to the actual pods that are backing the service.

So, if there's a pod with an ip of 10.32.0.15 exposed via a service with IP 10.96.0.12,
an iptable rule will be created to forward traffic going to 10.96.0.12 -> 10.32.0.15.

For installing kube-proxy on worker nodes download it from the kubernetes release page.
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kube-proxy
```

If kube-proxy is running as a linux service in each of the worker nodes, you can view it's status via the linux services.

kube-proxy.service definition is available in
`cat /etc/systemd/system/kube-proxy.service`

The path of the kube-proxy config is usually passed as an option to the kube-proxy binary.
`cat /var/lib/kube-proxy/kube-proxy-config.yaml`

OR

you can deploy it via the kubeadm tool.
If deployed it via the kubeadm tool, you can view the kube-proxy pods via

`kubectl get pods -n kube-system | grep kube-proxy`

When deployed via kubeadm, kube-proxy is in fact deployed as a daemonset so that it runs on every node.

`kubectl get daemonset kube-proxy -n kube-system`
