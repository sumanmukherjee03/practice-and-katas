CREATE cks-master VM using gcloud command.
This is not necessary if created using the browser interface
```
gcloud compute instances create cks-master --zone=europe-west3-c \
--machine-type=e2-medium \
--image=ubuntu-2004-focal-v20220419 \
--image-project=ubuntu-os-cloud \
--boot-disk-size=50GB
```

CREATE cks-worker VM using gcloud command.
This is not necessary if created using the browser interface
```
gcloud compute instances create cks-worker --zone=europe-west3-c \
--machine-type=e2-medium \
--image=ubuntu-2004-focal-v20220419 \
--image-project=ubuntu-os-cloud \
--boot-disk-size=50GB
```

You can use a region near you - https://cloud.google.com/compute/docs/regions-zones


INSTALL cks-master
```
gcloud compute ssh cks-master
sudo -i
bash <(curl -s https://raw.githubusercontent.com/killer-sh/cks-course-environment/master/cluster-setup/latest/install_master.sh)
```


INSTALL cks-worker
```
gcloud compute ssh cks-worker
sudo -i
bash <(curl -s https://raw.githubusercontent.com/killer-sh/cks-course-environment/master/cluster-setup/latest/install_worker.sh)
```

To open up firewall rules for nodeports
```
gcloud compute firewall-rules create nodeports --allow tcp:30000-40000
```

To list the compute instances
```
gcloud compute instances list
```
