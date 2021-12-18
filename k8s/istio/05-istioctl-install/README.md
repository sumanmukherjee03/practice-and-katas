## istioctl

Easiest way to get started with istio is via istioctl.
"x" refers to experimental features of istioctl.
```
brew install istioctl
istioctl x precheck
istioctl install
kubectl get pods -n istio-system
```
After installation via this method you are gonna find that only 2 pods are running in the `istio-system` namespace.
One pod is for istiod and the other one is for istio-ingressgateway.
We wouldnt be seeing the other components like kiali, grafana, jaeger etc.

There are ways to uninstall istio if it was installed via istioctl.
```
istioctl experimental
istioctl x uninstall --purge
```
This still retains the istio-system namespace btw.

Istio has builtin configuration profiles that come with various predefined configs like - default, demo, external, minimal etc.
The graphical components have been switched off in these profiles by default now anyways.
These components are rather considered addons to istio.
The `default` profile is the most common one to go with when installing into a production cluster.
That's the one which is installed when running `istioctl install`.

If you want to install the demo profile for example you can run this command.
```
istioctl install --set profile=demo
```
The command above can be used to switch profiles even if a profile was already installed.
The demo profile is unsuitable for production because it has a very high volume of tracing enabled. So, on a
high traffic system it can be very resource intensive.


### Integrations

There are several addons available for istio - kiali, grafana, jaeger, prometheus, cert-manager.
You can find the installation guides for these addons here - https://istio.io/latest/docs/ops/integrations/

You can find the manifests for addons installations if you follow the download link for istio - https://istio.io/latest/docs/setup/getting-started/#download
In the unpacked directory go to the path samples > addons which has the manifests for installing the addons.


### Profiles

You can dump istio profiles onto files like so to be able to tune it afterwards.
```
istioctl profile dump demo > raw-demo-profile.yaml
```

Remember yaml files of the kind `IstioOperator` are yaml files that can be applied via istioctl.
```
istioctl apply -f tuned-default-settings.yaml
```
But we can also run those files through the operator to get back yaml that can be interpreted by kubectl.

```
istioctl manifest generate -f tuned-default-settings.yaml > tuned-default-k8s.yaml
```
You can apply this generated yaml.
The only thing missing from the generated yaml is the definition of `istio-system`. But ofcourse we can manually add that in.
Another gotcha from the generated file is all the custom resource definiiton (`kind: CustomResourceDefinition`) need to be applied first.
So, it is best to separate those out into it's own file. Then you can run the whole manifest with kubectl.


For the addons, you can change the service types to be NodePort instead of ClusterIP so that you can use them directly via the node ip.
Ofcourse that is a personal choice. For our minikube installation you can always tunnel it via minikube service command.
