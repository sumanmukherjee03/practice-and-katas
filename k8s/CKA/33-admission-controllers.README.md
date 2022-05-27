## Admission controllers

Admission controllers run after authentication and authorization in k8s.
Admission controllers can not only validate and reject api requests from the users but it can also mutate or change the requests made by users.
In general Mutating Admission Controllers are invoked before the Validating Admission Controllers,
so that whatever was added or removed by the mutation can also be validated before persisting on etcd.

Different types of admission controllers available are :
- AlwaysPullImages
- DefaultStorageClass
- EventRateLimit
- NamespaceExists
- ResourceQuota
- ValidatingAdmissionWebhook

..... etc

- NamespaceAutoProvision : This admission controller is not enabled by default


To view the list of allowed admission controller plugins or the ones enabled by default run
```
kubectl exec kube-apiserver-controlplane -n kube-system -- kube-apiserver -h | grep 'enable-admission-plugins'
```

To view the list of extra admission controller plugins enabled
```
grep enable-admission-plugins /etc/kubernetes/manifests/kube-apiserver.yaml
```

So, if you want to include a new admission controller for example, you can update the
`/etc/kubernetes/manifests/kube-apiserver.yaml`

```
apiVersion: v1
kind: Pod
metadata:
  annotations:
    kubeadm.kubernetes.io/kube-apiserver.advertise-address.endpoint: 10.28.225.3:6443
  creationTimestamp: null
  labels:
    component: kube-apiserver
    tier: control-plane
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:
  - command:
    - kube-apiserver
    - --advertise-address=10.28.225.3
    - --allow-privileged=true
    - --authorization-mode=Node,RBAC
    - --client-ca-file=/etc/kubernetes/pki/ca.crt
    - --etcd-cafile=/etc/kubernetes/pki/etcd/ca.crt
    - --etcd-certfile=/etc/kubernetes/pki/apiserver-etcd-client.crt
    - --etcd-keyfile=/etc/kubernetes/pki/apiserver-etcd-client.key
    - --etcd-servers=https://127.0.0.1:2379
    - --insecure-port=0
    - --secure-port=6443
    - --enable-bootstrap-token-auth=true

    ...

    - --enable-admission-plugins=NodeRestriction, NamespaceAutoProvision
    - --disable-admission-plugins=DefaultStorageClass

    ...
```

Btw, the NamespaceAutoProvision and NamespaceExists admissions are now replaced by the `NamespaceLifecycle` admission controller.
Besides handling namespace validations it ensures that the kube-system and kube-public namespaces are auto provisioned and cant be deleted.

Since the kube-apiserver runs as a pod in the controlplane, you can also check the process
`ps -ef | grep kube-apiserver | grep admission-plugins`

-------------------------------------------------------------------------------

To support custom Admission Controllers we have 2 other admission controllers available
- MutatingAdmissionWebhook
- ValidatingAdmissionWebhook

We can point these webhooks to point to an external server whether that is hosted within the kubernetes cluster or is external to the cluster.
Once the server which is running the Custom Admission Controller processes that request and responds back to the webhook request,
and based on that either some resource is modified with the response or an api request is accepted/rejected.
The webhook to the custom admission controller is invoked with an `AdmissionReview` object.

For example the request object to the custom admission controller would look like below.
It usually contains all the details of the type of api request made, the resource in question etc.
```
{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "request": {
    "uid": "......",
    "kind": {"group": "autoscaling", "version": "v1", "kind": "Scale"},
    "resource": {"group": "apps", "version": "v1", "resource": "deployments"},
    ....
  }
}
```

The custom admission controller then responds with a json response of object `AdmissionReview` as well.
```
{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "response": {
    "uid": "same uid as the request",
    "allowed": true
  }
}
```

The response of a mutating webhook request can look somewhat like this
```
{
  "apiVersion": "admission.k8s.io/v1",
  "kind": "AdmissionReview",
  "response": {
    "uid": "same uid as the request",
    "allowed": true
    "patch": "base64 encoded value of a patch object - similar to what we would provide for a kubectl command",
    "patchtype": "JSONPatch"
  }
}
```

If your webhook service is running within the k8s cluster then make sure to expose it with a service.

Once the webhook service is created, we need to create a validate webhook object with our api server.
There are 2 types of objects - ValidatingWebhookConfiguration or MutatingWebhookConfiguration.
This is an example of ValidatingWebhookConfiguration only.
Something to remember here is that the webhook configuration object itself is not namespaced. It is a cluster wide resource.
```
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: "pod-policy.example.com"
webhooks:
- name: "pod-policy.example.com"

  # The rules specify when to call the webhook server
  rules:
  - apiGroups:   [""]
    apiVersions: ["v1"]
    operations:  ["CREATE"]
    resources:   ["pods"]
    scope:       "Namespaced"

  # The clientConfig can also have a key url : <some_external_service_url>
  # Provide the caBundle for TLS
  clientConfig:
    service:
      namespace: "example-namespace"
      name: "example-service"
      path: "/validate"
    caBundle: "Ci0tLS0tQk...<`caBundle` is a PEM encoded CA bundle which will be used to validate the webhook's server certificate.>...tLS0K"
  admissionReviewVersions: ["v1", "v1beta1"]
  sideEffects: None
  timeoutSeconds: 5
```

To retrieve a webhook configuration object run :
`kubectl get MutatingWebhookConfiguration demo-webhook`

For webhooks that get deployed within the cluster or outside, generally the mode of communication will be
requiring a client cert to present to the webhook processing server and a key to encrypt the symetric encryption key.
And usually this will be stored in a k8s secret. Obviously, we would expect the webhook service to be exposing https port whether that is external or a k8s service.
Example command to do so
```
kubectl create secret tls webhook-server-tls --cert=/path-to/webhook-server-tls.crt --key=/path-to/webhook-server-tls.key
```
