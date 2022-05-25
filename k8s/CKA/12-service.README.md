## service

Services help connecting applications with other services or end users.

Service is also an object just like a deployment or replicaset etc.

Imagine that the host of a kubernetes worker is at 192.168.1.2
And imagine that a pod is on an internal network of 10.244.0.0/24, may be with an ip of 10.244.0.2.
Obviously a machine on the same network as the worker cant directly access the pod ip. This is where a service comes in.

One purpose of a service is to listen to a port on a node and forward the request onto a port on the pod, to which a container is listening to.
This is an example of a nodeport service

3 types of services
  - NodePort : exposes an internal port of a pod via a port on the host
  - ClusterIp : creates a virtual ip inside the cluster to enable communication between the services, for example a frontend pod to a backend pod
  - LoadBalancer : creates a loadbalancer in supported cloud providers to balance the load to your pods for an application


-----------------------------------------------

### NodePort

The service is like a virtual server inside the node. It is listening to a port of 80.
And forwarding to the target port on the pod - 80.
Inside the cluster the service gets it's own IP address.
This is called the cluster ip of the service.
And finally there is the port on the host to access the service externally which is known as the node port, ie 30008.
NodePorts can be in a valid range of 30000 - 32767

    ____________________________________________________
    |                                                   |
    |                                                   |
    |          ________________         ___________     |
    |          |  Service     |       |           |     |
    |          |              |       |   Pod     |     |
    |NodePort  |Port          |       |           |     |
    |30008     |80            |       |TargetPort |     |
    |          |              |       |80         |     |
    |          ________________       |___________|     |
    |            10.106.1.12            10.244.0.2      |
    |                                                   |
    |                                                   |
    |                                                   |
    |                                                   |
    |                                                   |
    ____________________________________________________
                      192.168.1.2

So, 192.168.1.2:30008 --> 10.106.1.12:80 --> 10.244.0.2:80
You can curl the host on the node port - `curl http://192.168.1.2"30008`.

For sample node port service definition - `cat nodeport-service-definition.yaml`
```
apiVersion: v1
kind: Service
metadata:
  name: frontend-service
spec:
  type: NodePort
  selector:
    app: nginx
    type: frontend
  ports:
    - targetPort: 80
      port: 80
      nodePort: 30008
```
If you dont provide a `nodePort`, then a port in the valid range on the node is automatically allocated.
If you dont provide `targetPort` it is assumed to be the same as `port`.

The selector section is important because that implies which pods to target in the service based on the labels in the pod definition template.

You can have multiple port mappings inside a service.

`kubectl create -f nodeport-service-definition.yaml`

To view the services
`kubectl get services`

Services use `Algorithm: Random` to forward requests to pods.
Also, services can have session affinity via the attribute `SessionAffinity: Yes`.

A service spanning multiple nodes in a multi node cluster have the NodePort exposed on the same port in all hosts.
So, all hosts expose 30008 and forward requests to the cluster ip of the service which in turn forwards it to the pods.

To expose a deployment using nodeport, start by creating a service template with some details filled in
```
kubectl expose deployment webapp-deployment --name=webapp-service --target-port=8080 --port=8080 --type=NodePort --dry-run=client -o yaml
```
Then edit the yaml file to add a known node port. There is no cli option to pass the desired node port.

OR

If you dont care about the port of the node port and a random one would work just fine, then
```
kubectl expose deployment webapp-deployment --name=webapp-service --target-port=8080 --port=8080 --type=NodePort
```

OR

In this case the service does not accept any selectors, so you will have to edit the document.
But at least you can specify the nodeport unlike the one above.
```
kubectl create service nodeport webapp-service --tcp=8080:8080 --node-port=30080 --dry-run=client -o yaml
```


--------------------------------------------------------------------------------


### ClusterIP

Pods can go up and down, meaning their IPs can keep changing. As such, we need a more static IP within the cluster
so that frontend pods for example can reliably communicate with the backend pods. This is where the service with ClusterIP comes in.
Also, services handle balancing load.

Each service gets a name and IP assigned to it and other pods must use the name to access the pods behind the service.

For sample node port service definition - `cat clusterip-service-definition.yaml`
```
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  type: ClusterIP
  selector:
    app: golang
    type: backend
  ports:
    - targetPort: 8080
      port: 8080
```
The type by default is `ClusterIP`.
If you dont provide `targetPort` it is assumed to be the same as `port`.
The `targetPort` is the port that the backend pod is listening to.

You can have multiple port mappings inside a service.

`kubectl create -f clusterip-service-definition.yaml`

To view the services
`kubectl get services`



--------------------------------------------------------------------------------


### loadBalancer

NodePort is great for exposing a service externally but it is not a convenient way to use the apps for end users,
because you need the ip of all the hosts and the 30000+ port number to access the apps. This inconvenience is solved
by using a loadbalancer type service. Kubernetes has capability to create loadbalancers in azure, gcp and aws.

If you dont have a cloud which is supported by kubernetes, it is analogous to create a separate VM with something
like nginx configured to serve as the load balancer.

In a single node cluster, and when there is no cloud, the LoadBalancer type service is the same as a NodePort type service.
