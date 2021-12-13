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
will have the same `app` label, there is no distinction in the App graph view. kiali puts a special meaning for the `app` label.

Now if you look at the "Versioned app graph" view then you would notice that there is a box around the 2 different deployments that
we are using to test out the canary. Adding a `version` label around the deployment will help with displaying that in the kiali UI.

So, kiali puts special meaning to the `app` and the `version` labels.

Here we are referring to `template > metadata > labels`.

The box around the workloads in the versioned graph actually represents the corresponding app.
And that box is clickable. You can right click on that box and go to "Show Details" to see the details of the App.
Where you can click on the service link to go to the service definition in kiali and click into "Actions" to create weighted rules.
Once you create the Destination Rules it will create the VirtualService and DestinationRules resources.
You can click on the triangle service icon to see how the traffic is split into multiple versions.

So, under the hood, the virtual service and destination rule that istio generates looks similar to this
```
apiVersion: networking.istio.io/v1beta1
kind: VirtualService
metadata:
  name: fleetman-staff-virtual-service
  namespace: default
spec:
  hosts:
    # This the name of the k8s service that istio is applying routing rules to.
    # And we are using the fully qualified in-cluster name because the k8s service could have been in a different namespace
    # So, for the proxy to be able to discover the pods, it's best to use the full name of the k8s service.
    # You could have also just used `fleetman-staff-service`
    - fleetman-staff-service.default.svc.cluster.local
  http:
  - route:
    # We use the full service name again in the actual routes.
    # The reason why we have the service repeated here is because it is possible
    # that traffic intended for fleetman-staff-service.default.svc.cluster.local
    # could be routed to 2 completely different services like
    #   - fleetman-staff-service-1.default.svc.cluster.local
    #   - fleetman-staff-service-2.default.svc.cluster.local
    - destination:
        host: fleetman-staff-service.default.svc.cluster.local # This is the target DNS name
        subset: risky-destination # This is the name of the subset in destination rules
      weight: 10
    - destination:
        host: fleetman-staff-service.default.svc.cluster.local # This is the target DNS name
        subset: safe-destination # This is the name of the subset in destination rules
      weight: 90

---
apiVersion: networking.istio.io/v1beta1
kind: DestinationRule
metadata:
  name: fleetman-staff-destination-rules
  namespace: default
spec:
  # This the name of the k8s service that istio is applying routing rules to.
  # And we are using the fully qualified in-cluster name because the k8s service could have been in a different namespace
  # So, for the proxy to be able to discover the pods, it's best to use the full name of the k8s service.
  # You could have also just used `fleetman-staff-service`.
  host: fleetman-staff-service.default.svc.cluster.local
  subsets:
  # The labels in the subset are selector labels
  - labels:
      version: risky # This means the select the pod with the label version and value risky
    name: risky-destination
  - labels:
      version: safe # This means the select the pod with the label version and value risky
    name: safe-destination
```

Now this configuration is added in the file `6-istio-rules.yaml`. So, go to the service details in the UI
and delete whatever Destination Rules you had created and apply the yaml config in the file.

When we apply the yaml for a VirtualService and DestinationRule the pilot, ie the controlplane component istiod
uses this configuration to change/update the proxies.The pilot will modify the proxies dynamically
with the custom routing rules. So, under the hood, really it is Envoy which is doing all this Canary stuff.
Remember a k8s service is in no way replaced or enhanced by a VirtualService. We still need the k8s service
for DNS lookup and get the IP addresses of the pods that serve the service. VirtualService on the otherhand
is really a proxy (Envoy) level component.

So, an outgoing request from webapp to the staff service for example will go through the webapp proxy first,
which will do service discovery using k8s service objects and then based on the VirtualService weighted configuration
in the webapp proxy will direct traffic to one of the pods.

If in the virtual service yaml we entered a wrong service name we would be getting 503 gateway errors because the proxy
wont be able to find the service to send traffic to.
And this is actually rightly reflected in kiali. kiali checks the host for the virtual service and destination rules.


### Stickiness

One common situation that accompanies Canary Releases is session stickiness. Subsequent requests from the same client
should keep getting the same version of the software, otherwise it can turn into a very frustrating situation in real life.
That's where session stickiness comes in.

DestinationRules define policies to traffic intended for a service after routing has occurred.
These can contain rules that specify configuration for LoadBalancing, connection pool size, remove unhealthy hosts from connection pool etc.

Stickiness is based on consistent hashing by the loadbalancer in istio proxy. The consistent hash is produced
based on the User cookie as the hash keythat specify configuration for LoadBalancing, connection pool size, remove unhealthy hosts from connection pool etc.

Stickiness is based on consistent hashing by the loadbalancer in istio proxy. The consistent hash is produced
based on several types of inputs coming from the client as the hash key. If the input is same, then the hash output is the same
and traffic consistently goes to the same subset. The various things that can be used as inputs for consistent hashing are
source ip, a http cookie, a http header value http query params etc.

Unfortunately though the session stickiness and weighted destination rules dont mix together and work.
There is an open istio issue related to this and envoy didnt seem to support it either for the longest time but seems to have fixed now.
However, it is still not supported by istio.

So, essentially, weighted subset based canary releases and http session stickiness doesnt work together in istio yet.

However, if we remove the weighted subsets and just have one subset, ie 1 destination, we can have session stickiness.
This is essentially because the weighted rule applies before the consistent hash based stickiness rule applies.
So, the pod the traffic is supposed to go to is already chosen.

So, this would not work :

                                 |---hash 1---->pod 1
           |----10% subset 1---->|
           |                     |---hash 2---->pod 2
           |
client ----|
           |
           |                     |---hash 1---->pod 1
           |----90% subset 2---->|
                                 |---hash 2---->pod 2


But, this would work :

                       |---hash 1---->pod 1
client ---1 subset---->|
                       |---hash 2---->pod 2

However, the downside is that you cant control how much of the traffic you want to divert to these canary pods.

One way to test this thing is by changing the input to the hashing algo to a http header instead of the source IP.
And we can test it with a curl request : `curl -H "x-test-canary: test1" http://localhost:<port-of-fleetman-webapp-service-exposed>/api/vehicles/driver/City%20Truck`
Ofcourse you have to make a change in the `DestinationRules` configuration to perform consistent hashing based on that http header name.
However, doing just that is not gonna fix our problems because we have multiple microservices calling each other
and unless this `x-` header is getting propagated across microservices, the fleetman-staff-service pods
arent gonna receive service with this header in the request and hence the session stickiness wouldnt work as expected.
So, we also need to make sure that the jaeger client library configures the propagation of these `x-` headers in the microservices.

Our microservices are already forwarding the `x-` header, so session affinity would work in this case.

Consistent hashing and session affinity can come in handy for performance enhancement, for example with something like caching.
