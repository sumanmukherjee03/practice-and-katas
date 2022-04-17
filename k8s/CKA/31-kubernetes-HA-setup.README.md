## Kubernetes control plane HA setup

During a HA setup of the controlplane/master it is necessary to consider 3 things

  - apiserver
      The apiserver can be up and running on all master nodes.
      You however need a loadbalancer infront of the apiserver to load balance the requests from kubectl.
      This loadbalancer is what should be configured in the kubeconfig file.
      Also, the loadbalancer should be listening to 6443, the default port to which the apiserver listens.

  - controller-manager and scheduler
      The controller components like the replication controller, deployment controller etc and the scheduler needs to run as active and standby mode on the different nodes for a HA cluster.
      You cant have the scheduler or controllers active on multiple master nodes because that
      might result in creating more pods than are needed and the states wont reconcile properly.
      The active and standby mode is achieved via a leader election process.
          For example, when starting the controller-manager start with the leader elect option as true.
          This causes the kube-controller-manager to create a lock on "kubernetes-controller-manager endpoint object".
          Which ever process for the controller manager creates the lock gets the lease to act as the controller manager
          for a duration after which it can either renew the lock or if that master node goes away, a new controlplane
          node can acquire the lock during re-election and become the active kube-controller-manager.
              ```
              kube-controller-manager --leader-elect true --leader-elect-lease-duration 15s --leader-elect-renew-deadline 10s --leader-elect-retry-period 2s ...
              ```
      The scheduler has the same cli options and follows a similar approach.

  - etcd
      The etcd components can either be in the same nodes as the master nodes in a stacked topology
      OR they can live outside the master nodes in an external etcd topology. The apiserver is the only controlplane
      component that talks to the etcd servers and the kube-apiserver needs to be passed the addresses of the
      etcd nodes during startup. etcd does peer replication via raft, so it does not need a loadbalancer before itself.
      You can write state to any etcd node and it will be replicated to the other 2 etcd nodes.
      To be more precise a write request sent to a follower node is forwarded to the leader node.
      And once the leader node writes the data, it distributes a copy of the data to the follower nodes.
      A write is considered to be complete if the data has been replicated to the majority of the nodes or if the write had a quorum (n/2 + 1)

      For providing the etcd endpoints to the kube-apiserver, see the example below (etcd endpoints would be localhost and addresses of the rest of the controlplane nodes in case of a stacked deployment):
          `/usr/local/bin/kube-apiserver --etcd-servers=https://10.242.10.21:2379,https://10.242.10.22:2379,https://10.242.10.23:2379 ...`
