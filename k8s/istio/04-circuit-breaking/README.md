### Circuit breaker
We can have circuit breaking behavior with istio. If a microservice we are making calls to gets overloaded and is experiencing
trouble responding to our requests, the circuit breaker comes into play and stops relaying the outbound network calls.
When the upstream microservice gets back on it's feet, the circuit breaker starts relaying those calls again.
This prevents the applications from encountering cascading failures. Cascading failures are caused by one overloaded
service causing other services which rely on it to get overloaded and eventually starting to fail.

Istio does this by removing a bad pod from the loadbalacing pool of upstream of backends from a proxy.
If you got only 1 replica of a pod and that's the one removed by the circuit breaker, then your downstream client service
should be able to handle some backpressure.
