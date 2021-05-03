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
