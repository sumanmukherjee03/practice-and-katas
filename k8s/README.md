Some commands to help with the project

```sh
helm repo add kong https://charts.konghq.com
helm repo update
helm install kong/kong --generate-name --set ingressController.installCRDs=false --set admin.enabled=true
kubectl get svc
kubectl get pods
curl -s localhost:32547 | jq -r '.plugins.available_on_server | ."rate-limiting"'
curl localhost:32547/plugins
curl localhost:32547/routes
curl localhost:32547/upstreams
curl localhost:32547/targets
helm install sample-api-server sample-api-server/
helm upgrade sample-api-server sample-api-server/
kubectl get kongplugins.configuration.konghq.com
kubectl describe Ingress sample-api-server
kubectl logs kong-1615734873-kong-7b996d657f-x4jk4 -f ingress-controller
curl -vvv http://localhost:31615/foo
curl -vvv -H "HOST: sample-api-server.local.dev" -H "PLAN-NAME: free" http://localhost/test
```
