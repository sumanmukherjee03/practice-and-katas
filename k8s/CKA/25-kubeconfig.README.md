## kubeconfig

With kubectl when the kubeconfig does not have any details or does not exist, you can still run the commands with cli flags like so
`kubectl get pods --server <kubernetes-api-server-url>:6443 --client-key kube-admin.key --client-certificate kube-admin.crt --certificate-authority ca.crt`

If you move these into a configuration file you can pass the path of that file into kubectl
`kubectl get pods --kubeconfig /path/to/kubeconfig`

If you create the kubeconfig file in `$HOME/.kube/config`, then you wont need to specify the path to the config file in your kubectl command everytime.

### Format of kubeconfig
The kubeconfig file has a specific format.
`cat $HOME/.kube/config`
```
apiVersion: v1
kind: Config
current-context: kube-admin@kubernetes-cluster
clusters:
  - name: kubernetes-cluster
    cluster:
      certificate-authority: ca.crt
      server: https://kubernetes-cluster:6443
users:
  - name: kube-admin
    user:
      client-certificate: kube-admin.crt
      client-key: kube-admin.key
      namespace: kube-system
contexts:
  - name: kube-admin@kubernetes-cluster
    context:
      user: kube-admin
      cluster: kubernetes-cluster
```

Some useful kubeconfig commands
```
kubectl config view
kubectl config use-context engineer@kubernetes-cluster
```

Regarding the certificate for the cluster, instead of using the field `certificate-authority` and providing
the full path to the server cert file, you can use the field `certificate-authority-data` and provide the actual certificate in a base64 encoded format of the original cert.
