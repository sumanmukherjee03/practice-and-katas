### Minikube

To get Istio up and running we are gonna want to use a kubernetes cluster provisioned with minikube.
After provisioning the kubernetes cluster we need to run some quick tests to make sure that everything is ok.

```
brew install minikube
minikube start --memory 7168 --cpus 4 --feature-gates=EphemeralContainers=true
minikube ip
kubectl get pods --all-namespaces
minikube logs
minikube ssh
minikube status
minikube dashboard
kubectl create deployment hello-minikube --image=k8s.gcr.io/echoserver:1.4
kubectl expose deployment hello-minikube --type=NodePort --port=8080
kubectl port-forward service/hello-minikube 8080:8080
minikube service hello-minikube
```

Once you enable the ingress and ingress-dns addons for minikube, make sure to update the resolver for DNS in OSX
Follow the link here to get a better understanding of how to get ingress dns working
  - https://minikube.sigs.k8s.io/docs/handbook/addons/ingress-dns/
```
minikube addons enable ingress
minikube addons enable ingress-dns
kubectl apply -f https://raw.githubusercontent.com/kubernetes/minikube/master/deploy/addons/ingress-dns/example/example.yaml
sudo killall -HUP mDNSResponder
scutil --dns
```

If you want your docker client to talk to the minikube docker daemon there is an easy way to do that.
```
minikube docker-env
eval $(minikube -p minikube docker-env)
docker image ls
docker ps -a
```
