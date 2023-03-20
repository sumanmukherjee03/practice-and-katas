## cheatsheet
```
kubectl api-resources --namespaced=true
kubectl api-resources --namespaced=false
kubectl api-resources --api-group=extensions

kubectl cluster-info
kubectl cluster-info dump
kubectl config view
kubectl config current-context
kubectl config use-context engineer@kubernetes-cluster
kubectl get all --all-namespaces
kubectl get pods --all-namespaces
kubectl -n kube-system get pods -o wide
kubectl describe pod webapp-pod | grep -i image
kubectl get pod nginx --watch
kubectl get deployments
kubectl get svc -o wide
kubectl get endpoints -o wide
kubectl get nodes -o wide
kubectl create serviceaccount sa1
kubectl get serviceaccounts
kubectl api-resources --namespaced=true
kubectl api-resources --namespaced=false
kubectl get roles
kubectl get rolebindings
kubectl describe role engineer
kubectl describe rolebinding engineering-user-role-binding
kubectl version --short
kubectl get configmaps
kubectl describe configmap app-config
kubectl get daemonsets
kubectl describe daemonsets monitoring-agent-daemon
kubectl -n kube-system get daemonset kube-proxy
kubectl create ns dev
kubectl get ns --no-headers
kubectl describe quota compute-quota
kubectl get quota --namespace=dev
kubectl get pv
kubectl describe pv pv-vol1
kubectl get pvc
kubectl describe pvc pv-vol1-claim
kubectl delete pvc pv-vol1-claim
kubectl get storageclass portworx-vol
kubectl logs web -f


kubectl describe pod <pod_name> | grep -i node
kubectl describe pod <pod_name> | grep -i image
kubectl run redis --image=redis123 --dry-run=client -o yaml
kubectl get replicasets
kubectl describe replicaset new-replica-set | grep -i -E '(image|container)'
kubectl api-resources --namespaced=true | grep -i replicaset
kubectl delete -f replicaset-definitions.yaml
kubectl get replicaset new-replica-set -o yaml > replicaset-def.yaml
kubectl scale replicaset new-replica-set --replicas=0
kubectl scale replicaset new-replica-set --replicas=4
kubectl edit replicaset new-replica-set
kubectl get deployments
kubectl describe deployment frontend-deployment | grep -i image
kubectl create deployment http-frontend --image=httpd:2.4-alpine --dry-run=client -o yaml > deployment-httpd.yaml
kubectl get ns --no-headers | wc -l
kubectl run redis --image=redis -n finance
kubectl get pods -A -o wide | grep -i blue
kubectl describe svc kubernetes | grep -i type
kubectl describe svc kubernetes | grep -i targetport
kubectl get endpoints kubernetes
kubectl describe deployment simple-webapp-deployment | grep -i image
kubectl run nginx-pod --image=nginx:alpine
kubectl run redis --image=redis:alpine --labels=tier=db
kubectl run redis --image=redis:alpine -l tier=db
kubectl expose pod redis --name redis-service --port 6379 --target-port 6379 --type ClusterIP
kubectl expose pod redis --name redis-service --port 6379 --target-port 6379 --type ClusterIP --dry-run=client -o yaml
kubectl run --restart=Never --image=busybox static-busybox --dry-run=client -o yaml --command -- sleep 1000 > /etc/kubernetes/manifests/static-busybox.yaml
kubectl create deployment webapp --image=kodekloud/webapp-color --replicas=3
kubectl run custom-nginx --image=nginx --port=8080 --expose
kubectl create deployment redis-deploy -n dev-ns --image=redis --replicas=2
kubectl run httpd --image=httpd:alpine --port=80 --expose
kubectl describe pod ubuntu-sleeper | grep -C 5 -i command
kubectl describe pod webapp-color | grep -C 4 -i environment
kubectl get pods --field-selector=status.phase=Running
kubectl explain pod.spec.containers.envFrom.configMapRef
kubectl get secret
kubectl get secret <secret> -o yaml
kubectl describe secret <secret> | grep -i type
kubectl create secret generic db-secret --from-literal=DB_Host=sql01 --from-literal=DB_User=root --from-literal=DB_Password=<whatever>
kubectl explain pod.spec.containers.securityContext.capabilities
kubectl explain pod.spec.containers.securityContext.runAsUser
kubectl describe pod elephant | grep -C 5 -i terminated
kubectl describe sa default
kubectl get pod <pod_name> -o yaml | grep -C 3 -i serviceaccount
kubectl exec <pod_name> -c <container_name> -it -- cat /var/run/secrets/kubernetes.io/serviceaccount/token
kubectl get secret dashboard-sa-token -o jsonpath='{.data.token}'
kubectl explain deployment.spec.template.spec | grep -i serviceaccount
kubectl explain pod --recursive | grep -C5 tolerations
kubectl explain deployment --recursive | less
kubectl taint node node01 color=blue:NoExecute
kubectl taint node node01 color=blue:NoExecute-
kubectl describe node kubemaster | grep -i taint
kubectl label nodes node01 size=large
kubectl get nodes
kubectl explain deployment --recursive | grep -C 5 -i nodeselector
kubectl taint node node01 spray=mortein:NoSchedule
kubectl describe node node01 | grep -i taints
kubectl taint node controlplane node-role.kubernetes.io/master=:NoSchedule-
kubectl -n elastic-stack logs kibana
kubectl -n elastic-stack exec app -c app -it -- tail -f /log/app.log
kubectl explain pod.spec.containers.readinessProbe --recursive
kubectl explain pod.spec.containers.livenessProbe --recursive
kubectl logs webapp-2 -c simple-webapp -f | grep -i -E '(fail|error)'
kubectl top node
kubectl top pod
kubectl get pods --no-headers --selector env=dev | wc -l
kubectl get all --no-headers --selector env=prod
kubectl get pods --no-headers --selector env=prod,bu=finance,tier=frontend
kubectl set image deployments/frontend simple-webapp=kodekloud/webapp-color:v2
kubectl create job throw-dice-job --image=kodekloud/throw-dice
<!-- kubectl create cronjob throw-dice-cron-job --image=kodekloud/throw-dice --schedule="30 21 * * *" -->
kubectl get networkpolicies
kubectl describe networkpolicy payroll-policy
kubectl exec internal -c internal -it -- nc -v -z 10-50-0-6.default.pod.cluster.local 8080
kubectl explain networkpolicy.spec --recursive
kubectl get all -A | grep -v -i kube-system
kubectl get deployments -n app-space --no-headers
kubectl describe ingress -n app-space ingress-wear-watch
kubectl describe svc default-http-backend -n app-space
kubectl edit ingress -n app-space ingress-wear-watch
kubectl explain pod.spec.containers.volumeMount --recursive
kubectl explain pod.spec.volumes --recursive
kubectl explain persistentvolume.spec --recursive
kubectl explain persistentvolumeclaim.spec --recursive | grep -i 'persistentvolumereclaimpolicy'
kubectl get persistentvolumeclaim
kubectl delete persistentvolumeclaim myclaim
kubectl explain storageclass --recursive
kubectl get persistentvolume
kubectl describe persistentvolume pv-vol1
kubectl describe storageclass local-storage
kubectl get nodes -o wide
kubectl explain statefulset.spec --recursive
kubectl explain statefulset.spec.volumeclaimtemplates --recursive
kubectl explain networkpolicy.spec --recursive | grep -C 10 -i -E '(ingress|egress)'
kubectl logs web --previous
kubectl run -it --rm --restart=Never busybox --image=busybox sh
kubectl get pods -l app=hostnames -o jsonpath='{{range .items}}{{.status.podIP}}{{"\n"}}{{end}}'



kubectl create configmap app-config --from-literal=APP_COLOR=blue --from-literal=APP_MODE=prod
kubectl create configmap app-config --from-file=app_config.properties
kubectl describe configmap app-config
kubectl create secret generic app-secret --from-literal=USER=foo --from-literal=ENGINE=mysql
kubectl create secret generic app-secret --from-file=app_secret.properties


kubectl create serviceaccount dashboard-sa
kubectl get serviceaccount
kubectl describe serviceaccount dashboard-sa | grep -i token
kubectl describe secret dashboard-sa-token-kubeddm
curl https://apiserverurl/api -insecure --header "Authorization: Bearer <token_value_from_above>"
kubectl describe pod dashboard-app | grep -i secrets | grep -i serviceaccount
kubectl exec -it dashboard-app cat /var/run/secrets/kubernetes.io/serviceaccount/token


kubectl get ingress
kubectl describe ingress test-ingress
kubectl create ingress test-ingress --rule="wear.onlinestore.com/wear=wear-service:80"


kubectl create ns ingress-space
kubectl config set-context --current --namespace=ingress-space
kubectl create configmap nginx-configuration
kubectl create serviceaccount ingress-serviceaccount
kubectl describe rolebinding ingress-role-binding
kubectl describe role ingress-role
kubectl describe deployment ingress-controller
kubectl expose deployment ingress-controller --name=ingress --port=80 --target-port=80 --type=NodePort --dry-run=client -o yaml > ingress-svc.yaml
kubectl create ingress app-ingress -n app-space --rule="/wear*=wear-service:8080" --rule="/watch*=video-service:8080" --annotations 'nginx.ingress.kubernetes.io/rewrite-target=/'


echo -n 'mysecretvalue' | base64
echo -n 'mysecretfromk8s' | base64 -d


ls /usr/include/linux.capability.h
docker run --user 1000 ubuntu sleep 3600
docker run --cap-add MAC_ADMIN ubuntu
docker run --cap-drop KILL ubuntu
docker run --privileged ubuntu



kubectl get pods --server https://master-loadbalancer:6443 --client-key kube-admin.key --client-certificate kube-admin.crt --certificate-authority ca.crt
kubectl get pods --kubeconfig /path/to/kubeconfig
kubectl get pods --no-headers --selector tier=frontend,env=prod
kubectl exec -it webapp -c app -- tail -f /var/log/app.log
kubectl run --rm --restart=Never redis --image redis --port 6379 --expose
kubectl run --rm --restart=Never ubuntu --image ubuntu --command sleep 900 --dry-run=client -o yaml
kubectl run --rm --restart=Never busybox --image=gcr.io/google-containers/busybox --command sleep 900 --overrides='{"apiVersion": "v1", "spec": {"template": {"spec": {"nodeSelector": {"kubernetes.io/hostname": "node01"}}}}}'
kubectl exec -it webapp -c app -- nslookup db-service
kubectl run --rm --restart=Never --image=nikola/netshoot --command /bin/sh -c 'while true; do echo "HTTP/1.1 200 OK\n SUCCESS" | nc -l -p 80 -q 1; done' --port 80 --expose
<!-- kubectl get pods -l tier=webapp -o jsonpath='{range .items[*]}{.status.podIP}{"\n"}{end}' -->
kubectl create deployment frontend --image nginx --replicas=2
kubectl scale deployment frontend --replicas=3
kubectl set image deployment/frontend nginx=nginx:1.18 --record
kubectl expose deployment nginx --port 80
kubectl expose deployment webapp --name=webapp-service --port=8080 --target-port=8080 --type=NodePort --dry-run=client -o yaml

kubectl label nodes node01 size=large
curl -v -k -u '<username>:<password>' https://master-loadbalancer:6443/api/v1/pods
curl -v -k -H 'Authorization: Bearer <token>' https://master-loadbalancer:6443/api/v1/pods
curl -v -k https://master-loadbalancer:6443/api/v1/pods --key admin.key --cert admin.crt --cacert ca.crt

kubectl get csr
kubectl csr approve john
kubectl get csr john -o yaml
kubectl certificate deny rogue-req

kubectl proxy; curl -k http://localhost:6443/apis;

kubectl create role engineer --verb=create --verb=get --verb=list --verb=delete --resource=pods --dry-run=client -o yaml
kubectl create rolebinding engineering-user-role-binding --role=engineer --user=engineering-user --dry-run=client -o yaml

kubectl auth can-i create deployments
kubectl auth can-i delete pods
kubectl auth can-i delete pods --as john
kubectl auth can-i delete pods --as john --namespace prod

kubeadm token list
apt-get -y -qq update; apt install -y kubeadm=1.20.0-00; kubeadm upgrade plan; kubeadm upgrade apply v1.20.0; kubectl get nodes; kubectl drain controlplane --ignore-daemonsets; apt-get install -y kubelet=1.20.0-00; apt-get install -y kubectl=1.20.0-00; systemctl daemon-reload; systemctl restart kubelet; kubectl uncordon controlplane;
kubectl drain node01 --ignore-daemonsets
ssh node01; apt-get install -y kubeadm=1.20.0-00; apt-get install -y kubelet=1.20.0-00; kubeadm upgrade node; systemctl daemon-reload; systemctl restart kubelet
kubectl uncordon node01

ETCDCTL_API=3 etcdctl cluster-health --endpoints=https://127.0.0.1:2379
ETCDCTL_API=3 etcdctl endpoint health --endpoints=https://127.0.0.1:2379
ETCDCTL_API=3 etcdctl snapshot save snapshot.db --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=https://127.0.0.1:2379
ETCDCTL_API=3 etcdctl snapshot status snapshot.db --cert=/etcd/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=https://127.0.0.1:2379
systemctl kube-apiserver stop
ETCDCTL_API=3 etcdctl snapshot restore snapshot.db --data-dir=/var/lib/etcd-from-backup --cert=/etc/kubernetes/pki/etcd/server.crt --key=/etc/kubernetes/pki/etcd/server.key --cacert=/etc/kubernetes/pki/etcd/ca.crt
sed -i -e 's#--data-dir=/var/lib/etcd#--data-dir=/var/lib/etcd-from-backup#g' /etc/systemd/system/etcd.service
systemctl daemon-reload; systemctl etcd restart; systemctl kube-apiserver restart

kubectl create configmap app-config --from-literal=COLOR=blue --from-literal=ENVIRONMENT=staging
kubectl create configmap app-config --from-file=app-config.properties

cat /etc/systemd/system/kube-apiserver.service
ps -aux | grep kube-apiserver
cat /etc/kubernetes/manifests/kube-apiserver.yaml
cat /etc/kubernetes/manifests/kube-apiserver.yaml | grep -i '--authorization-mode'

cat /etc/systemd/system/kube-controller-manager.service
ps -aux | grep kube-controller-manager
cat /etc/kubernetes/manifests/kube-controller-manager.yaml | grep -i '--controllers'

cat /etc/systemd/system/kubelet.service | grep '--config'
cat /var/lib/kubelet/kubelet-config.yaml | grep '--pod-manifest-path'

kubectl get events | grep -i <custom-scheduler-name>

<!-- kubectl get nodes -o jsonpath='{range .items[*]}{.status.nodeInfo.osImage}{"\n"}{end}' -->
<!-- kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.cpu}{"\n"}{end}' -->
kubectl get nodes -o custom-columns=NODE:.metadata.name,CPU:.status.capacity.cpu
kubectl get pv -o custom-columns=NAME:.metadata.name,CAPACITY:.spec.capacity.storage --sort-by=.spec.capacity.storage
kubectl config view --kubeconfig=/path/to/kubeconfig -o jsonpath='{$.contexts[?(@.context.user == "aws-user")].name}{"\n"}'
<!-- kubectl get nodes -o jsonpath='{range .items[*]}{range .status.addresses}{"InternalIP of "}{@[1].address}{" "}{@[0].address}{" "}{end}{end}' > /root/CKA/node_ips -->

kubectl create quota dev-ns-count --hard=count/deployments.apps=2,count/replicasets.apps=4,count/pods=10,count/secrets=4 --namespace=dev

kubectl taint node node01 color=green:NoExecute
kubectl taint node node01 color=green:NoExecute-
kubectl explain pod --recursive | grep -A5 tolerations
kubectl taint node controlplane node-role.kubernetes.io/master:NoSchedule-

kubectl logs web --previous
kubectl exec cassandra -- cat /var/log/cassandra/system.log

kubectl -n kube-system exec etcd-master -- /bin/sh -c "ETCDCTL_API=3 etcdctl get / --prefix --keys-only --limit=10 --key=/etc/kubernetes/pki/etcd/server.key --cert=/etc/kubernetes/pki/etcd/server.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt"
kubectl -n kube-system exec etcd-controlplane -c etcd -- /bin/sh -c 'ETCDCTL_API=3 etcdctl endpoint health --key=/etc/kubernetes/pki/etcd/server.key --cert=/etc/kubernetes/pki/etcd/server.crt --cacert=/etc/kubernetes/pki/etcd/ca.crt --endpoints=https://127.0.0.1:2379'

kubectl describe ingress ingress-online-store

kubectl -n kube-system get deployment coredns -o yaml
kubectl -n kube-system get svc kube-dns
kubectl -n kube-system get configmap coredns -o yaml | grep -A 15 '.:53'
cat /var/lib/kubelet/config.yaml | grep -C5 -i -E '(clusterDNS|clusterDomain)'
kubectl run --rm --restart=Never busybox --image=busybox --command cat /etc/resolv.conf

kubectl -n kube-system edit system:kube-scheduler
kubectl get events | grep -i 'scheduled'
kubectl exec busybox -it -- nslookup nginx-resolver-service
kubectl exec busybox -it -- nslookup 10-244-1-7.default.pod.cluster.local
kubectl expose pod nginx-resolver --name=nginx-resolver-service --port=80 --target-port=80

kubectl describe pod kube-apiserver-controlplane
kubectl get pods -n kube-system | grep -i dns
kubectl get svc kube-dns -n kube-system
kubectl get pods -n kube-system -l k8s-app=kube-dns
kubectl describe configmap coredns -n kube-system
kubectl -n kube-system describe deployment coredns | grep -C3 -i args | grep -i corefile
kubectl exec test-pod -c test -it -- nc -v -z web-service 80
kubectl exec test-pod -c test -it -- nslookup <another_pod_ip>
kubectl exec hr -c web -it -- nslookup mysql.payroll

kubectl run curl --image=alpine/curl --rm -it -- sh

kubectl run ubuntu --image ubuntu --overrides='{"apiVersion": "v1", "spec": {"template": {"spec": {"nodeSelector": {"kubernetes.io/hostname": "node01"}}}}}' --restart=Never --command sleep 300
kubectl get pods --show-labels
kubectl get pods --selector app=nginx
kubectl get pods -l app=nginx
kubectl get pod nginx --watch
kubectl get pods --no-headers --selector env=prod,bu=finance,tier=frontend


kubectl logs -n kube-system <weave-net-pod> -c weave | grep -i range

for i in {1..35}; do
   kubectl exec --namespace=kube-public curl -- sh -c 'test=`wget -qO- -T 2  http://webapp-service.default.svc.cluster.local:8080/info 2>&1` && echo "$test OK" || echo "Failed"';
   echo ""
done

kubectl exec busybox -- /bin/bash -c 'until nslookup db-service; do echo waiting for db to be up and running; sleep 3; done;';
kubectl exec nikola/netshoot -- /bin/bash -c 'while true; do echo -e "HTTP/1.1 200 OK\n SUCCESS" | nc -l -p 80 -q 1; done';
kubectl set image deployment/nginx nginx=nginx:1.18 --record
kubectl explain resourcequota
kubectl create quota dev-ns-counts --hard=count/deployments.apps=2,count/replicasets.apps=4,count/pods=10,count/secrets=4 --namespace=dev
kubectl get pods -o jsonpath='{.items[0].spec.containers[*].image}'
kubectl get nodes -o jsonpath='{range .items[*]}{.metadata.name}{"\t"}{.status.capacity.cpu}{"\n"}{end}'
kubectl get persistentvolume --sort-by=.spec.capacity.storage -o custom-columns=NAME:.metadata.name,CAPACITY:.spec.capacity.storage
kubectl taint node controlplane node-role.kubernetes.io/master:NoSchedule-
kubectl explain limitrange.spec --recursive
kubectl rollout status deployment/webapp-deployment
kubectl rollout history deployment/webapp-deployment
kubectl rollout undo deployment/webapp-deployment
kubectl rollout undo deployment/frontend --to-revision=2
kubectl rollout history deployment/frontend
kubectl annotate pods my-pod color=blue
kubectl autoscale deployment foo --min=2 --max=10
kubectl drain node01 --ignore-daemonsets
kubectl drain node01 --ignore-daemonsets --force
kubectl cordon node01
kubectl uncordon node01
kubectl get all --all-namespaces -o yaml > everything.yaml
kubectl get pods --server <kubernetes-api-server-url>:6443 --client-key kube-admin.key --client-certificate kube-admin.crt --certificate-authority ca.crt
kubectl config use-context engineer@kubernetes-cluster
kubectl cp /tmp/foo_dir my-pod:/tmp/bar_dir
kubectl cp my-namespace/my-pod:/tmp/foo /tmp/bar

kubectl get pods -l app=blue -o=custom-columns='NAME:.metadata.name' --no-headers=true

kubectl --kubeconfig /etc/kubernetes/kubelet.conf run nginx-critical --image=nginx --dry-run=client -o yaml > /etc/kubernetes/manifests/nginx-critical.yaml

kubectl proxy && curl http://localhost:6443 -k
kubectl proxy && curl http://localhost:6443/apis -k

kubectl create role engineer --verb=list --verb=create --verb=delete --resource=pods --dry-run=client -o yaml
kubectl get roles
kubectl get rolebindings
kubectl describe role engineer
kubectl create rolebinding engineering-user-role-binding --role=engineer --user=engineering-user --dry-run=client -o yaml
kubectl describe rolebinding engineering-user-role-binding

kubectl create clusterrole pvviewer-role --verb=list --resource=persistentvolumes
kubectl create clusterrolebinding pvviewer-role-binding --clusterrole=pvviewer-role --serviceaccount=default:pvviewer

kubectl auth can-i create deployments
kubectl auth can-i delete nodes
kubectl auth can-i create deployments --as john
kubectl auth can-i create pods --as john --namespace test
kubectl api-resources --namespaced=true
kubectl api-resources --namespaced=false
kubectl describe pod kube-apiserver-controlplane -n kube-system | grep -C 7 -i -E '(authorization-mode)'
kubectl explain pod.spec.containers.securityContext --recursive | grep -i privilege
kubectl patch svc http-svc -p '{"spec":{"type": "ClusterIP"}}'
kubectl patch svc http-svc -p '{"spec":{"type": "NodePort"}}'

kubectl annotate serviceaccount ebs-csi-controller-sa -n kube-system eks.amazonaws.com/role-arn=arn:aws:iam::YOUR_AWS_ACCOUNT_ID:role/AmazonEKS_EBS_CSI_DriverRole

kubectl proxy 8001&
curl localhost:8001/apis/networking.k8s.io
kubectl convert -f ingress-old-api-spec.yaml --output-version networking.k8s.io/v1 | kubectl apply -f -

aws eks --region us-west-2 update-kubeconfig --name <cluster_name>
aws ssm start-session --target <instance_id>
METADATA_API_TOKEN="$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600")"
curl -H "X-aws-ec2-metadata-token: $METADATA_API_TOKEN" http://169.254.169.254/latest/meta-data/
aws-iam-authenticator token -i web-k8s-eks --region us-west-2
cat /etc/kubernetes/kubelet/kubelet-config.json
kubectl describe daemonset aws-node -n kube-system
kubectl port-forward pods/<pod-name> 28015:27017
journalctl -u kubelet
journalctl -u kube-proxy
```
