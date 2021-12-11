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


### Peek into distributed tracing

`minikube service tracing -n istio-system &`

On the jaeger UI you can choose "Custom time" and select a custom time range to get a subset of traces.
You can also check the checkboxes beside the collection of spans to compare traces.
That might for example give you a view of where the requests are going separate ways - for example it might
point to the obvious that your requests are going separate paths from the api-gateway onwards.

Remember in the span view, you are gonna see 2 entries for each k8s service. That's because the requests also go through the proxies.
Just something to keep in mind so that you dont get confused.

Jaeger adds the header `x-request-id` for tracing. So, we need to enable propagation of the headers in our services,
ie propagate the trace context.
In particular, istio relies on propagation of B3 trace headers and Envoy-generated request ID.
B3 propagation is the propagation of headers that start with `b3` or `x-b3-`. Here's a bunch of headers
that are needed to be propagated for istio tracing to work properly. The header propagation can be done with
jaeger or zip client libraries.
  - x-request-id
  - x-b3-traceid
  - x-b3-spanid
  - x-b3-parentspanid
  - x-b3-sampled
  - x-b3-sampled
  - x-b3-flags
  - b3
  - x-ot-span-context (if you wanna use opentracing)

Here's a link - https://istio.io/latest/about/faq/distributed-tracing/#how-to-support-tracing
to point to how to setup distributed tracing for the applications.
Here's another link to depict how to setup tracing context propagation - https://istio.io/latest/docs/tasks/observability/distributed-tracing/overview/#trace-context-propagation


### Monitoring

`minikube service grafana -n istio-system &`

2 dashboards that are specifically of interest to us in grafana are
  - Istio Service Dashboard
  - Istio Workload Dashboard


--------------------------------------------------------------------------------

## Traffic management

Istio has a few traffic management features. Canary releases is one of them.

### Canary release

One way of thinking about canary in kubernetes is to have a pod with a different container deployed as part of a service,
i.e., create a separate deployment object with a different name but use the same label on the deployment and the pods.
And of course change the container image to be the one with the experimental tag.
The round robin load balancing will do the rest. This will effectively drive 33% of the traffic to this experimental pod.
But it is hard to control the percent of requests that can go through the canary.
If you want 1 canary pod to receive only 10% of traffic you have to increase the total number of normal pods 10
so that the load balancing only forwards 10% of traffic to the canary pod. This is an expensive way to have canaries in kubernetes.

Once you have configured 2 different deployments with 2 different images for testing out a canary release,
now if you go to the graph UI in kiali and look at the workloads graph, you will notice traffic flowing to 2 workloads.
What is referred to as workload in Kiali, can be thought of as a deployment in k8s.
However, if you look at the App graph, you wont notice 2 different entities representing the canary. That is because
what istio calls an app is based on the `app` label in your k8s manifests. And because both the deployments, ie the normal and the canary
will have the same `app` label, there is no distinction in the App graph view.


