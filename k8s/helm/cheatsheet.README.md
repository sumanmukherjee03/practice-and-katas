## helm cheatsheet

Cheatsheet commands to help with helm

```
helm env
helm repo add jetstack https://charts.jetstack.io
helm search repo cert-manager
helm template cert-01 jetstack/cert-manager --namespace cert-manager --version 1.5.3 --values cert-manager-values.yaml --output-dir helm-generated-yaml

helm search hub wordpress

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
helm search repo joomla
helm search repo joomla --debug --version 12

helm fetch bitnami/postgresql --version 11.6.2

helm repo list
helm search repo drupal --versions
helm install --set name=prod myredis ./redis
helm install bravo bitnami/drupal
helm uninstall bravo
helm show all bitnami/drupal

helm list
helm pull --untar bitname/apache
helm install mywebapp ./apache

helm repo add kong https://charts.konghq.com
helm repo update
helm install kong/kong --generate-name --set ingressController.installCRDs=false --set admin.enabled=true
helm install sample-api-server sample-api-server/
helm upgrade sample-api-server sample-api-server/
```
