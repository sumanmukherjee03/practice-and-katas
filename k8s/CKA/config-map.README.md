## config-map

Config maps are used to pass key/value pairs of data into kubernetes pod specs.

You can load the entire config map into a pod definition yaml file or load a single
property as an env var from the config map. For example :

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
        - configMapRef:
            name: app-config
```

OR load a single keys' value from the config map

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
        - name: COLOR
          valueFrom:
            configMapKeyRef:
              name: app-config
              value: COLOR
```

OR load it as a volume in the pod spec and read the config map in the application

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
        - name: app-config-volume
          mountPath: /data/app-config.properties
          readOnly: true
  volumes:
    - name: app-config-volume
      configMap:
        name: app-config
```


You can create the config map in an imperative or declarative way using a manifest.

`kubectl create configmap app-config --from-literal=COLOR=blue --from-literal=STAGE=dev`

OR

`kubectl create configmap app-config --from-file=app-config.properties`
where `app-config.properties` can be as simple as
```
COLOR=blue
STAGE=dev
```

The structure for a declarative approach is provided below :
cat config-map.yaml
```
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  COLOR: blue
  STAGE: dev
```

`kubectl create -f config-map.yaml`

Some useful config map commands
```
kubectl get configmaps
kubectl describe configmap app-config
```
