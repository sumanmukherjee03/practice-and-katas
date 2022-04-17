## authentication

There are different ways to authenticate with the kubernetes server
  1. username/password - stored in static files
  2. username/token - stored in static files
  3. certificates
  4. external auth providers like LDAP/SAML etc
  5. service accounts - for machines

While Kubernetes does not natively create users, it can create and manage serviceaccounts.

```
kubectl create serviceaccount sa1
kubectl get serviceaccounts
```

All user access is managed by the kube-apiserver.

### basic auth (with username/password or bearer token based)

For basic auth, you can have a csv file with `password,username,userid` and use that static file as the source of authentication.
The file has 3 columns `password,username,userid`. This file is passed as an option to the kube-apiserver.

For example `cat /etc/systemd/system/kube-apiserver.service`
```
ExecStart=/usr/local/bin/kube-apiserver \\
  --advertise-address=${INTERNAL_IP} \\
  ....
  --basic-auth-file=user-credentials.csv
```

OR

If you setup the kube-apiserver using the kubeadm tool, you need to modify the manifest file for the static pod.

For example `cat /etc/kubernetes/manifests/kube-apiserver.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:
    - command:
        - kube-apiserver
        - --advertise-address=172.10.10.101
        ....
        - --basic-auth-file=user-credentials.csv
```
The change to this file will automatically restart the kube-apiserver

Essentially look for the option `--basic-auth-file` passed to the `kube-apiserver`.

While authenticating pass the basic auth info to the api
`curl -v -k -u "<username>:<password>" https://master-node-ip:6443/api/v1/pods`

You can also have a group in the csv file to assign users to a group.

Similar to a static username/password/userid file, you can also have a credentials file with username/tokens/userid instead.
While authenticating pass the bearer token info to the api
`curl -v -k -H "Authorization: Bearer <token>" https://master-node-ip:6443/api/v1/pods`

In a kubeadm setup consider volume mounting this static credentials file.
Here's a more detailed setup.

For example - `cat /tmp/credentials/user-credentials.csv`
```
welcome2kube,john_doe,uuid001
```
The format of the file above is <password>,<username>,<userid>

For example `cat /etc/kubernetes/manifests/kube-apiserver.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:
  - command:
      - kube-apiserver
      - --advertise-address=172.10.10.101
      ....
      - --basic-auth-file=/tmp/credentials/user-credentials.csv
    image: k8s.gcr.io/kube-apiserver-amd64:v1.20.0
    name: kube-apiserver
    volumeMounts:
    - mountPath: /tmp/credentials
      name: user-details
      readOnly: true
  volumes:
  - hostPath:
      path: /tmp/credentials
      type: DirectoryOrCreate
    name: user-details
```

To actually tie some permissions to the users whose creds have been passed into the kube-apiserver,
you must create a Role for RBAC and consecutively create a RoleBinding to tie the user to that role
```
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: pod-reader
rules:
- apiGroups: [""] # "" indicates the core API group
  resources: ["pods"]
  verbs: ["get", "watch", "list"]

---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: read-pods
  namespace: default
subjects:
- kind: User
  name: john_doe
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: pod-reader
  apiGroup: rbac.authorization.k8s.io
```

This should allow us to use the user defined in the static file to query the pods in the cluster.
`curl -v -k -u "john_doe:welcome2kube" https://localhost:6443/api/v1/pods`

### TLS certificates

In the context of kubernetes, the servers private key and public key signed by the CA are referred to as server certs.
The public and private key of the CA that signs the server certs are referred to as the root certs.
And a server can request the client to verify themselves using certificates and these are referred to as client certs.

Certificates or public keys would be with `*.crt` or `*.pem` extensions.
Private keys would be with `*.key` or `*-key.pem` extensions.

COMPONENTS THAT HAVE CERTS :
________________________________
  Server certs and keys :
      - kube-apiserver
            The kube-apiserver has a apiserver.crt and a apiserver.key.
            It exposes a https service that users and/or other services use to communicate with the kube-apiserver.
      - etcd
            The etcd server also has a etcdserver.crt and etcdserver.key.
            It exposes a https service that the kube-apiserver talks to for storing state of the cluster.
      - kubelet
            The kubelets on the worker nodes have a kubelet.crt and kubelet.key.
            The api server talks to the kubelets on the nodes.

  Client certs and keys :
      - kubectl
            The clients, ie admins use a admin.crt and admin.key to authenticate themselves with the kube-apiserver.
      - kube-scheduler
            The scheduler talks to the kube-apiserver and uses client certs scheduler.crt and scheduler.key
      - kube-controller-manager
            The controller manager talks to the kube-apiserver and uses client certs controller-manager.crt and controller-manager.key
      - kube-proxy
            The kube-proxy talks to the kube-apiserver and uses client certs kube-proxy.crt and kube-proxy.key
      - kube-apiserver
            The kube-apiserver is the only component that talks to etcd. As such it can use it's apiserver.crt and apiserver.key as client crt and key
            Or it can have separate client certs and keys for talking to etcd server, such as apiserver-etcd-client.crt and apiserver-etcd-client.key.
            The kube-apiserver also talks to the kubelets to monitor the status of the nodes/pods.
            As such it can use it's apiserver.crt and apiserver.key
            Or it can have separate client certs and keys to talk to kubelets, such as kubelet-client.crt and kubelet-client.key.

Kubernetes needs at least one certificate authority for the cluster to generate all these server and client certs.
You can have more than one CA - one for all the components of the cluster and one specifically for etcd.
If you have a separate CA for etcd then use that to sign the etcd server certs and also the kube-apiserver client certs used for talking to etcd.

But for convenience lets assume we have just one CA. The CA also has it's own cert and key, say ca.cert and ca.key. These are also referred to as the root certs.


COMMANDS TO GENERETE CERTS :
__________________________________

CA certs
    First generate the ca key and cert. Since this cert and key is for the CA itself
    we self sign the CSR generated for the cert.
    ```
    openssl genrsa -out ca.key 2048
    openssl req -new -key ca.key -subj "/CN=KUBERNETES-CA" -out ca.csr
    openssl x509 -req -in ca.csr -signkey ca.key -out ca.crt
    ```

Client certs for user or kubectl
    Next we generate the client key and certs for the admin user
    The subject name CN in the csr is the name the kubectl client authenticates with.
    So, in audit logs this is the name we will see. However, the name can be anything.
    While signing, use the ca.crt and ca.key instead of self signing the crt. That's the main difference from the commands above for generating the ca.crt.
    Also, when generating the csr, remember to mention the group that this kube-admin user belongs to, with something like `"/O=system:masters"` in the `-subj` parameter.
    This `system:masters` group should already exist in the kubernetes cluster with admin privileges.
    ```
    openssl genrsa -out admin.key 2048
    openssl req -new -key admin.key -subj "/CN=kube-admin/O=system:masters" -out admin.csr
    openssl x509 -req -in admin.csr -CA ca.crt -CAkey ca.key -out admin.crt
    ```

    This client cert and key generated for the kube-admin and the CA cert can be used when making api calls to the kube-apiserver later on like so
    `curl https://kube-apiserver:6443/api/v1/pods --key admin.key --cert admin.crt --cacert ca.crt`.

    This is one way of using the client certs. The other way is to move all this to `kube-config.yaml` and use it via kubectl
    ```
    apiVersion: v1
    kind: Config
    clusters:
      - cluster::
          certificate-authority: ca.crt
          server: https://kube-apiserver:6443
        name: kubernetes
    users:
      - name: kube-admin
        user:
          client-certificate: admin.crt
          client-key: admin.key
    ```

    Similarly generate all the other client certificats. Only difference from the above client certs and keys being
    that the names should be prefixed with `system:` in the `CN` section of the `-subj` parameter.
    So, `CN=system:kube-scheduler` , `CN=system:kube-controller-manager` , `CN=system:kube-proxy` and so on.

ETCD server certs
    Here's the commands to generate the certs for etcd-server.
    The CSR ofcourse needs to be signed by the CA cert with the CA key and not self signed.
    ```
    openssl genrsa -out etcdserver.key 2048
    openssl req -new -key etcdserver.key -subj "/CN=etcd-server" -out etcdserver.csr
    openssl x509 -req -in etcdserver.csr -CA ca.crt -CAkey ca.key -out etcdserver.crt
    ```
    If you are using a HA etcd cluster, then you will have to configure certs for the etcd peers as well.
    You can call those `etcdpeer.crt`, `etcdpeer.key`

    The manifest for etcd.yaml would contain a section like this
    ```
    - etcd
    - --advertise-client-urls=https://127.0.0.1:2379
    - --listen-client-urls=https://127.0.0.1:2379
    - --cert-file=/etc/kubernetes/pki/etcd/etcdserver.cert
    - --key-file=/etc/kubernetes/pki/etcd/etcdserver.key
    - --trusted-ca-file=/etc/kubernetes/pki/etcd/ca.cert
    - --client-cert-auth=true
    - --name=master
    - --initial-advertise-peer-urls=https://127.0.0.1:2380
    - --listen-peer-urls=https://127.0.0.1:2380
    - --data-dir=/var/lib/etcd
    ......
    - --peer-cert-file=/etc/kubernetes/pki/etcd/etcdpeer.cert
    - --peer-key-file=/etc/kubernetes/pki/etcd/etcdpeer.key
    - --peer-trusted-ca-file=/etc/kubernetes/pki/etcd/ca.cert
    - --peer-client-cert-auth=true
    ```

kube-apiserver certs
    Here's the commands to generate the certs for kube-apiserver.

    The first step is to generate the key as usual.
    ```
    openssl genrsa -out apiserver.key 2048
    ```

    A lot of people refer to the kube-apiserver in many ways because this is the main point of interaction of everything and everyone with kubernetes.
    Thus, the cert for this has to support many alternate names. These alternate names can also include the kube-apiserver node IP and pod IP.
    To specify the alternate names we need to create a `openssl.cnf` file and specify the alternate names there.
    `cat openssl.cnf`
    ```
    [req]
    req_extensions = v3_req
    [v3_req]
    basicConstraints = CA:FALSE
    keyUsage = nonRepudiation
    subjectAltName = @alt_names
    [alt_names]
    DNS.1 = kubernetes
    DNS.2 = kubernetes.default
    DNS.3 = kubernetes.default.svc
    DNS.4 = kubernetes.default.svc.cluster.local
    IP.1 = 10.96.0.1
    IP.2 = 172.17.0.87
    ```
    As you can see the alternate names need to be specified in the cnf file.
    Now use that file while generating the CSR
    ```
    openssl req -new -key apiserver.key -subj "/CN=kube-apiserver" -config openssl.cnf -out apiserver.csr
    ```

    Lastly, the CSR ofcourse needs to be signed by the CA cert with the CA key and not self signed.
    ```
    openssl x509 -req -in apiserver.csr -CA ca.crt -CAkey ca.key -out apiserver.crt
    ```

    It's important to note that the kube-apiserver acts as a client when talking to the kubelet and etcd.
    As such we need to generate certs for kube-apiserver to work as client certs for talking to etcd cluster and kubelet.
    The process would be similar to the above except for that the alternate names are not going to be needed.
    Lets say we call these files `apiserver-etcd-client.crt` and `apiserver-etcd-client.key` and `apiserver-kubelet-client.crt` and `apiserver-kubelet-client.key`.

    If you had a non kubeadm setup then `cat /etc/systemd/system/kube-apiserver.service`
    ```
    ExecStart=/usr/local/bin/kube-apiserver \
      --advertise-address=${INTERAL_IP} \
      --allow-privileged=true \
      --apiserver-count=3 \
      --authorization-mode=Node,RBAC \
      --bind-address=0.0.0.0 \
      --enable-swagger-ui=true \
      --etcd-ca-file=/var/lib/kubernetes/ca.crt \
      --etcd-certfile=/var/lib/kubernetes/apiserver-etcd-client.crt \
      --etcd-keyfile=/var/lib/kubernetes/apiserver-etcd-client.key \
      --etcd-servers=https://127.0.0.1:2379 \
      --kubelet-certificate-authority=/var/lib/kubernetes/ca.crt \
      --kubelet-client-certificate=/var/lib/kubernetes/apiserver-kubelet-client.crt \
      --kubelet-client-key=/var/lib/kubernetes/apiserver-kubelet-client.key \
      --kubelet-https=true \
      --client-ca-file=/var/lib/kubernetes/ca.crt \
      --tls-cert-file=/var/lib/kubernetes/apiserver.crt \
      --tls-private-key=/var/lib/kubernetes/apiserver.key \
      ......
      --v=2
    ```
    The api server needs the ca.crt to verify it's client. Also, note it is the same ca.crt that is passed as an option
    for the kube-apiserver to verify the server certs of the etcd cluster and the kubelets.
    Note how the client certificates are separate for the kube-apiserver here.

kubelet certs
    The certs and keys for the kubelets will be different on each node. They need the node name as part of the certificate name.
    This is because when the kubelet talks to the kube-apiserver, the apiserver needs to know which node is interacting.
    It uses this information for authorization purposes. Also, these certs needs to have the proper group information for the nodes.
    This group is again system level group. So the group name would be like `system:nodes` and cert name would be like `system:node:node01`.
    ```
    openssl genrsa -out node01.key 2048
    openssl req -new -key node01.key -subj "/CN=system:node:node01/O=system:nodes" -out node01.csr
    openssl x509 -req -in node01.csr -CA ca.crt -CAkey ca.key -out node01.crt
    ```
    You can use this information in the kubelet-config.yaml
    Again, in this config you can see we use the same ca.crt file for validating client cert that kube-apiserver would be using to talk to the kubelet.
    ```
    apiVersion: kubelet.config.k8s.io/v1beta1
    kind: KubeletConfiguration
    authentication:
      x509:
        clientCAFile: /var/lib/kubernetes/ca.crt
    authorization:
      mode: Webhook
    clusterDomain: cluster.local
    clusterDNS:
      - "10.32.0.10"
    podCIDR: "${POD_CIDR}"
    resolvConf: /run/systemd/resolve/resolv.conf
    runtimeRequestTimeout: 15m
    tlsCertFile: /var/lib/kubelet/node01.crt
    tlsPrivateKeyFile: /var/lib/kubelet/node01.key
    ```

COMMANDS TO INSPECT CERTS
_______________________________

To decode certs first obtain the path of the certs by looking at the command to start the kube-apiserver.
Then use the openssl utility to decode certs. For example :
```
openssl x509 -in /etc/kubernetes/pki/apiserver.crt -text -noout
```
Look for the Subject > CN to get the name of the cert, Validity > Not After to get the expiry date, Issuer > CN to get the CA who issued the cert etc.

Inspect logs to debug cert related issues.
If you manually setup the cluster, search for service logs
`journalctl -u etcd.service -l`
OR ELSE search for logs in the pods if you set up the cluster with kubeadm
`kubectl logs etcd-master`
If kube-apiserver is down then you might have to view the logs from the docker container in the master nodes.


KUBECTL COMMANDS TO SIGN CERTIFICATES
___________________________________________

Assume that a new user has generated a certificate `openssl genrsa -out john.key 2048`.
Then he creates a CSR `openssl req -new -key john.key -subj "/CN=john" -out john.csr`
He then sends that CSR to an admin who has to get it signed by the kubernetes CA.

The admin creates `john-csr.yaml` with the base64 encoded value of the CSR `cat john.csr | base64`
```
apiVersion: certificates.k8s.io/v1beta1
kind: CertificateSigningRequest
metadata:
  name: john
spec:
  groups:
    - system:authenticated
  usages:
    - digital signature
    - key encipherment
    - server auth
  request:
    .....encoded value of csr......
```
The admin submits this CertificateSigningRequest which can be signed by the CA using these commands
```
kubectl get csr
kubectl certificate approve john
kubectl get csr john -o yaml
```
Get the signed certificate from the yaml output and base64 decode it.
This can then be provided to john so that john can access the kube-apiserver using the new cert through kubectl.

The kubernetes component that signs certificates is the controller-manager.
The controller manager has controllers called `csr-approving`, `csr-signing` for dealing with certs.
The controller manager when starting up has these 2 options to get the paths of the root CA cert and key
  - `--cluster-signing-cert-file`
  - `--cluster-signing-key-file`

You can deny a CSR as an administrator through `kubectl`.
`kubectl certificate deny rogue-req`
