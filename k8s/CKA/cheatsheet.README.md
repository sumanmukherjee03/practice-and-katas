## cheatsheet

This is a file with not much explanation but only commands.

```
kubectl cluster-info
kubectl cluster-info dump
kubectl get all --all-namespaces
kubectl get pods --all-namespaces
kubectl -n kube-system get pods -o wide
kubectl get pods --no-headers --selector tier=frontend,env=prod
kubectl exec -it webapp -c app -- tail -f /var/log/app.log
kubectl run --rm --restart=Never redis --image redis --port 6379 --expose
kubectl run --rm --restart=Never ubuntu --image ubuntu --command sleep 900 --dry-run=client -o yaml
kubectl run --rm --restart=Never busybox --image=gcr.io/google-containers/busybox --command sleep 900 --overrides='{"apiVersion": "v1", "spec": {"template": {"spec": {"nodeSelector": {"kubernetes.io/hostname": "node01"}}}}}'
kubectl exec -it webapp -c app -- nslookup db-service
kubectl describe pod webapp-pod | grep -i image
kubectl get pod nginx --watch
kubectl run --rm --restart=Never --image=nikola/netshoot --command /bin/sh -c 'while true; do echo "HTTP/1.1 200 OK\n SUCCESS" | nc -l -p 80 -q 1; done' --port 80 --expose
kubectl get pods -l tier=webapp -o jsonpath='{{range .items}}{{.status.podIP}}{{"\n"}}{{end}}'

kubectl get deployments
kubectl create deployment frontend --image nginx --replicas=2
kubectl scale deployment frontend --replicas=3
kubectl set image deployment/frontend nginx=nginx:1.18 --record
kubectl expose deployment nginx --port 80
kubectl expose deployment webapp --name=webapp-service --port=8080 --target-port=8080 --type=NodePort --dry-run=client -o yaml
kubectl get svc -o wide
kubectl get endpoints -o wide
```
