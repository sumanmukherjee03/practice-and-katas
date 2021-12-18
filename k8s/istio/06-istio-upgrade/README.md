## Istio upgrade

There are a few different ways of upgrading an istio installation

### In place upgrade

Official documentation of istio upgrade is 1 step minor version upgrade. Meaning you can only upgrade 1 minor version at a time.

To get an older version of istio
```
curl -L https://istio.io/downloadIstio | ISTIO_VERSION=1.11.0 TARGET_ARCH=x86_64 sh -
cd ~/Downloads/istio-1.11.0 && cp bin/istioctl ~/bin/istioctl1_11_0
istioctl1_11_0 x precheck
istioctl1_11_0 install --set profile=demo
```

Once you have deployed the applications and generated some traffic head over to grafana dashboard and look into the
Mesh Dashboard. Scroll to the very bottom and there is a graph for Istio Version which shows how many pilot components
and how many proxy components exist for a specific version in the kubernetes cluster right now.
As application pods start up, the sidecar numbers go up.
If you have deployed an istio gateway then that accounts for 1 extra istio proxy.

When you perform an in-place upgrade, make sure to use the istioctl binary of the new version.
```
istioctl upgrade
watch kubectl get pods -n istio-system
```

With the in-place upgrade we might notice that if we are using an egressgateway then that doesnt upgrade.
So, we are left with the control plane in an upgraded state and the ingress gateway in an upgraded state. But the
application pod proxies remain in old state and the egressgateway does not seem to get upgraded.

So, for the application pods we need to be able to restart them for the new proxy to be picked up.
This command that do that all in one go
```
kubectl -n default rollout restart deploy
```

Except for the egressgateway the rest of the in-place upgrade went fairly well.

### Staged upgrade or Canary upgrade

For staged upgrade when installing the older version of istio add a tag called revision like so
```
istioctl1_11_0 install --set profile=demo --revision=1-11-0
```

Remember that doing this will actually label the pod of istiod. Also, the pod name will have it's version in the pod name unlike before.
However, we need to label the default namespace differently. The old label of `istio-injection: enabled` isnt
going to work for injecting sidecar proxies into the pods.

if you had accidentally labeled the default namespace you can remove it like this
```
kubectl label namespace default istio-injection-
```

This time the label we want to use on the namespace is `istio.io/rev: 1-11-0`
This is important for the sidecar proxy to be injected.
Once the application pods come up you will see that the proxies are connected to the istiod daemon of 1.11.0.

To check that, run the command
```
istioctl1_11_0 proxy-status
```
This command tells us which proxy is connected to which istiod daemon


Now for the rolling upgrade install the new version of istio. Again with the correct revision set.
```
istioctl install --set profile=demo --set revision=1-12-1
```
This prevents the old proxies to connect to the new istiod.

Now, if you run `istioctl proxy-status` you will notice that the only proxies that picked up the new version
are the ingressgateway and egressgateway in the istio-system namespace. All the application pods are still connected to
the old proxy.

You will notice that there are 2 mutatingwebhookconfigurations as well. This is what does the proxy injection into the pods.
```
kubectl get mutatingwebhookconfigurations
```

At this point upgrade the `istio.io/rev` tag to 1-12-1 and reapply the yaml file.
Now you can do a rolling restart of all pods or 1 pod at a time.
```
kubectl -n default rollout restart deploy
```

Run `istioctl proxy-status` to check if the proxies are talking to the new istiod
if all is good at this point we would want to remove the old istio version. Be sure not to use the --purge flag
and also be sure to mention the revision flag and of course use the correct binary.
```
istioctl1_11_0 x uninstall --revision=1-11-0
```

The only concern here is that the ingressgateway gets upgraded instantly when installing the new version of istio.
