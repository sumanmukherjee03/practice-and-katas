The newer versions of istio do not require us to have a basic auth with username and password for kiali.
Kiali would generally be behind an ingress gateway which would handle authentication for us.
For the purposes of this lesson we can skip applying the 3-kiali-secret.yaml

If you have already created pods and havent applied the proper label to the default namespace,
then the envoy proxy container wont be injected into the pods.
In that case after you apply the proper label to the namespace, it is best to delete the existing pods so
that the new pods pick up the envoy containers.
`kubectl delete pods --all`
