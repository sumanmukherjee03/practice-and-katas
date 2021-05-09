## persistent volume

Docker stores data in the host filesystem within

/var/lib/docker
  - aufs
  - containers
  - image
  - volumes

In docker there are 2 kinds of volumes - volume mount and bind mount.
"Volume mount" mounts a created volume in docker into the docker container
where as a "bind mount" mounts a host directory into the docker container.

```
docker volume create data_volume
docker run -v data_volume:/var/lib/mysql mysql
docker run -mount type=volume,src=data_volume,dst=/var/lib/mysql mysql
```

OR

```
docker run -v /data/mysql:/var/lib/mysql mysql
docker run -mount type=bind,src=/data/mysql,dst=/var/lib/mysql mysql
```

The thing that manages the layered file system architecture, moving layers, copy on write for layers etc is storage drivers.
Different docker storage drivers available for docker are :
  - aufs : default for ubuntu
  - device mapper : could be the default for centos, rhel
  - btrfs
  - overlay
  - overlay2
  - zfs

The creation and management of volumes is not handled by storage drivers but volume driver plugins.
Different volume drivers available for docker are :
  - local : for creating volumes on the docker host in /var/lib/docker/volumes
  - azure file storage
  - convoy
  - flocker
  - gce-docker
  - RexRay : can be used to provision aws storage in aws ebs, s3 etc.

```
docker run -it --name mysql --volume-driver rexray/ebs --mount src=ebs-vol, target=/var/lib/mysql mysql
```

Volumes maintain state even after the docker container has exited.

Similar to Container Runtime Interface (CRI) or Container Network Interface (CNI), kubernetes also has Container Storage Interface (CSI).
Different vendors implement the CSI for allowing pods to store data in the cloud.

------------------------------------------------------

Below is an example of a simple `pod-definition.yaml` file that stores data in the hosts' `/data` directory.
Inside the container this volume is being bind mounted to `/opt` directory.
```
apiVersion: v1
kind: Pod
metadata:
  name: rand-num-gen
spec:
  containers:
    - image: alpine
      name: rand-num-gen-container
      command: ["/bin/sh", "-c"]
      args: ["shuf -i 0-100 -n 1 >> /opt/numbers.out"]
      volumeMounts:
        - mountPath: /opt
          name: data-volume
  volumes:
    - name: data-vol
      hostPath:
        path: /data
        type: Directory
```
This form of volume mount isnt useful in a cluster though because the directory is not being shared across multiple hosts.

To fix this in aws for instance you could store the data in ebs volumes by doing something like this.
Remember however that ebs volumes are zone specific. In kubernetes that is accounted for with allowed topologies.
```
apiVersion: v1
kind: Pod
metadata:
  name: rand-num-gen
spec:
  containers:
    - image: alpine
      name: rand-num-gen-container
      command: ["/bin/sh", "-c"]
      args: ["shuf -i 0-100 -n 1 >> /opt/numbers.out"]
      volumeMounts:
        - mountPath: /opt
          name: data-volume
  volumes:
    - name: data-vol
      awsElasticBlockStore:
        volumeID: vol-1234567890
        fsType: ext4
        type: gp2
```

To get away from the users having to define storage in pod definition files and manage storage more centrally
persistent volumes is used. Users can then carve out a chunk of storage from these centrally managed storage.
The access modes are:
  - ReadWriteOnce : the volume can be mounted as read-write by a single node
  - ReadOnlyMany : the volume can be mounted read-only by many nodes
  - ReadWriteMany : the volume can be mounted as read-write by many nodes

For example `cat pv-definition.yaml`
```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv-vol1
spec:
  accessModes:
    - ReadWriteMany
  capacity:
    storage: 10Gi
  awsElasticBlockStore:
    volumeID: vol-1234567890
    fsType: ext4
    type: gp2
  persistentVolumeReclaimPolicy: Retain
```
To create the persistent volume `kubectl create -f pv-definition.yaml`

The persistent volume is not deleted by default when the claim is deleted.
However that can be controlled with the `persistentVolumeReclaimPolicy` field.
The valid values are `Retain`, `Recycle`, `Delete`. The Recycle is an useful option as it wipes out the data
and frees up the volume again to be reclaimed.

Useful commands to inspect persistent volumes
```
kubectl get persistentvolume
kubectl describe persistentvolume pv-vol1
```

`cat pvc-definition.yaml`
```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: pv-vol1-claim
spec:
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 500Mi
```
To create the persistent volume claim `kubectl create -f pvc-definition.yaml`
Useful commands to inspect persistent volume claims
```
kubectl get persistentvolumeclaim
kubectl describe persistentvolumeclaim pv-vol1-claim
kubectl delete persistentvolumeclaim pv-vol1-claim
```

Persistent volume claims can be used in pods like so :
```
apiVersion: v1
kind: Pod
metadata:
  name: webapp
spec:
  containers:
    - name: webapp
      image: backend-webapp:latest
      volumeMounts:
      - mountPath: "/opt/data"
        name: webapp-vol
  volumes:
    - name: webapp-vol
      persistentVolumeClaim:
        claimName: pv-vol1-claim
```
