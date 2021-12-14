### Istio Gateways

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

