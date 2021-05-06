## authorization

There are several ways to define what an authenticated user/serviceaccount can do in a kubernetes cluster
  1. RBAC
  2. ABAC
  3. Node authorization
  4. Webhook Mode


### apigroups

These are the various resources and verbs on those resources that authorization is supposed to be protecting.

Some of these apis are mainly for cluster health and inspection. For example :
  - `/version` - api group for looking up the version of the kube-apiserver.
  - `/metrics`, `/healthz`, `/logs` - are some pretty self explanatory apis.

The other 2 main categories of apis that control functionality of the cluster itself are

`/api` - used for the `core` group
    - /v1
        - namespaces, pods, rc, events, endpoints, nodes, bindings, PV, PVC, configmaps, secrets, services

`/apis` - used for the `named` group
    - /apps
        - /v1
            - /deployments
                - verbs : list, get, create, delete, update, watch
            - /replicasets
                - verbs : list, get, create, delete, update, watch
            - /statefulsets
                - verbs : list, get, create, delete, update, watch
    - /extensions
    - /networking.k8s.io
        - /v1
            - /networkpolicies
                - verbs : list, get, create, delete, update, watch
    - /storage.k8s.io
    - /authentication.k8s.io
    - /certificates.k8s.io
        - /v1
            - /certificatesigningrequest
                - verbs : list, get, create, delete, update, watch


To list the available api groups you can also simply just hit the kube-apiserver using curl
  - `curl http://localhost:6443 -k --key admin.key --cert admin.crt --cacert ca.crt`
  - `curl http://localhost:6443/apis -k --key admin.key --cert admin.crt --cacert ca.crt | grep -i "name"`

Of course as you can see above it requires you to pass the key and certs to the curl command.
An alternative is to use kubernetes proxy command. This proxies the kube-apiserver to localhost:6443.
Also, it saves you from passing the authentication information in the curl commands.
```
kubectl proxy
curl http://localhost:6443 -k
curl http://localhost:6443/apis -k
```


### Node Authorizer

kubelets interact with the kube-apiserver to Read Services, Endpoints, Nodes, Pods etc and
also to Write Node Status, Pod Status, Events etc.
To be able to perform these actions the kube-apiserver authorizes the kubelet via the NodeAuthorizer mechanism.
In the kubelet client certificate, the cert should have the group specified via `/O=system:nodes` in the `-subj` parameter of the cert
and the cert name to be something like `CN=system:node:node01`.
Meaning any user which has a name prefix of `system:node:` and a group of `system:nodes` will be granted the privileges similar to a kubelet.


### ABAC

You can create a policy for an user like so in a policy file
```
{"kind": "Policy", "spec": {"user": "eng-user", "namespace": "*", "resource": "pods", "apiGroup": "*"}}
{"kind": "Policy", "spec": {"user": "qa-user", "namespace": "test", "resource": "pods", "apiGroup": "*"}}
```
and pass the policy file to the kube-apiserver at start for the policies to take effect.
This can be done for groups too, as in add another line with a policy to the policy file for a group.
However, everytime this policy file is changed, the kube-apiserver needs to be restarted. As such this is not the preferred solution for authorization.


### Webhook

If you want an external authorization mechanism, that's where the webhook mechanism of kubernetes authorization comes into play.
For example `eng-user` requests access for a resource to `kube-apiserver`. The `kube-apiserver` asks `Open Policy Agent` if the user can access a resource
and that response is returned back and the user is allowed access or NOT.


### AlwaysAllow

Allows all requests without performing any authorization checks

### AlwaysDeny

Denies all requests without performing any authorization checks

### RBAC

This is the most commonly used mode for authorization.
Create a role object `cat engineer-role.yaml`
```
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: engineer
rules:
  - apigroups: [""]
    resources: ["pods"]
    verbs: ["list","get","create","update","delete"]
  - apigroups: [""]
    resources: ["ConfigMaps"]
    verbs: ["list","get","create","delete"]
    resourceNames: ["backend-app", "frontend-app"]
```
For the core group leave the group section as an empty string. For other groups specify the group.
Roles and RoleBindings are namespace specific.
In the resources, you can add resourceNames field to further harden which specific resources to target for the authorization.
Create the role using `kubectl create -f engineer-role.yaml`.

To link a user to the role, create another object called RoleBinding.
`cat john-engineer-binding.yaml`
```
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: john-engineer-binding
subjects:
  - kind: User
    name: john
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: engineer
  apiGroup: rbac.authorization.k8s.io
```
Create the role binding using `kubectl create -f john-engineer-binding.yaml`

Helpful kubectl commands to inspect roles and role bindings.
```
kubectl create role engineer --verb=list --verb=create --verb=delete --resource=pods --dry-run=client -o yaml
kubectl get roles
kubectl get rolebindings
kubectl describe role engineer
kubectl create rolebinding engineering-user-binding --role=engineer --user=engineering-user --dry-run=client -o yaml
kubectl describe rolebinding engineering-user-binding
```

To check if the current user can access a certain action for a resource use this
```
kubectl auth can-i create deployments
kubectl auth can-i delete nodes
```

An admin can even impersonate another user
```
kubectl auth can-i create deployments --as john
kubectl auth can-i create pods --as john --namespace test
```


-------------------------------------------

The mode of authorization is set on the kube-apiserver as an option at start.
`--authorization-mode=AlwaysAllow`. The value is set to `AlwaysAllow` by default.

You can set multiple modes too like `--authorization-mode=Node,RBAC,Webhook`.
Authorization happens in the order in which it is specified when the kube-apiserver is started,
ie Node -> RBAC -> Webhook . If a module denies the request, it asks the next module in order.
As soon as a module allows a request, it breaks from the chain.
