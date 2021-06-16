## Advanced kubectl commands

Some sample jsonpath queries to help with querying the kubernetes api later on.

### jsonpath

JSONPath can be used to query data from json or yaml both.
Remember all results if a jsonpath query are wrapped within a pair of [].
```
$.vehicles.car.color                                # Use the dot notation to query dictionaries and $ represents the root element
$[0]                                                # Gets the first element of a list
$[-1]                                               # Gets the last element of a list. If that doesnt work try these - $[-1:0] OR $[-1:]
$[0,3]                                              # Gets you the first and fourth element in the list
$[0:3]                                              # Gets you the first till third element, ie not including the fourth element in the list
$.car.wheels[1].model
$[?(@ > 4)]                                         # ?() -> means a criteria. @ -> means each item in the list that is being iterated upon.
$.car.wheels[?(@.location == "rear-right")].model
$.*.color                                           # * -> means all elements. So here it means all objects in the root object
$[*].model                                          # Here * means all objects in the list
$.car.wheels[*].model
$.*.wheels[*].model
```
Operators can be @ == 4, @ != 4, @ in [1,2,3,4,5], @ nin [1,2,3,4,5] etc



### kubectl + jsonpath

Some sample jsonpath queries with kubectl
```
kubectl get nodes -o json
kubectl get pods -o json
kubectl get pods -o jsonpath='{.items[0].spec.containers[*].image}'
kubectl get pod nginx -o jsonpath='{$.spec.containers[*].image}'
kubectl get nodes -o jsonpath='{.items[*].metadata.name}'
kubectl get nodes -o jsonpath='{.items[*].status.nodeInfo.architecture}'
kubectl get nodes -o jsonpath='{.items[*].status.capacity.cpu}'
kubectl get nodes -o jsonpath='{.items[*].metadata.name}{"\n"}{.items[*].status.capacity.cpu}'
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.cpu}{"\n"}{end}'
kubectl get nodes -o custom-columns=NODE:.metadata.name,CPU:.status.capacity.cpu
kubectl get nodes --sort-by=.metadata.name
kubectl config view --kubeconfig=my-kube-config -o jsonpath='{range $.users[*]}{.name}{"\n"}{end}'
kubectl config view --kubeconfig=my-kube-config -o jsonpath='{$.contexts[?(@.context.user == "aws-user")].name}{"\n"}'
kubectl get persistentvolume --sort-by=.spec.capacity.storage -o custom-columns=NAME:.metadata.name,CAPACITY:.spec.capacity.storage
```
