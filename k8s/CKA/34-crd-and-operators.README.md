## Custom Resource Definition

You have resources in k8s and the state of those resources are persisted in etcd. Controllers watch for changes to these resources
and communicate with the api server to persist the state of these resources in etcd.

You can have custom objects in kubernetes as well.

For example an object like `graphedge.yml`

```
apiVersion: graph.org/v1alpha1
kind: GraghEdge
metadata:
  name: edge1
  namespace: kruskals-demo
spec:
  src: node1
  dest: node2
  weight: 3
  directed: false
```

This resource GraphEdge will then get persisted in etcd. But for that ofcourse you will need a `GraphEdgeController`.

That will allow you to perform operations like these with kubectl
```
kubectl apply -f graphedge.yml
kubectl -n kruskals-demo get graphedge edge1
kubectl -n kruskals-demo delete graphedge edge1
```

However, besides the controller GraphEdgeController you will also need a CRD (custom resource definition) for the `GraphEdge` resource.
For example `graphedge-crd.yml` could look like this:
```
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: graphedge.graph.org
spec:
  scope: Namespaced
  group: graph.org
  names:
    kind: GraphEdge
    singular: graphedge
    plural: graphedges
    shortnames:
      - ge
  versions:
    - name: v1alpha1
      served: true
      storage: true
  schema:
    openAPIV3schema:
      type: object
      properties:
        spec:
          type: object
          properties:
            src:
              type: string
            dest:
              type: string
            weight:
              type: integer
              minimum: 0
              maximum: 10
            directed:
              type: bool
```
Remember that in the versions here, only 1 version can have `storage: true`.
Then you can create the CRD and subsequently query/update the custom resources via kubectl like so
```
kubectl -n kruskals-demo get ge
kubectl api-resources
```

-----------------------------------------------------------------------------------------------------------------------

### operators

We have an `EtcdCluster` CRD and a `EtCDController` which looks for changes to this CRD.
Similarly there is also a CRD `EtcdBackup` and `EtcdRestore`. There are operators which operate on these CRDs and
backup and restore the etcd cluster. So, operators are essentially more than just controllers.
They can manage the entire lifecycle of an application, like initializing, taking backups, restoring etc.

Lots of prebuilt operators are available in operatorhub.io .
