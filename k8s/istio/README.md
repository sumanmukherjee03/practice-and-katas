### What is kubernetes not good at?

Kubernetes is not good at inspecting interpod communication.
There is no visibility or control on inter pod comm. This is what a service mesh would help solve.

All network calls are gonna be routed through the service mesh layer.
The service mesh could have logic for pre routing logic or post response logic.

In case of Istio, it adds a proxy container inside of each pod as a sidecar to the application container.
The collection of all these proxy pods are collectively called "Istio Data Plane"

There is also a namespace called istio-system where "Istio Control Plane" components run.
Previously istio used to have multiple pods in this namespace each doing it's own thing.
However, there has been some major refactors and now we have a single pod that runs in that namespace called istiod (istio daemon).
There are some other pods in the control plane like Kiali UI etc.

Proxy pods call into the istiod pod to report metrics for telemetry and other sorts of functionalities.

### Getting Istio running

