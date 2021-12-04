## What is kubernetes not good at?

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

## Getting Istio running

A pre-requisite to get istio working is to have a kubernetes cluster. Try the minikube installation and setup documentation in the parent dir.
To create istio resources follow these steps
```
kubectl apply -f 1-istio-init.yaml
kubectl apply -f 1-istio-minikube.yaml
```

For injecting the proxy containers into the data plane, we can use istioctl which will be able to inject
a sidecar into our yaml definition.
Or we can enable a flag named istio-injection in  a kubernetes resource like a namespace for example via labels
so that the injection of a sidecar proxy container happens automatically when a yaml is applied, by the admission controller
when pods are created in that namespace.
Basically istiod takes care of that proxy container injection based on the label value of `istio-injection` in the namespace.
Without the label on the namespace we arent gonna see any data from the data plane. So, we wont be seeing any of the diagnostics.
```
kubectl describe ns default
kubectl label ns default istio-injection=enabled
kubectl get ns default -o yaml
```

Now start the application
`kubectl apply -f 4-application-full-stack.yaml`
If you do not see 2 containers initializing in the pods, then we probably didnt label the namespace properly.

To diagnose issues with the app above, we can use kiali, which is a tool provided by istio.
```
kubectl get svc -n istio-system
minikube service kiali -n istio-system
```
Now you can go to the kiali UI

You can also open up the jaeger UI for tracing service calls with this command
```
minikube service tracing -n istio-system
```

From the example in the 4-application-full-stack.yaml the calls from the staff service to the external
service `fleetman-driver-monitoring`, the calls are taking too long (20s - 30s) or are failing.
This is a good place to have short circuit in place, or drop outgoing requests from the `staff-service` for
the `fleetman-driver-monitoring`. Or add a timeout for those requests, so that we dont have to wait for the application code to be modified.
So, we can modify the application yaml and modify the `fleetman-driver-monitoring` VirtualService and add a timeout
```
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: fleetman-driver-monitoring
spec:
  hosts:
  - 2oujlno5e4.execute-api.us-east-1.amazonaws.com
  http:
  - match:
      - port: 80
    route:
      - destination:
          host: 2oujlno5e4.execute-api.us-east-1.amazonaws.com
          port:
            number: 443
    timeout: 1s
```
