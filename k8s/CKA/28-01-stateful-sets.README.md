### stateful sets

Statefulsets are similar to deployments with some differences.
With statefulsets for example pods are created in sequence. 1 pod comes up first and is in a running state,
then the next one comes up. This guarantees an order of deployment which is useful for applications
that require persistence - for example a persistent rabbitmq cluster for example.

Pods also get a unique name based on the index of the order in which they should get created.
Meaning pod names will be like rabbitmq-0, rabbitmq-1, rabbitmq-2 etc. Thus DNS addresses can also remain the same,
and helps with replication and/or disaster recovery. This makes it reliable for something like mysql master/slave deployment as well.
So, in essence stateful sets maintain a sticky identity.

Here's an example of a stateful set.
It is important to remember that deployment of stateful sets are ordered and graceful, ie one after the other.
In the Service object definition, notice how clusterIP is set to None. That's what creates a headless service for us.

```
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: google-storage
provisioner: kubernetes.io/gce-pd
parameters:
  type: pd-ssd
  replication-type: none
reclaimPolicy: Retain
allowVolumeExpansion: true
volumeBindingMode: Immediate

---
apiVersion: v1
kind: Service
metadata:
  name: mysql
  labels:
    app: mysql
spec:
  ports:
  - port: 3306
    name: db
  clusterIP: None
  selector:
    app: mysql

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: db
spec:
  serviceName: "mysql"
  replicas: 3
  selector:
    matchLabels:
      app: mysql
  template:
    metadata:
      labels:
        app: mysql
    spec:
      containers:
      - name: mysql
        image: mysql
        ports:
        - containerPort: 3306
          name: db
        volumeMounts:
        - name: data
          mountPath: /var/data/mysql
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes:
        - ReadWriteOnce
      storageClassName: google-storage
      resources:
        requests:
          storage: 1Gi
```

You can also scale stateful sets in an ordered fashion.
```
kubectl scale statefulset mysql --replicas=5
```
When you scale down the oldest one is deleted last and the newest one goes down first.

You can delete a stateful in a similar way to the above as well
`kubectl delete statefulset mysql`

You can change the `podManagementPolicy` field in stateful sets to `parallel` to not follow an ordered deployment approach.

In a master slave topology, for example in mysql, you wanna point the writes to the master and reads can happen over all nodes.
So, having a simple service for mysql that loadbalances across all pods is not gonna address this issue.
It's for cases like these that you can use a pods DNS. But that DNS is based on IP, so it can change.
That's where a headless service comes in.
A headless service just creates a DNS entry for a pod name and a known subdomain. It does not perform any loadbalancing for example.

so, if you create a headless service with the name `mysql-h`, then the corresponding pods get DNS entries like
<pod-name.headless-service-name.namespace.svc.cluster-domain>.

So, it could look like `mysql-0.mysql-h.default.svc.cluster.local`
A headless service definition is the same as any other service. But make sure to set `ClusterIP: None` in the yaml.

Dont forget to mention `serviceName: mysql` in the StatefulSet definition.
The StatefulSet uses this `serviceName` to create the DNS entries for the pods.



---------------------------------------

Persistent volume has a 1:1 relation to a PersistentVolumeClaim, which then has a 1:1 relation to a pod definition.
So, this is very static in nature and is not suitable for a deployment for example.
One improvement with StorageClass is that it makes the creation of the volume lazy. But still, the
volume claim will have a 1:1 relation with the pod. So, you still cant use that in a deployment.
Otherwise, all pods will use the same volume claim, as in they will be sharing the same storage. May be,
sometimes that's the desirable thing to do. But often it is not.

This creates the need for a separate `pvc` for each `pod`.
This is where the `volumeClaimTemplates` field comes into play. It allows creation of a pvc per pod in the deployment.

So, as the first pod is created, it provisions a pvc, which provisions a pv internally for the storage class and a storage is created in gcp
and finally that gets bound to the first pod. Then the same happens for the 2nd pod and so on.

When pods gets deleted, the pvc is not automatically deleted. But when a new pod takes it's place, the pvc is reattached to the
new pod. This enables stable storage for pods.
