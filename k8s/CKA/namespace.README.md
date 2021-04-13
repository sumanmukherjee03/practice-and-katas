## namespace

Three namespaces are created by default - the `default` namespace, the `kube-system` namespace and the `kube-public` namespace.

The kube-system namespace is used by kubernetes for it's own resources while the default namespace is generally used by users.

kube-public is the namespace where resources that should be made available to all the users is created.

You can have multiple namespaces in the cluster.
To reach db pod in the same namespace you can reach via the dns `db-service`.
However, to reach a DB service in another namespace, say in database namespace, you can reach via
`db-service.database.svc.cluster.local`.

In this DNS above, the `cluster.local` is the default domain name of the kubernetes cluster.

```
name-of-service.namespace.subdomain-of-services.domain-of-cluster
  db-service    database     svc                  cluster.local
```

To create pods in a namespace
```
kubectl get pods -n kube-system
kubectl create -f pod-definition.yaml -n kube-system
```
OR move the namespace into the pod definition file under metadata section
```
metadata
  name: nginx-pod
  namespace: kube-system
```

This is an example of `namespace-definition.yaml`
```
apiVersion: v1
kind: Namespace
metadata:
  name: dev
```

To create the namespace
`kubectl create -f namespace-definition.yaml`
OR
`kubectl create namespace dev`

To look at the namespaces
```
kubectl get ns
kubectl get ns --no-headers
```

To permanently switch the context of the client in kubernetes, use this command
`kubectl config set-context $(kubectl config current-context) --namespace=dev`
After this point you will be able to use the dev namespace without passing the namespace flag.

To get pods in all namespaces
```
kubectl get pods --all-namespaces
```

To limit resources in a namespace, a ResourceQuota object can be created.
compute-quota-definition.yaml
```
apiVersion: v1
kind: ResourceQuota
metadata:
  name: compute-quota
  namespace: dev
spec:
  hard:
    pods: 10
    requests.cpu: 4
    requests.memory: 5Gi
    limits.cpu: 10
    limits.memory: 10Gi
```
To create the resource quota run `kubectl create -f compute-quota.yaml`
