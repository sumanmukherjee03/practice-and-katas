### Circuit breaker
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
