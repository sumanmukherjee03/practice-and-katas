## ingress

A service can expose an application externally via NodePort. You can have a DNS pointing to the IP of the service.
However clients will still have to keep typing the random port allocated during the NodePort. Lets say a port like 35000
which forwards traffic via the service to a container port of 5000 say.
For clients to not have to remember the port 35000, we can add a proxy server like nginx in front of the service.
This is where an ingress comes into play. Ingress does what the proxy server does. Point your DNS to the ingress proxy server.

In the cloud a service type load balancer will create a loadbalancer per service.
So, external DNS -> loadbalancer:35000 -> pod:5000
If you add another service that will have a loadbalancer as well.
But you would want to forward traffic to the appropriate loadbalancer based on path.
You also want to handle SSL at one place. This again is a good candidate for a reverse proxy. And hence an Ingress.
Also, a loadbalancer per service will rack up cloud cost pretty rapidly.

Essentially Ingress is a layer 7 loadbalancer. Of course ingress still has to be exposed via a NodePort or a cloud loadbalancer.

Ingress controller and ingress resources combined together does what a reverse proxy would do.
Ingress controllers are not enabled by default when you create a kubernetes cluster. So, it has to be created via some kind of deployment.

There are a few different kinds of ingress controllers
  - GCP loadbalancer
  - Nginx
  - Contour
  - HA Proxy
  - Traefik
  - Istio

Nginx and GCP loadbalancer is maintained by the kubernetes team.




--------------------------------------------------------------------------------

### ingress controller

This is an example of an ingress controller deployment for nginx.
You need a configmap first to manage the configuration for the nginx server and that is passed along to the binary when starting the nginx server.
This configmap contains nginx specific configuration like timeouts.
For the nginx deployment we use a special image from kubernetes.
The container the POD_NAME and POD_NAMESPACE env vars to operate.
Expose the ports 80 and 443 via a service.
The ingress controller monitors the ingress resources and configures the nginx server as things change in our cluster.
However, to be able to do this the ingress controller needs additional permissions.
That's why we also need a ServiceAccount object with the proper Roles and RoleBindings.

Example nginx ingress controller deployment - `cat ingress-controller.yaml`

```
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configuration
data:
  proxy-connect-timeout: "10"
  proxy-read-timeout: "120"
  proxy-send-timeout: "120"

---

apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: nginx-ingress-controller
spec:
  replicas: 1
  selector:
    matchLabels:
      name: nginx-ingress
  template:
    metadata:
      labels:
        name: nginx-ingress
    spec:
      containers:
        - name: nginx-ingress-controller
          image: quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.21.0
      args:
        - /nginx-ingress-controller
        - --configmap=$(POD_NAMESPACE)/nginx-configuration
      env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
      ports:
        - name: http
          containerPort: 80
        - name: https
          containerPort: 443

---

apiVersion: v1
kind: Service
metadata:
  name: nginx-ingress
spec:
  type: NodePort
  ports:
    - name: http
      port:80
      targetPort: 80
      protocol: TCP
    - name: https
      port:443
      targetPort: 443
      protocol: TCP
  selector:
    name: nginx-ingress

---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: nginx-ingress-serviceaccount
```




--------------------------------------------------------------------------------

### ingress resources

Sample routing `cat ingress-warehouse.yaml`
This one is extremely simple as it forwards all traffic to a backend service.
```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress-warehouse
spec:
  backend:
    serviceName: warehouse-service
    servicePort: 8080
```


This one is a more complicated example for an ingress resource with rules and paths and domains.
`cat ingress-online-store.yaml`

```
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: ingress-online-store
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
    - host: accessories.bestdealscom
      http:
        paths:
          - path: /watches
            backend:
              serviceName: watch-store-service
              servicePort: 8080
          - path: /hats
            backend:
              serviceName: hat-store-service
              servicePort: 8081
          - path: /belts
            backend:
              serviceName: belt-store-service
              servicePort: 8082
    - host: clothes.bestdealscom
      http:
        paths:
          - path: /shirts
            backend:
              serviceName: shirt-store-service
              servicePort: 8083
          - path: /pants
            backend:
              serviceName: pant-store-service
              servicePort: 8084
          - path: /blouses
            backend:
              serviceName: blouse-store-service
              servicePort: 8085
          - path: /skirts
            backend:
              serviceName: skirt-store-service
              servicePort: 8086
```
REMEMBER : If none of the paths above match, the ingress resource is gonna forward the traffic to `default-http-backend:80`.
So, we should not forget to deploy such a service.
If you dont have a hostname it will match anything (or wildcard) for the hostname.
Take note of the annotation `nginx.ingress.kubernetes.io/reqrite-target` which can be used to rewrite urls.

Here's another example with a rewrite
```
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /$2
  name: rewrite
  namespace: default
spec:
  rules:
  - host: rewrite.bar.com
    http:
      paths:
      - backend:
          serviceName: http-svc
          servicePort: 80
        path: /something(/|$)(.*)
```
Here rewrite.bar.com/something rewrites to rewrite.bar.com/ and rewrite.bar.com/something/new rewrites to rewrite.bar.com/new

More on nginx ingress controller here : https://kubernetes.github.io/ingress-nginx/examples/


Some helpful kubectl commands for ingress resources
```
kubectl create -f ingress-online-store.yaml
kubectl get ingress
kubectl describe ingress ingress-online-store
```
