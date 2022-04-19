## kube-scheduler

The kube-scheduler only decides which pod goes to which node. It does not actually run the pod in the host. That's the job of the kubelet.

The scheduler picks a node to place the pod on depending on a 2-step process.
  - filter out nodes that do not fit the requirement
  - rank nodes based on a prioritization algorithm and pick the node with the best score
    - for resources, the scoring is based on how much resources would be left on the host if the pod is scheduled there
      - the more resources left after pod placement, the better the score of that node


If we are not using kubeadm, kube-scheduler is available as a binary to download
```
wget https://storage.googleapis.com/kubernetes-release/release/v1.20.0/bin/linux/amd64/kube-scheduler
```

OR view the running process with
`ps -aux | grep kube-scheduler`

When running as a service the kube-scheduler is passed an option `--config` which points to the file containing the configuration.
`cat /etc/kubernetes/config/kube-scheduler.yaml`

if you have set it up with kubeadm tool, then kube-scheduler is available as a pod in the kube-system namespace.
`cat /etc/kubernetes/manifests/kube-scheduler.yaml`


--------------------------------------------------------------------------------

Kubernetes allows you to have multiple schedulers but they need to be of different names.
Also, if you are running the scheduler in HA mode, we need to ensure that the
corresponding lock object used in leader election by the scheduler is also named differently.
The reason for needing a lock object is because only one scheduler in HA cluster can be active at a time and others on standby.
Otherwise, multiple schedulers may try and place a new pod in multiple different nodes and cause more than the desired number of pods to be started.
To be able to do this the leader election process uses a lock object, the default lock object being `leases`.

If using the kubeadm tool, the manifest for the default scheduler from `/etc/kubernetes/manifests/kube-scheduler.yaml`
can be copied over to something like `/etc/kubernetes/manifests/custom-scheduler.yaml` and updated with customer scheduler
image and options etc.

Relevant sections from `/etc/kubernetes/manifests/custom-scheduler.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: custom-scheduler
  namespace: kube-system
spec:
  containers:
    - command:
      - kube-scheduler
      - --address=127.0.0.1
      - --kubeconfig=/etc/kubernetes/custom-scheduler.conf
      - --leader-elect=true                                   # This is required if the master is a HA cluster
      - --leader-elect-resource-lock=leases
      - --scheduler-name=custom-scheduler
      - --lock-object-name=custom-scheduler                   # This is DEPRECATED in favour of --leader-elect-resource-name. The default value for the flag is "kube-scheduler".
.....
      image: custom-scheduler
      name: custom-scheduler
```
The `--leader-elect-resource-lock=leases` option which is the default, determines what kind of object is used for locking during leader election.
Possible options are - leases, endpoints, configmaps, endpointleases, configmapleases.

Use the kubectl create command to create this new scheduler in the kube-system namespace.

To use this new scheduler in the pod definition yaml, pass the `schedulerName: custom-scheduler` setting
in the pod definition spec section.

pod-definition.yaml
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: nginx-container
      image: nginx
  schedulerName: custom-scheduler
```

To see which scheduler picked up the pod, view the events in the current namespace
`kubectl get events | grep -i "scheduled"`

-------------------------------------------------------------------------------

Relevant sections of a running kube-scheduler
```
spec:
  containers:
  - command:
    - kube-scheduler
    - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
    - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
    - --bind-address=127.0.0.1
    - --kubeconfig=/etc/kubernetes/scheduler.conf
    - --leader-elect=true
    - --port=0
    image: k8s.gcr.io/kube-scheduler:v1.19.0
```

And a complete kubernetes manifest for custom scheduler. All the necessary resources are created in the `kube-system` namespace.
```
apiVersion: v1
kind: ServiceAccount
metadata:
  name: my-scheduler
  namespace: kube-system

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: my-scheduler-as-kube-scheduler
subjects:
- kind: ServiceAccount
  name: my-scheduler
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: system:kube-scheduler                 # We have bound the system cluster role for kube-scheduler to this custom service account so that it gets the same privileges as the kube-scheduler
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: my-scheduler-as-volume-scheduler
subjects:
- kind: ServiceAccount
  name: my-scheduler
  namespace: kube-system
roleRef:
  kind: ClusterRole
  name: system:volume-scheduler               # We have bound the system cluster role for kube-scheduler to this custom service account so that it gets the same privileges as the kube-scheduler
  apiGroup: rbac.authorization.k8s.io

---
apiVersion: v1
kind: Pod
metadata:
  creationTimestamp: null
  labels:
    component: my-scheduler
    tier: control-plane
  name: my-scheduler
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-scheduler
    - --authentication-kubeconfig=/etc/kubernetes/scheduler.conf
    - --authorization-kubeconfig=/etc/kubernetes/scheduler.conf
    - --bind-address=127.0.0.1
    - --kubeconfig=/etc/kubernetes/scheduler.conf
    - --leader-elect=false
    - --port=10459
    - --scheduler-name=my-scheduler    # This is matched with 'spec.schedulerName' on pod spec, to determine which scheduler is responsible for scheduling a pod
    image: k8s.gcr.io/kube-scheduler:v1.19.0
    imagePullPolicy: IfNotPresent
    livenessProbe:
      failureThreshold: 8
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10459
        scheme: HTTP
      initialDelaySeconds: 10 # This tells the kubelet to wait for 10 seconds before starting it's first liveness probe
      periodSeconds: 5
      timeoutSeconds: 7
    name: kube-scheduler
    resources:
      requests:
        cpu: 100m
    # Define a startup probe for containers that are slow to initially start up
    startupProbe:
      failureThreshold: 5
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10459
        scheme: HTTP
      initialDelaySeconds: 10 # This tells the kubelet to wait for 10 seconds before starting it's first startup probe
      periodSeconds: 5
      timeoutSeconds: 7
    volumeMounts:
    - mountPath: /etc/kubernetes/scheduler.conf
      name: kubeconfig
      readOnly: true
  hostNetwork: true
  priorityClassName: system-node-critical
  volumes:
  - hostPath:
      path: /etc/kubernetes/scheduler.conf
      type: FileOrCreate
    name: kubeconfig
```
To get the custom scheduler fully working, we finally need to get the cluster role `system:kube-scheduler` updated
such that it allows the `resourceNames` with `my-scheduler` for the resources `endpoints` and `leases`.
To do this run `kubectl -n kube-system edit system:kube-scheduler` and update the cluster role.

Here's an example of how to write a custom scheduler - https://kubernetes.io/blog/2017/03/advanced-scheduling-in-kubernetes/
