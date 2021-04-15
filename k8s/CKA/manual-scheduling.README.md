### manual-scheduling

To manually schedule a pod, add a `nodeName` field in the pod definition yaml file.

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
  nodeName: worker-01
```

Without a scheduler, the easiest way is to assign a `nodeName` to the pod at creation time.
But if you want to move the pod to a different node, you wont be allowed to do that because
kubernetes does not allow you to modify the nodeName property.


-------------------------------------------------------------------------------

Another way to schedule pod on a node without using the scheduler is to create a pod binding object.

pod-binding.yaml
```
apiVersion: v1
kind: Binding
metadata:
  name: nginx
target:
  apiVersion: v1
  kind: Node
  name: worker-02
```

Send a json request to the kubernetes api-server with the binding object definition in json format
```
curl -H 'Content-Type:application/json' -X POST -d '{"apiVersion": "v1"...}' http://$SERVER/api/v1/namespace/default/pods/$PODNAME/binding/
```
