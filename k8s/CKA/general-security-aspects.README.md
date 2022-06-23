## K8s security Starter Pack
These are some guidelines for kubernetes security best practices

### General Security Principles

These are some security principles. Although they are independent concepts, you can think of them as a chain when considering an external attacker

1. Limit the Attack surface
  - reduce unnecessary code
  - reduce unnecessary packages in containers or operating systems
  - reduce the number of open ports
2. Defense in Depth
  - multiple layers
  - redundancy is a good thing (unlike DRY for application development)
3. Least Privilege
  - provide pods, containers and apps only the access that is required to run them

### Various categories for considering a security architecture
1. Host operating system security
  - kubernetes nodes whether master or worker components should only be running kubernetes resources
  - reduce the attack surface by removing unnecessary applications
  - keep OS upgrades up to date
  - keep packages up to date
  - Use a Runtime security tool
  - Find and identify malicious processes
  - Restrict IAM and SSH access
2. Kubernetes cluster security which builds on top of the host OS security
  - kubernetes components like kubelet, kube-apiserver etc should be up to date
  - restrict external access
    - keep etcd server privately accessible only
    - keep network rules tight for kubernetes components to talk to each other and not have blanket rules
    - if possible, keep kube-apiserver privately accessible only
    - use authentication and authorization for access to the kube-apiserver
    - use admission controller to inspect every request coming to the kube-apiserver
      - you can have NodeRestriction admission controllers
      - you can also have Custom Policies via OPA (Open Policy Agent)
    - Enable audit logging
    - Use some Security Benchmarking tool like CIS to dive into how to securely run the different components
    - Encrypt etcd at rest
    - Encrypt traffic where ever possible, especially between kube-apiserver and etcd
3. Application security for the containers
  - Use secrets as opposed to hard coded credentials
  - Use RBAC
  - You can look at container sandboxing
  - Container hardening
    - reduce attack surface by not running as root
    - container can have a readonly file system if possible to make the container truly immutable at runtime
    - container image is minimal in nature and doesnt have unnecessary packages
  - Vulnerability scanning
    - container image scanning with tools like Claire and openscap
    - application code scanning
  - mTLS for container communications
    - service meshes make this very easy

### Requirements for a CNI plugin implementation
2 important requirement for a CNI plugin implementation are:
  - Every pod can communicate with every other pod, at least by default
  - For every pod to communicate with every other pod, even if the other pod is on a different host, the network should not require a NAT

--------------------------

### Listing the certs in kubernetes

```
cd /etc/kubernetes/pki
openssl x509 -in ca.crt -text
openssl x509 -in apiserver.crt -text
openssl x509 -in apiserver-kubelet-client.crt -text

cd /etc/kubernetes/pki/etcd
openssl x509 -in ca.crt -text
openssl x509 -in server.crt -text
openssl x509 -in peer.crt -text
openssl x509 -in healthcheck-client.crt -text
```

The certs outside the etcd directory are all signed by the kubernetes ca server, ie signed by `kubernetes`.
The apiserver server for example has subject name of `kube-apiserver` but it also has a whole bunch of alternate names for the cert as well.
The cert for the apiserver to communicate with the kubelet is `apiserver-kubelet-client.crt`.
This has a subject name of `kube-apiserver-kubelet-client` and uses the group `system:masters`.

The etcd certs are all signed by the `etcd-ca`.
The etcd server and peer certs have the node hostname and node address as the subject names and alternate names.
However the etcd healthcheck client cert has a subject name of `kube-etcd-healthcheck-client` and has group `system:masters`.


If you navigate to the directory `/etc/kubernetes`, you will see a few config files
  - kubelet.conf
  - scheduler.conf
  - controller-manager.conf

These are all kubeconfig files and have the CA cert embeded into the kubeconfigs.
The users in these kubeconfigs are different though
  - `system:kube-scheduler` for scheduler.conf
  - `system:node:cks-master` for the kubelet.conf
  - `system:kube-controller-manager` for the controller-manager.conf

In a kubeadm setup, the kubelet also runs in the master nodes and it runs other kubernetes components as static pods in the kube-system namespace.



-------------------------

### Containers in details

There are 4 namespaces :
  - PID : Pid 1 can exist multiple times inside the pid namespace, ie inside a container. The PID namespace isolates the processes running inside containers.
  - Mount : Restricts access to mounts or the root file system
  - Network : Each container can have different firewall rules, routing rules, not able to see all traffic, only have access to certain network devices etc.
  - User : Each container have an user id of 0, which is root inside that container. And this means the user id 0 inside a container is different from the user 0 on the host.

Containers use linux namespaces to restrict what the users can see and combines that with cgroups to restrict how much resources the containers can use.

Different container tools :
  - docker : container runtime + tool to manage containers and images
  - containerd : container runtime
  - crictl : cli for CRI compatible container runtimes. CRI means container runtime interface. crictl can work with docker, containerd etc.
  - podman : tool to manage container and images

podman behaves the same as docker but only for the management part of it, ie
```
podman ps
podman images
podman build -t simple .
podman run simple
```

crictl is the tool used to communicate with the container runtime, for example if you have containerd as the runtime.
You can check crictl config at `cat /etc/crictl.yaml`

Also, your runtime can be different from docker. It can be containerd.
containerd config is at `cat /etc/containerd/config.toml`

An example of running a container while attaching to the PID namespace of another container :
```
docker run --name c2 -d ubuntu sh -c 'sleep 1d'
docker run --name c2 --pid=container:c1 -d ubuntu sh -c 'sleep 9d'
docker exec c2 ps aux
docker exec c1 ps aux
```
If you list the processes inside the c1 or c2 container now, you will see both the processes for `sleep 1d` and `sleep 9d`.


--------------------------

### Network policies

By default every pod can access every other pod with out doing any network address translation.
This is a feature in the spec for CNI. So, CNI plugins like calico or weave have to implement it.

You can define ingress and egress network policies with `podSelector` or `namespaceSelector`
and these selectors are based on labels.
Besides the ones above you can also have an `ipBlock` definition for to/from traffic.
