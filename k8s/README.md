### Docker
If you are on mac and using the docker-desktop app, make sure you go into settings
and allocate 5 GB of memory, 3 GB swap and 4 CPUs for minikube to work well.
Also, this might be a good to cleanup and prune unused resources being hogged by docker.
```
docker system prune
docker volume prune
docker image prune
```

### Minikube

To get Istio up and running we are gonna want to use a kubernetes cluster provisioned with minikube.
After provisioning the kubernetes cluster we need to run some quick tests to make sure that everything is ok.

```
brew install minikube
minikube start --memory 4096 --feature-gates=EphemeralContainers=true
minikube ip
minikube status
minikube logs
minikube ssh
kubectl get pods --all-namespaces
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

To cleanup unused resources from minikube
```
minikube ssh -- docker system prune
minikube ssh -- docker volume prune
minikube ssh -- docker image prune
```

To stop and remove minikube
```
minikube delete
```
