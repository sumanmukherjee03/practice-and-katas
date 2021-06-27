## cheatsheet
```
kubectl cluster-info
kubectl cluster-info dump
kubectl get pods --server https://master-loadbalancer:6443 --client-key kube-admin.key --client-certificate kube-admin.crt --certificate-authority ca.crt
kubectl get pods --kubeconfig /path/to/kubeconfig
kubectl config view
kubectl config current-context
kubectl config use-context engineer@kubernetes-cluster
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

kubectl get nodes -o wide
kubectl label nodes node01 size=large
kubectl create serviceaccount sa1
kubectl get serviceaccounts
curl -v -k -u '<username>:<password>' https://master-loadbalancer:6443/api/v1/pods
curl -v -k -H 'Authorization: Bearer <token>' https://master-loadbalancer:6443/api/v1/pods
curl -v -k https://master-loadbalancer:6443/api/v1/pods --key admin.key --cert admin.crt --cacert ca.crt

kubectl get csr
kubectl csr approve john
kubectl get csr john -o yaml
kubectl certificate deny rogue-req

kubectl proxy; curl -k http://localhost:6443/apis;

kubectl api-resources --namespaced=true
kubectl api-resources --namespaced=false

kubectl create role engineer --verb=create --verb=get --verb=list --verb=delete --resource=pods --dry-run=client -o yaml
kubectl get roles
kubectl get rolebindings
kubectl describe role engineer
kubectl create rolebinding engineering-user-role-binding --role=engineer --user=engineering-user --dry-run=client -o yaml
kubectl describe rolebinding engineering-user-role-binding

kubectl auth can-i create deployments
kubectl auth can-i delete pods
kubectl auth can-i delete pods --as john
kubectl auth can-i delete pods --as john --namespace prod

kubectl version --short
kubeadm token list
apt-get -y -qq update; apt install -y kubeadm=1.20.0-00; kubeadm upgrade plan; kubeadm upgrade apply v1.20.0; kubectl get nodes; kubectl drain controlplane --ignore-daemonsets; apt-get install -y kubelet=1.20.0-00; apt-get install -y kubectl=1.20.0-00; systemctl daemon-reload; systemctl restart kubelet; kubectl uncordon controlplane;
kubectl drain node01 --ignore-daemonsets
ssh node01; apt-get install -y kubeadm=1.20.0-00; apt-get install -y kubelet=1.20.0-00; kubeadm upgrade node; systemctl daemon-reload; systemctl restart kubelet
kubectl uncordon node01

ETCDCTL_API=3 etcdctl snapshot save snapshot.db --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=https://127.0.0.1:2379
ETCDCTL_API=3 etcdctl snapshot status snapshot.db --cert=/etcd/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=https://127.0.0.1:2379
systemctl kube-apiserver stop
ETCDCTL_API=3 etcdctl snapshot restore snapshot.db --data-dir=/var/lib/etcd-from-backup --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt
sed -i -e 's#--data-dir=/var/lib/etcd#--data-dir=/var/lib/etcd-from-backup#g' /etc/systemd/system/etcd.service
systemctl daemon-reload; systemctl etcd restart; systemctl kube-apiserver restart

kubectl create configmap app-config --from-literal=COLOR=blue --from-literal=ENVIRONMENT=staging
kubectl create configmap app-config --from-file=app-config.properties
kubectl get configmaps
kubectl describe configmap app-config

kubectl get daemonsets
kubectl describe daemonsets monitoring-agent-daemon

cat /etc/systemd/system/kube-apiserver.service
ps -aux | grep kube-apiserver

cat /etc/systemd/system/kube-controller-manager.service
ps -aux | grep kube-controller-manager

cat /etc/systemd/system/kubelet.service | grep '--config'
cat /var/lib/kubelet/kubelet-config.yaml | grep '--pod-manifest-path'

kubectl -n kube-system get daemonset kube-proxy
kubectl get events | grep -i <custom-scheduler-name>

kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.cpu}{"\n"}{end}'
kubectl get nodes -o custom-columns=NODE:.metadata.name,CPU:.status.capacity.cpu
kubectl get pv -o custom-columns=NAME:.metadata.name,CAPACITY:.spec.capacity.storage --sort-by=.spec.capacity.storage
kubectl config view --kubeconfig=/path/to/kubeconfig -o jsonpath='{$.contexts[?(@.context.user == "aws-user")].name}{"\n"}'

kubectl create ns dev
kubectl get ns --no-headers
kubectl create quota dev-ns-count --hard=count/deployments.apps=2,count/replicasets.apps=4,count/pods=10,count/secrets=4 --namespace=dev
kubectl describe quota compute-quota
kubectl get quota --namespace=dev

kubectl get pv
kubectl describe pv pv-vol1
kubectl get pvc
kubectl describe pvc pv-vol1-claim
kubectl delete pvc pv-vol1-claim
kubectl get storageclass portworx-vol

kubectl taint node node01 color=green:NoExecute
kubectl taint node node01 color=green:NoExecute-
kubectl explain pod --recursive | grep -A5 tolerations
kubectl taint node controlplane node-role.kubernetes.io/master:NoSchedule-
```
