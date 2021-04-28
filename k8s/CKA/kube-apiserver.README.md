## kube-apiserver

kubectl talks to the api server, which in turn talks to etcd and retrieves information to form responses.

```
kubectl get nodes
```

The api can be invoked with simple curl commands too.

For example - this is the equivalent of creating a pod
`curl -X /api/v1/namespaces/default/pods ...`

Steps involved :
  - kube api server authenticates the incoming request
  - kube api server validates the incoming request
  - kube api server creates a pod but does not assign it to any node
  - kube api server updates that information in the etcd cluster
  - the scheduler which continously monitors the api server notes that a new pod has been created
  - the scheduler identifies which node the pod can be placed in and communicates that information back to the api server
  - the api server then updates that information in the etcd cluster
  - the api server then passes that information to the kubelet in the appropriate node
  - the kubelet then creates the pod and instructs the container runtime to deploy the image
  - the kubelet passes the information of the state back to the api-server
  - the api server then updates the etcd cluster

If not using the kubeadm the kube-apiserver is available for download through the kubernetes release page on google
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kube-apiserver
```

Documentation for flags when starting the kube apiserver :
 - https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/

Some security related flags
```
  --etcd-ca-file=/var/lib/kubernetes/ca.pem
  --etcd-certfile=/var/lib/kubernetes/kubernetes.pem
  --etcd-keyfile=/var/lib/kubernetes/kubernetes-key.pem

  --kubelet-certificate-authority=/var/lib/kubernetes/ca.pem
  --kubelet-client-certificate=/var/lib/kubernetes/kubernetes.pem
  --kubelet-client-key=/var/lib/kubernetes/kubernetes-key.pem
  --kubelet-https=true
```

kubeadm deploys the kube-apiserver as a pod in the kube-system namespace.
The pod is `kube-apiserver-master`.
You can view the pod definition file located in
`cat /etc/kubernetes/manifests/kube-apiserver.yaml`.

In a non kubeadm setup the options can be visible from the kube-apiserver linux service
`cat /etc/systemd/system/kube-apiserver.service`
OR
`ps -aux | grep kube-apiserver`
