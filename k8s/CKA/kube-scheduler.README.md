## kube-scheduler

The scheduler only decides which pod goes where. It doesnt actually run the pod in the host. That's the job of the kubelet.

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
      - --leader-elect=true
      - --scheduler-name=custom-scheduler
      - --lock-object-name=custom-scheduler
.....
      image: custom-scheduler
      name: custom-scheduler
```

Use the kubectl create command to create this new scheduler in the kube-system namespace.

To use it in the pod definition yaml pass the `schedulerName: custom-scheduler` setting
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

And a complete kubernetes manifest for custom scheduler
```
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
    - --scheduler-name=my-scheduler
    image: k8s.gcr.io/kube-scheduler:v1.19.0
    imagePullPolicy: IfNotPresent
    livenessProbe:
      failureThreshold: 8
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10459
        scheme: HTTP
      initialDelaySeconds: 10
      periodSeconds: 5
      timeoutSeconds: 7
    name: kube-scheduler
    resources:
      requests:
        cpu: 100m
    startupProbe:
      failureThreshold: 5
      httpGet:
        host: 127.0.0.1
        path: /healthz
        port: 10459
        scheme: HTTP
      initialDelaySeconds: 10
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

Here's an example of how to write a custom scheduler - https://kubernetes.io/blog/2017/03/advanced-scheduling-in-kubernetes/
