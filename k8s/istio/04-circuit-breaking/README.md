## Circuit breaker

We can have circuit breaking behavior with istio. If a microservice we are making calls to gets overloaded and is experiencing
trouble responding to our requests, the circuit breaker comes into play and stops relaying the outbound network calls.
When the upstream microservice gets back on it's feet, the circuit breaker starts relaying those calls again.
This prevents the applications from encountering cascading failures. Cascading failures are caused by one overloaded
service causing other services which rely on it to get overloaded and eventually starting to fail.

Istio does this by removing a bad pod from the loadbalacing pool of upstream of backends from a proxy.
If you got only 1 replica of a pod and that's the one removed by the circuit breaker, then your downstream client service
should be able to handle some backpressure.

Something to remember is that you need a virtualservice when you want to configure some sort of routing
and you need a destinationrule when you need to configure the loadbalancer. One can exist without the other depending on your need.

Circuit breaker counts or maintains monitoring metrics to each individual pod. However, the configuration of
it via destinationrule applies to the service object.

Something to remember is that in older versions of istio circuit breaker is triggered by HTTP 502, 503, 504.
The newer versions of istio however have added a new field to let you specify any 500 error to trigger the circuit breaker.

When the circuit breaker triggers, a pod isnt evicted from the loadbalancer forever. It is evicted for a specified amount
of time only and then it is added back to the pool. However the ejection time follows an exponential backoff pattern,
ie if first time it got ejected for 30s, 2nd time it'll be for 60s, 3rd time it'll be for 90s and so on.



## Mutual TLS

A service container can call another service container with http, but since the calls go through
a proxy, the proxy upgrades the call to https. The receiving side then downgrades it again and forwards the http call to the receiving service container.

In earlier versions of istio, mutual tls used to be implemented by a pod named `istio-citadel`, but now that
component is baked into istiod.

So, the traffic flow inside the cluster would look somewhat like this.
  ---->container            container
          |                     |
         HTTP                 HTTP
          |                     |
         proxy-----mTLS-------proxy

Also, this is not just applicable to http, this is true for protocols like rmq or grpc as well.
mTLS means that the caller verifies the callee and the callee verifies the caller, both verifications are done with
certificates issued by citadel. The certs are issued and distributed to the proxies by citadel.

istio can block all non TLS traffic OR istio can upgrade non TLS traffic to mTLS.

In the newer versions of istio, mTLS is automatically switched on. So, there's nothing to do.


### Permissive mTLS

You can have a namespace where you have not setup the sidecar proxy injection, ie kubernetes pods without proxies.
When these pods talk to pods that are using istio, then a regular HTTP call cant be upgraded to HTTPS because
the sending side cant present a certificate.

Istio by default uses permissive mTLS, meaning if the sender doesnt present a cert then the proxy
allows the connection to proceed as it is.

One way to test permissive mTLS is to edit a service, in this case say the fleetman-position-tracker service and
change the service type to `NodePort` and give it a `nodePort` of say, 32323. Then tunnel to that nodeport via minikube.
Now, if you curl an endpoint on that service with http (not https) it returns back a valid response. This is good because
the sending side, ie our curl is not sending any certs for the client to validate and also, the request is not being upgraded
to https in any way because the sending side does not go through a proxy.

If you dont like the permissive mTLS behavior, ie, we want to block HTTP connections, then we can
update the permissive mTLS to strict mTLS.

To turn on strict mode in mTLS
```
kubectl apply -f - <<EOF
apiVersion: security.istio.io/v1beta1
kind: PeerAuthentication
metadata:
  name: "default"
  namespace: "istio-system"
spec:
  mtls:
    mode: STRICT
EOF
```

On turning the mTLS to strict mode, now all the curl reuqests will be failing because on the client side it is not going through the proxy.

Ofcourse this is mutual TLS, so, a pod with proxy if makes a call to a pod without a proxy running in some other namespace
will also fail to make a http call if the mTLS mode is strict
