The newer versions of istio do not require us to have a basic auth with username and password for kiali.
Kiali would generally be behind an ingress gateway which would handle authentication for us.
For the purposes of this lesson we can skip applying the 3-kiali-secret.yaml

If you have already created pods and havent applied the proper label to the default namespace,
then the envoy proxy container wont be injected into the pods.
In that case after you apply the proper label to the namespace, it is best to delete the existing pods so
that the new pods pick up the envoy containers.
`kubectl delete pods --all`

If you wanna access kiali over minikube
```
minikube service kiali -n istio-system
```

It is important to remember that the Graph view of the Kiali dashboard for a specific namespace and service
is dynamic in nature. If there is no traffic between 2 services the graph edges become grey and ultimately disappear.

A workload represents all the pods backing a service.

Right click on the service triangle of a service in the graph view and look at the details of a service.
You can view the Inbound metrics of the service as well in that view.

Similarly in the workloads view you can right click on a workload icon, ie the circle and go to it's details view.
This will show you that the workload is a deployment and it's inbound metrics etc.


### Peek into traffic management

```
kubectl get virtualservices
kubectl get vs

kubectl get destinationrules
kubectl get dr
```

If one service is misbehaving for example, you can suspend traffic from the UI with kiali by going to the details
view of a service and clicking on the actions button.

Once you suspend traffic, go back to the cli and try the commands above again and you are gonna see some entries.

If you are satisfied with the experiment and want to turn traffic back on, go to the Actions button in the
details view of the service and "Delete destination rules".
