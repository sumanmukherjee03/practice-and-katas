## kubeadm

We are assuming that you already have 3 VMs ready with ubuntu bionic base image.
The very bare minimum required for the VMs is 2 GB RAM, 2 CPUs, Linux (Debian/Rhel), Unique hostnames and Swap disabled (for kubelet to work).
The sample cluster we will be creating here consists of a single master and 2 workers.
The nodes have been setup in the `192.168.56.0/24` CIDR range.

### Pre-requisites before creating cluster with kubeadm

In your VM switch to the root user to install and configure all the pre-requisites.

This is a set of instructions provided in this documentation for setting up a kubernetes cluster via kubeadm
  - https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/

1. To ensure that iptables on the nodes can see bridged traffic we must check that `br_netfilter` module is loaded by the kernel.
```
lsmod | grep 'br_netfilter'
modprobe --dry-run br_netfilter
modprobe overlay
modprobe br_netfilter
```

2. We must set the following to 1 in `sysctl` config.
sysctl is the linux utility that allows for setting and getting attributes for the kernel like system limits, security settings etc.
```
cat <<EOF | tee /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
net.ipv4.ip_forward = 1
EOF
sysctl --system
```
Another way to update the required sysctl settings is via
```
sysctl net.bridge.bridge-nf-call-iptables = 1
sysctl net.bridge.bridge-nf-call-ip6tables = 1
sysctl net.ipv4.ip_forward = 1
```

3. Install docker container runtime
```
apt-get update; apt-get install -y apt-transport-https ca-certificates curl software-properties-common gnupg2
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
apt-get update; apt-get install -y containerd.io docker-ce docker-ce-cli
```

4. Setup the configuration for the docker daemon (See the link here - https://kubernetes.io/docs/setup/production-environment/container-runtimes/#docker)
In particular configure the docker daemon to use systemd for the management of container's cgroups.
In linux systems where systemd is the init system for the processes, then systemd acts as the cgroup manager.
So, it is prudent to use the same cgroup manager for the docker daemon. You could use something like `cgroupfs` for the cgroup manager of kubelet and container runtime.
But that can cause 2 panes of system resources and the system can become unstable under pressure.
That's why it is best to choose systemd as the cgroup manager for the container runtime.
```
cat <<EOF | sudo tee /etc/docker/daemon.json
{
"exec-opts": ["native.cgroupdriver=systemd"],
"log-driver": "json-file",
"log-opts": {
  "max-size": "100m"
},
"storage-driver": "overlay2"
}
EOF
```

5. Reload the systemd daemon so that it picks up the new docker daemon config and then Restart the docker daemon once the configuration is set
```
systemctl enable docker
systemctl daemon-reload
systemctl restart docker
systemctl status docker
```

6. Add the kubernetes apt repository for installing kubelet, kubeadm and kubectl
```
curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | tee /etc/apt/sources.list.d/kubernetes.list
apt-get update
```

7. Install kubelet, kubeadm and kubectl and mark them as hold to prevent automatic upgrades in ubuntu machines.
By default ubuntu systems upgrade packages automatically and that's a problem for kubernetes. So, this step is necessary.
```
apt-get install -y kubelet kubeadm kubectl
apt-mark hold kubelet kubeadm kubectl
```


### Setting up the cluster using kubeadm

1. The cluster components need to be initialized on the `master` node.
Some important things to consider are the pod networking CNI plugin - like weave or calico etc.
The pod networking CIDR should not coincide with the node CIDR.
We will also need the apiserver advertise url ie, the master node ip if it is a single master node cluster or the loadbalancer url if it is a HA cluster (ie multiple masters).
```
kubeadm init --pod-network-cidr 10.244.0.0/16 --apiserver-advertise-address=192.168.56.2
```

2. Switch to a non-root user in the master node, ie your normal user in the master node and copy over the kube config to be able to interact with the apiserver via kubectl.
```
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
kubectl get nodes
```
NOTE : Copy the command provided for `kubeadm join` and keep it somewhere or echo it on the master node so that you can retrieve it from history later.
This needs to be run from the worker nodes afterwards.


3. Install calico operator and CRDs as the pod networking CNI plugin (Documentation here - https://docs.projectcalico.org/getting-started/kubernetes/quickstart)
As a pre-requisite enable the kernel bridge module setting to let the bridge traffic go to iptable rules for processing - `sysctl net.bridge.bridge-nf-call-iptables = 1`
```
kubectl create -f https://docs.projectcalico.org/manifests/tigera-operator.yaml
```

4. Create the calico custom resources for pod networking. More details here - https://docs.projectcalico.org/getting-started/kubernetes/quickstart
Calico uses a different pod cidr by default. Make sure what you put in the calico custom-resources.yaml is what you have provided to kubeadm for pod network cidr.
```
cat << EOF | tee $HOME/calico-custom-resources.yaml
# This section includes base Calico installation configuration.
apiVersion: operator.tigera.io/v1
kind: Installation
metadata:
  name: default
spec:
  # Configures Calico networking.
  calicoNetwork:
    # Note: The ipPools section cannot be modified post-install.
    ipPools:
    - blockSize: 26
      cidr: 10.244.0.0/16
      encapsulation: VXLANCrossSubnet
      natOutgoing: Enabled
      nodeSelector: all()
EOF
```
To see the calico pods run `watch kubectl get pods -n calico-system`
Then run `kubectl create -f $HOME/calico-custom-resources.yaml`

5. Next go to a worker node and make that worker node join the master
```
kubeadm join 192.168.56.2:6443 --token st1zft.1eh6dk8ygvpx2eym --discovery-token-ca-cert-hash sha256:94d4d8184229faa660259cf35b294367c3659be3368f570ae8e8b6d12d135b41
```
How to list/create the kubeadm join token that can be used by the worker nodes to join a kubernetes cluster if you lost the one from before.
This needs to be run on the master.
```
kubeadm init
kubeadm token list
kubeadm token create --print-join-command
```

6. Verify that pods can be scheduled
kubectl get nodes
kubectl run nginx --image=nginx
kubectl -n kube-system get cm kubeadm-config -o yaml



An alternative to using calico networking is to use the weave net for pod networking.
Here's how to apply weave works pod network solution. More documentation is here - https://www.weave.works/docs/net/latest/kubernetes/kube-addon/
```
curl -fsSLo weave-daemonset.yaml "https://cloud.weave.works/k8s/net?k8s-version=$(kubectl version | base64 | tr -d '\n')"
kubectl apply -f weave-daemonset.yaml
```

Another alternative to using calico is to use flannel for pod networking.
However, flannel does not support NetworkPolicy in kubernetes.
```
kubectl apply -f https://raw.githubusercontent.com/coreos/flannel/2140ac876ef134e0ed5af15c65e414cf26827915/Documentation/kube-flannel.yml
```

NOTES
----------------
Some trivial commands for getting OS/network related information
```
cat /etc/*-release
lsb_release -cs
apt-cache search kubeadm
ifconfig
```
