## security context

Kubernetes allows you to set security context on pods or individual containers.
Setting on a pod is applied to all containers and setting on a container overrides the setting on a pod.

Security context can control different ways you can handle ACL inside linux containers for access to files, directories, sockets etc.
  - uid/gid based ACL
  - SELinux based controls
  - give a container some linux capabilities but not all the privileges of a root user
  - AllowPrivilegeEscalation : This controls whether a container process can gain more access than it's parent process.
                               It is essentially controlling whether `no_new_privs` flag gets set on the PID 1 inside the container
  - readOnlyRootFilesystem : Mounts the containers root file system as read only
  - seccompProfile : If the nodes have been configured with the seccomp profile json files,
                     ie the kubelets have the seccomp profile files in their seccomp path `/var/lib/kubelet/seccomp/profiles`
                     then you can run the pods with the appropriate seccomp profile to restrict your
                     containers to only make certain syscalls.

This is an example of a security context to run a container as a non-root user
and non-root group and also to not allow privilege escalation for any proc within the container.
`cat pod-definition.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: nginx-container
      image: nginx
  securityContext:
    runAsUser: 1000
    runAsGroup: 2000
    allowPrivilegeEscalation: false
```

To set the security context at the container level as a non-root user
and running the container with some added capabilities.

`cat pod-definition.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
  labels:
    app: nginx
spec:
  containers:
    - name: nginx-container
      image: nginx
      securityContext:
        runAsUser: 1000
        capabilities:
          add: ["MAC_ADMIN", "SYS_TIME"]
```

To load a container with a seccomp profile
```
apiVersion: v1
kind: Pod
metadata:
  name: fine-pod
  labels:
    app: fine-pod
spec:
  securityContext:
    seccompProfile:
      type: Localhost
      localhostProfile: profiles/fine-grained.json
  containers:
  - name: test-container
    image: hashicorp/http-echo:0.2.3
    args:
    - "-text=just made some syscalls!"
    securityContext:
      allowPrivilegeEscalation: false
```
You could also load the default seccomp profile from kubernetes if you dont want to use a custom profile
Relevant section to reconfigure is
```
....
  securityContext:
    seccompProfile:
      type: RuntimeDefault
....
```

Here's an example of a seccomp profile definition json file.
The kubelets load them into a pod that the pod can refer to for use by the containers.
```
{
  "defaultAction": "SCMP_ACT_ERRNO",
  "architectures": [
    "SCMP_ARCH_X86_64",
    "SCMP_ARCH_X86",
    "SCMP_ARCH_X32"
  ],
  "syscalls": [
    {
      "names": [
        "read",
        "write",
        "close",
        "getpid",
        "getuid",
        "open",
        "poll",
        ......
          "ioctl"
      ],
      "action": "SCMP_ACT_ALLOW"
    }
  ]
}
```
