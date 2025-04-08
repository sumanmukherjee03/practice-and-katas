## Local Development Guide For Contribution
The local setup is not awfully complicated if you have `podman` installed.
Assuming that you are using podman for local development in mac, this section will help you with testing things locally.
You can install podman with homebrew.

```sh
wget https://github.com/kubernetes-sigs/cloud-provider-kind/releases/download/v0.3.0/cloud-provider-kind_0.3.0_darwin_arm64.tar.gz
tar xvf cloud-provider-kind_0.3.0_darwin_arm64.tar.gz
mv cloud-provider-kind ~/bin/
rm LICENSE
rm cloud-provider-kind_0.3.0_darwin_arm64.tar.gz
```

Setup podman locally on your macbook so that you can run the experiments
```sh
podman machine stop kindexp-cluster
podman machine rm kindexp-cluster
rm -rf ~/.config/containers
podman machine init --cpus 4 --memory=6144 kindexp-cluster
podman machine list
podman machine start kindexp-cluster; sleep 30; podman machine stop kindexp-cluster
podman machine set --rootful kindexp-cluster
podman machine start kindexp-cluster;
podman version
podman run quay.io/podman/hello
```

First things first, you need 2 kind clusters for testing things locally.
```sh
helm repo add cilium https://helm.cilium.io/
kind create cluster --name kindexp1 --config $(pwd)/kindexp1-cilium.yaml --image kindest/node:v1.28.12; sleep 30
helm install cilium cilium/cilium --version 1.14.13 --namespace kube-system --set image.pullPolicy=IfNotPresent --set ipam.mode=kubernetes --set nodeinit.enabled=true --set kubeProxyReplacement=partial --set hostServices.enabled=false --set externalIPs.enabled=true --set nodePort.enabled=true --set hostPort.enabled=true --set cluster.name=kind-kindexp1 --set cluster.id=1
kubectl get nodes --context kind-kindexp1
kubectl --context kind-kindexp1 label node kindexp1-control-plane node.kubernetes.io/exclude-from-external-load-balancers-
kubectl --context kind-kindexp1 apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml
kubectl --context kind-kindexp1 apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/kind/deploy.yaml
```

Now create the second kind cluster.
```sh
helm repo add cilium https://helm.cilium.io/
kind create cluster --name kindexp2 --config $(pwd)/kindexp2-cilium.yaml --image kindest/node:v1.28.12; sleep 30
helm install cilium cilium/cilium --version 1.14.13 --namespace kube-system --set image.pullPolicy=IfNotPresent --set ipam.mode=kubernetes --set nodeinit.enabled=true --set kubeProxyReplacement=partial --set hostServices.enabled=false --set externalIPs.enabled=true --set nodePort.enabled=true --set hostPort.enabled=true --set cluster.name=kind-kindexp2 --set cluster.id=2
kubectl get nodes --context kind-kindexp2
kubectl --context kind-kindexp2 label node kindexp2-control-plane node.kubernetes.io/exclude-from-external-load-balancers-
kubectl --context kind-kindexp2 apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.15.3/cert-manager.yaml
kubectl --context kind-kindexp2 apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.10.1/deploy/static/provider/kind/deploy.yaml
```

```sh
cilium clustermesh enable --context kind-kindexp1 --service-type LoadBalancer
cilium clustermesh status --context kind-kindexp1 --wait
cilium clustermesh enable --context kind-kindexp2 --service-type LoadBalancer
cilium clustermesh status --context kind-kindexp2 --wait
cilium clustermesh connect --context kind-kindexp1 --destination-context kind-kindexp2
cilium clustermesh status --context kind-kindexp1 --wait
cilium clustermesh status --context kind-kindexp2 --wait
cilium connectivity test --context kind-kindexp1 --multi-cluster kind-kindexp2
kubectl delete ns cilium-test
```

```sh
kubectl --context kind-kindexp1 -n default apply -f https://raw.githubusercontent.com/cilium/cilium/v1.14.3/examples/kubernetes/clustermesh/global-service-example/cluster1.yaml
kubectl --context kind-kindexp2 -n default apply -f https://raw.githubusercontent.com/cilium/cilium/v1.14.3/examples/kubernetes/clustermesh/global-service-example/cluster2.yaml
kubectl --context kind-kindexp1 -n default get deployment,pods,svc,cm
kubectl --context kind-kindexp2 -n default get deployment,pods,svc,cm
kubectl --context kind-kindexp1 -n default run --restart Never --rm -it --image ubuntu:22.04 sumantest -- /bin/bash -c 'apt-get -qq update; apt-get -qq -y install curl wget dnsutils telnet >/dev/null; for i in $(seq 1 30); do curl http://rebel-base/; done;'
kubectl --context kind-kindexp2 -n default run --restart Never --rm -it --image ubuntu:22.04 sumantest -- /bin/bash -c 'apt-get -qq update; apt-get -qq -y install curl wget dnsutils telnet >/dev/null; for i in $(seq 1 30); do curl http://rebel-base/; done;'
kubectl --context kind-kindexp1 -n default annotate svc rebel-base io.cilium/service-affinity="local"
kubectl --context kind-kindexp1 -n default run --restart Never --rm -it --image ubuntu:22.04 sumantest -- /bin/bash -c 'apt-get -qq update; apt-get -qq -y install curl wget dnsutils telnet >/dev/null; for i in $(seq 1 30); do curl http://rebel-base/; done;'
kubectl --context kind-kindexp1 -n default annotate svc rebel-base io.cilium/service-affinity-
kubectl --context kind-kindexp1 -n default annotate svc rebel-base io.cilium/service-affinity="remote"
kubectl --context kind-kindexp1 -n default run --restart Never --rm -it --image ubuntu:22.04 sumantest -- /bin/bash -c 'apt-get -qq update; apt-get -qq -y install curl wget dnsutils telnet >/dev/null; for i in $(seq 1 30); do curl http://rebel-base/; done;'
```

```sh
kubectl --context kind-kindexp1 -n default apply -f rebel-ingress.yaml
for i in $(seq 1 30); do curl http://127.0.0.1:18080/; done;
```

### Debugging podman for local development on mac
Sometimes when you are using podman locally for development purposes on your mac, it can give you grief.
These commands can help ease the pain
```sh
podman machine stop kindexp-cluster
podman machine start kindexp-cluster
podman ps -a
podman start kindexp1-control-plane
podman start kindexp1-worker
podman start kindexp1-worker2
podman start kindexp1-worker3
podman start kindexp2-control-plane
podman start kindexp2-worker
podman start kindexp2-worker2
podman start kindexp2-worker3
kubectl get nodes
kubectl --context kind-kindexp1 -n kube-system get pods
kubectl --context kind-kindexp1 -n bitgo-mongo-operator-system delete validatingwebhookconfiguration bitgo-mongo-operator-validating-webhook-configuration
kubectl --context kind-kindexp1 -n kube-system rollout restart daemonset/cilium
kubectl --context kind-kindexp1 -n kube-system rollout restart daemonset/kube-proxy
kubectl --context kind-kindexp1 -n kube-system rollout restart deployment/coredns
kubectl --context kind-kindexp1 -n kube-system rollout restart deployment/cilium-operator
kubectl --context kind-kindexp1 -n cert-manager rollout restart deployment/cert-manager
kubectl --context kind-kindexp1 -n cert-manager rollout restart deployment/cert-manager-webhook
kubectl --context kind-kindexp1 -n cert-manager rollout restart deployment/cert-manager-cainjector
kubectl --context kind-kindexp1 -n local-path-storage rollout restart deployment/local-path-provisioner
```
