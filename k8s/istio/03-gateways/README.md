### Istio Gateways

NOTE : In this set of examples `*-istio-rules-*.yaml` are mutually exclusive.
       You can only have one of them at a time. So, make sure to delete the resources from the other one if you are about to apply one.

In the versioned app graph, if you wanna see the version app view along with the service view,
then check "Service Nodes" in the "Display" dropdown menu.

If you try to get the canary working via the virtual services and destination rules only, without an istio gateway,
then the 90-10 split would not work. That is because the fleetman-webapp service is a frontend service.
The traffic splitting on the internal services is done by the proxy component on outgoing connections,
but in this case the connection to fleetman-webapp service is coming from an external client and that part of
the traffic is not going through any proxy. That's the reason we need a gateway to perform the traffic split, similar to a reverse proxy.

Envoy solves this problem via a "Edge proxy". Istio ingress-gateway does something very similar.
It is a pre-provided container run by istio in the `istio-system` namespace.
All incoming traffic from external sources come through this pod.

```
kubectl get pods -n istio-system | grep ingressgateway
kubectl get svc -n istio-system | grep ingressgateway
```

It is worth noting that the `istio-ingressgateway` service is of type `LoadBalancer`, meaning if you deploy this
in a cloud provider it'll provision a resource similar to a ALB in aws and would have an external IP, ie an ip outside the CIDR of the kubernetes cluster.

The ingress gateway istio configuration targets pods in the `istio-system` namespace that have labels `istio=ingressgateway`.
```
kubectl get pods -n istio-system --show-labels | grep ingressgateway
```

Once you have the ingressgateway setup, if you tunnel through minikube, you should be able to reach the node port
that exposes port 80 of the ingressgateway service in istio-system.
But if there are no paths and forwarding rules configured in the virtualservice you will be getting a 404.
This is because the ingressgateway pod is just running the proxy and traffic comes from outside to that proxy container
but when exiting the proxy it doesnt know where to go. And all these proxies are configured with the VirtualService
and DestinationRule via the istiod (pilot) component.

Remember that the same VirtualService configuration is all proxies, ie proxy of the internal pods as well as the one in the ingressgateway.
So, when configuring hosts for the VirtualService, you should consider hosts for internal and external traffic,
ie, other pods will try to reach the fleetman webapp via `fleetman-webapp.default.svc.cluster.local` internally
and browsers will try to reach it externally through it's public domain.
If you have put asterix in the hosts for external traffic that cancels out any other host you have put into place.

So, once you have configured and gotten ingressgateway to work, it is time to remove the `NodePort` on the service definition of fleetman-webapp service,
and change it to `ClusterIP`.


In matching rules for virtualservice, the matches can be based on many different things like uri prefix, headers, query params, scheme (http|https), http methods etc.

If you want to deploy 2 versions of a software, it is best to always avoid url based prefixes because that
involves change of application code in multiple places. It is quite troublesome. Instead the easiest way to do this is using subdomains.
That involves multiple subdomains, one for each version of the app to route traffic from external sources into the proper version.
For internal traffic you might use a different virtual services altogether, but with routing based on a header.
And that `x-` header would get added at the entrypoint application in the code based on whether it was the experimental version or original version
and then that header would get propagated all the way down to the rest of the traffic in the cluster.

We can achieve Dark releases with matching rules in the virtualservice via headers. Ofcourse you have to make sure that
the header gets propagated all the way down in all services for that to work. For the external traffic coming from outside
the cluster, we can use curl to pass headers or install a browser extension like ModHeader to add a custom header to all our
requests from the browser.
