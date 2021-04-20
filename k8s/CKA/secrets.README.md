## secrets

Secrets are like config maps except that they are supposed to store secret values in an encoded format.
Although secrets themselves are not encrypted and cant be considered safe, there are practices around handling secrets that can make it safe.
kubelet does not write secrets to disk but rather into tmpfs. If a pod that uses the secret is deleted,
kubelet will delete it's local copy of the secret from the node.
You can also enable encryuption at rest for etcd to make secrets safer to store - https://kubernetes.io/docs/tasks/administer-cluster/encrypt-data/ .

There's 2 ways to create a secret, the imperative and declarative way.

For the imperative approach there are 2 ways :

`kubectl create secret generic app-secret --from-literal=DB_USER=admin --from-literal=DB_PASSWORD=foobarbaz`
OR
`kubectl create secret generic app-secret --from-file=app-secret.properties`
where `app-secret.properties` is a simple file like so
```
DB_USER=admin
DB_PASSWORD=foobarbaz
```

For the delcarative approach :

While creating secrets in a declarative approach you must specify the secret values in a base64 encoded format

`echo -n admin | base64`
`echo -n foobarbaz | base64`

`cat app-secret.yaml`
```
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
data:
  DB_USER: base64_encoded_value_of_admin
  DB_PASSWORD: base64_encoded_value_of_foobarbaz
```

`kubectl create -f app-secret.yaml`

Some useful commands for secrets
```
kubectl get secrets
kubectl describe secret app-secret
```

To get the decoded values of the secret
`echo -n base64_encoded_value_of_foobarbaz | base64 -d`




You can load the entire secret into a pod definition yaml file or load a single
property as an env var from the secret. For example :

cat `pod-definition.yaml`
```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
    - name: nginx-container
      image: nginx
      envFrom:
        - secretRef:
            name: app-secret
```

OR load a single keys' value from the secret as an env var

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
    - name: nginx-container
      image: nginx
      env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: app-secret
              value: DB_PASSWORD
```

OR load it as a volume in the pod spec and read the secret in the application

```
apiVersion: v1
kind: Pod
metadata:
  name: nginx-pod
spec:
  containers:
    - name: nginx-container
      image: nginx
      volumeMounts:
        - name: app-secrets-volume
          mountPath: /opt/secrets
          readOnly: true
  volumes:
    - name: app-secrets-volume
      configMap:
        name: app-secrets
```

When mounting a secret as a volume, each secret will be available as a separate file in a decoded format in the mounted volume dir,
ie, /opt/secrets/DB_USER and /opt/secrets/DB_PASSWORD
