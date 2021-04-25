## authentication

There are different ways to authenticate with the kubernetes server
  1. username/password - stored in static files
  2. username/token - stored in static files
  3. certificates
  4. external auth providers like LDAP/SAML etc
  5. service accounts - for machines

While kubernetes does not natively create users, it can create and manage serviceaccounts.

```
kubectl create serviceaccount sa1
kubectl get serviceaccounts
```

All user access is managed by the kube-apiserver.

### basic auth (with username/password or bearer token based)

For basic auth, you can have a csv file with password/username/userid and use that static file as the source of authentication.
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

While authenticating pass the basic auth info to the api
`curl -v -k -u "<username>:<password>" https://master-node-ip:6443/api/v1/pods`

You can also have a group in the csv file to assign users to a group.

Similar to a static username/password file, you can also have a credentials file with username/tokens.
While authenticating pass the bearer token info to the api
`curl -v -k -H "Authorization: Bearer <token>" https://master-node-ip:6443/api/v1/pods`

In a kubeadm setup consider volume mounting this static credentials file.
Here's a more detailed setup.

`cat /tmp/credentials/user-credentials.csv`
```
welcome2kube,john_doe,uuid001
```

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
