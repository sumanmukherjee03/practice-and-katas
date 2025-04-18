## NetworkPolicy

By default all pods can communicate with each other within the kubernetes cluster.
This is because on the overlay network the kubernetes cluster has a network policy of "All Allow".
If we want selective traffic between pods/services and want to restrict communication from any
pod within the cluster then Ingress and Egress rules come into play. NetworkPolicy is the kubernetes object to handle this.

This is an example where there is a setup like so, nginx(listening to port 80) -> webapp(listening to port 8080) -> db(listening to port 3306).
This is a sample network policy for the db server.
The network policy is applied to the pods matching label `name: db`
and the ingress rule allows traffic from pod matching label `name: webapp` to pods matching label `name:db`.
Example: `cat network-policy.yaml`
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: db-policy
spec:
  podSelector:
    matchLabels:
      name: db                    # These are the pods that the policy is attached to and is applicable
  policyTypes:
    - Ingress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              name: webapp        # These are the pods that the network policy is allowing traffic from
      ports:
        - protocols: TCP
          port: 3306
```
Subsequently, you need to run `kubectl create -f network-policy.yaml` for creating the network policy.

But network policies can be more detailed than the one above.

Something to remember in the network policy above is that it is allowing the webapp labeled pods in any namespace
to reach the db pod. Meaning, the network policy is not namespaced.

However, if you want to restrict only webapp pods from a certain namespace for instance, then the network
policy also needs to have a namespaceSelector , similar to the podSelector .
Of course the namespace should have this label set for the selector to work.
Important to remember that in the rule you can combine the 2 selector types to only allow dev webapp pods.

Another kind of supported selector is ipBlock. Imagine if there was an external data warehousing tool that
needs to read from the database and run some ETL jobs to put transformed data in a data warehouse.

Also, imagine a case where there is an agent in the DB server that pushes metrics to a Saas solution.
To allow this we can have a Egress rule as well.

Example `cat network-policy.yaml`
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: db-policy
spec:
  podSelector:
    matchLabels:
      name: db
  policyTypes:
    - Ingress
    - Egress
  ingress:
    - from:
        - podSelector:
            matchLabels:
              name: webapp
          namespaceSelector:
            matchLabels:
              environment: dev
        - ipBlock:
            cidr: [192.160.10.0/24]
            except:
              - [192.160.10.1/32]
      ports:
        - protocols: TCP
          port: 3306
  egress:
    - to:
        - ipBlock
            cidr: [196.120.20.7/32]
      ports:
        - protocols: TCP
          port: 32000
          endPort: 32010
```

If we dont specify any rules in a network policy then it is deemed as default deny.
This is an example policy to deny all egress traffic from pods labelled marketing.
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: deny-egress-marketing
  namespace: default
spec:
  podSelector:
    matchLabels:
      app: marketing
  policyTypes:
    - Egress
```

This however is a more aggressive network policy to default deny all traffic from all pods.
For security reasons this might be a good network policy to have in place in your cluster.
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: default
spec:
  podSelector: {}
  policyTypes:
    - Egress
    - Ingress
```
When you have a default deny all policy, DNS resolution from within pods wont work.
So, make sure you enable that when allowing traffic for the pods via network policies.

That's why a better `default-deny-all` policy is this.
Also, keep in mind that the default deny is for that namespace only. It is not cluster wide.
```
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: default-deny-all
  namespace: default
spec:
  podSelector: {}
  policyTypes:
    - Egress
    - Ingress
  egress:
    - ports:
        - port: 53
          protocol: TCP
        - port: 53
          protocol: UDP
```

Remember that network policies always use labels. Make sure to add labels to pods, deployments, namespaces etc.
This might mean that you will have to add labels to the default namespace as well when allowing incoming traffic from default namespace into another namespace.

----------------------------------------------------
Pod networking solutions that support network policies
  - kube-router
  - calico
  - romana
  - weave-net

Network solutions that do not support network policies
  - flannel

If a pod has multiple network policies then the UNION of all network policies is applied on the pod.
And the order of the network policies on the pod doesn't matter.

