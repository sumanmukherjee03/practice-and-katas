## networking basics

### Route tables

Get the interfaces of the machine
`ip link`

If the switch is on network `192.168.1.0` and you want to connect a machine to it.
First assign eth0 on your machine an IP address from that same network.
`ip addr add 192.168.1.10 dev eth0`

Say you did this for another machine
`ip addr add 192.168.1.11/24 dev eth0`

To persist these settings across reboots update the settings in `/etc/network/interfaces` file.

Now these 2 machines are on the same network connected via a switch.
So from machine `192.168.1.10` you can connect to `192.168.1.11`.
`ping 192.168.1.11`

Switches help traffic flow within the same network.
But if a machine from `192.168.1.0/24` network wanted to connect to another machine on `192.168.2.0/24`
then it needs to go through a router. So a router helps with traffic flow between networks.

Since a router connects 2 networks, it's 2 different ports gets 2 different IPs assigned - one on each network.
So, lets say in this case the router gets 2 IPS - 192.168.1.1 and 192.168.2.1 .

To find the router for outgoing traffic from a network we need a gateway and routing tables.
Run the `route` command to get the routing table.
In this case for the first network, the gateway is the IP 192.168.1.1 because thats the door through which traffic leaves the network.

To add a route from first network to second network run
`ip route add 192.168.2.0/24 via 192.168.1.1`

Similarly for the second network to reach the first network, we need to add a route in there as well
`ip route add 192.168.1.0/24 via 192.168.2.1`
So, second networks gateway is `192.168.2.1`.

This means that the router has 2 IPs now 192.168.1.1 and 192.168.2.1, one on each interface eth0 and eth1.

However, if you want to reach an outside network like google, amazon etc, you need to add a default route via the gateway
`ip route add default via 192.1.1`


### IP tables explained

Example 1
--------------------

Provided below are some examples of routing tables and explanation of the routes.

`route`
Kernel IP routing table
Destination    Gateway     Genmask         Flags   Metric  Ref  Use  Iface

default        192.1.2.1   255.255.255.0   UG      0       0    0    eth0          | -> Both of these lines mean the same
0.0.0.0        192.1.2.1   255.255.255.0   UG      0       0    0    eth0          |    To reach any IP use the gateway 192.168.2.1

192.168.2.0    0.0.0.0     255.255.255.0   UG      0       0    0    eth0          | -> This means there is no need of a gateway to reach any machine in the 192.168.2.0/24 network.


Example 2
--------------------
In this example image there are 2 different routers, one for internal network comm and one for external network comm.

`route`
Kernel IP routing table
Destination    Gateway     Genmask         Flags   Metric  Ref  Use  Iface

default        192.1.2.1   255.255.255.0   UG      0       0    0    eth0          | -> This is the route to external networks
192.168.1.0    192.1.2.2   255.255.255.0   UG      0       0    0    eth0          | -> This is the route to internal networks



Example 3
----------------
`route -n`
Kernel IP routing table
Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
0.0.0.0         172.17.0.1      0.0.0.0         UG    0      0        0 ens3       | -> genmask 0.0.0.0 is used in the default gateway as the netmask
10.244.0.0      0.0.0.0         255.255.255.0   U     0      0        0 cni0       | -> gateway of 0.0.0.0 means that the host doesnt need a routing table for this network and doesnt need a gateway to reach this network
10.244.1.0      172.17.0.49     255.255.255.0   UG    0      0        0 ens3
172.17.0.0      0.0.0.0         255.255.0.0     U     0      0        0 ens3
172.18.0.0      0.0.0.0         255.255.255.0   U     0      0        0 docker0





..... Continuing on from earlier
In our example of a machine from `192.168.1.0/24` network wanting to connect to another machine on `192.168.2.0/24`
just having the routes defined is not enough for bidirectional communication. ip routes do not forward traffic by default.
Now, packet forwarding is enabled on the router through the network ip_forward setting.
`cat /proc/sys/net/ipv4/ip_forward` . If it is 0, no forward, changing to 1 will allow packet forwarding.
To persist this setting across reboots on the router, also update the `/etc/sysctl.conf` file
```
...
net.ipv4.ip_forward = 1
...
```


### DNS

DNS entries can be in `/etc/hosts` file. `cat /etc/hosts`
```
...
192.168.1.11 db
...
```
This is local entry. It is simply like an alias to an IP for referencing locally. The real DNS name of the IP
might not be what's in the /etc/hosts file. The host of that IP 192.168.1.11 could actually be called `host2.example.com`.


But managing DNS entries via /etc/hosts files can get unmanageable very quick if things span more than a couple of hosts.
Which is why we make the hosts ask a DNS server to resolve the DNS for us.
Which DNS server to ask is configured on the hosts via the `/etc/resolv.conf` file.
```
nameserver 192.168.1.101
nameserver 8.8.8.8
```

One way to think of it is that local or privately available DNS can be configured in `/etc/hosts`
whereas publicly resolvable DNS can be returned via a nameserver.

The order of name resolution is
  - first /etc/hosts file. If found, exit
  - second DNS server. If not found in the previous one, lookup in DNS server.

Of course you can have multiple nameserver entries in your `/etc/resolv.conf` file.
But that be tedious to maintain. A better solution is to have one nameserver and forward
every request it cant resolve to 8.8.8.8

Another useful directive in `/etc/resolve.conf` is the search directive.
For internal domains, like app.example.com, internally you would like to resolve `app` -> `app.example.com`,
ie, `ping app` will actually result in `ping app.example.com`. This can be achieved with the search directive.
```
nameserver 192.168.1.101
search example.com live.example.com
```
So, `app` can resolve to both `app.example.com` or `app.live.example.com`.

Tools to query DNS are nslookup and dig. Important to remember though that both of these tools dont query /etc/hosts.

A DNS server used in the kubernetes landscape is coredns. Coredns uses a config file called Corefile.
The configuration points to /etc/hosts file to load IP to hostname mappings.
When another server asks coredns to resolve a hostname it returns the IP based on the /etc/hosts file configured in the coredns server.

### Network namespaces

Create 2 network namespaces red and blue.
```
ip netns add red
ip netns add blue
ip netns
```

To view interfaces inside the network namespace. Both of these commands do the same thing.
```
ip netns exec red ip link
ip -n red link
```

Similarly for other things like arp, route etc
```
ip netns exec red arp
ip netns exec red route
```

You can connect 2 network namespaces using a veth pair, ie a veth interface on each network namespace.

Create a cable with 2 veth devices on each end
Attach veth devices to corresponding namespaces
Assign ip address to each veth device
Bring up each veth device
```
ip link add veth-red type veth peer name veth-blue
ip link set veth-red netns red
ip link set veth-blue netns blue
ip -n red addr add 192.168.15.1/24 dev veth-red
ip -n blue addr add 192.168.15.2/24 dev veth-blue
ip -n red link set veth-red up
ip -n blue link set veth-blue up
ip netns exec red ping 192.168.15.2
ip netns exec red arp
```
You can sever the connection above by deleting one end of the cable (ie, by deleting one veth device) like so
`ip -n red link del veth-red`
It automatically delete the other end of the cable, ie the veth-blue.


The above is a good way to connect 2 network namespaces. But it doesnt scale well for multiple namespaces.
Thats when you need a switch. Either the built in linux networking virtual switch Bridge or OpenVSwitch (OVS).
Bridge is the commonly used switch in docker. A bridge is like an interface device for the host and a switch for the namespaces.

Create a bridge network first
Bring up the device with the bridge network
Create a cable with 2 interfaces
Attach one end of the cable to the network namespace and other end to the bridge network
Assign ip address to the interface connected to the namespace and bring up that interface
```
ip link add v-net-0 type bridge
ip link set dev v-net-0 up
ip link add veth-red type veth peer name veth-red-br
ip link add veth-blue type veth peer name veth-blue-br
ip link set veth-red netns red
ip link set veth-red-br master v-net-0
ip link set veth-blue netns blue
ip link set veth-blue-br master v-net-0
ip -n red addr add 192.168.15.1/24 dev veth-red
ip -n blue addr add 192.168.15.2/24 dev veth-blue
ip -n red link set veth-red up
ip -n blue link set veth-blue up
ip netns exec red ping 192.168.15.2
ip netns exec red arp
```

At this point the network namespaces can all talk to each other over the brige.
However, this is a different network than the host. So, the host interface cant reach the internal network namespaces.
This is also doable. By assigning an ip address to the interface on the host for the switch, we can reach the namespaces from the host.
```
ip addr add 192.168.15.5/24 dev v-net-0
```

Now from the host you can `ping 192.168.15.1`

However this network on the bridge is internal to the host. A network namespace from the bridge network
cant ping another host on the outside. For instance say there is a host `192.168.1.4`. The network namespace
needs a gateway to do that. There are no routes for that in the network namespaces' routing tables.
The localhost itself can be a router in this case. Because one interface on the localhost has `v-net-0` and the other `eth0`
ie one interface that connects to the bridge network and one interface that connects to the external LAN network.



Add a route in the network namespace to route traffic destined for the LAN network via the bridges' `v-net-0` interface
The blue namespace can only reach the `v-net-0` because it is connected to the bridge via that interface.
The v-net-0 interface has the IP 192.168.15.5.

But for the namespace to successfully send and receive packets to the external host, we need a NAT of some kind on the host
that can mask the IP address of outgoing packets destined for the LAN with the host IP instead of the namespace IP.
This is because the external host 192.168.1.4 would only recognize and know how to respond back to 192.168.1.1 and not 192.168.15.1.

We add a NAT functionality to our host using iptables.
```
ip netns exec blue ip route add 192.168.1.0/24 via 192.168.15.5
iptables -t nat -A POSTROUTING -s 192.168.15.0/24 -j MASQUERADE
ip netns exec blue ping 192.168.1.4
```

Finally, from within the namespace you cant reach the internet because there is no default gateway.
But that can be achieved by making the host as the namesapces' default gateway because the gateway already knows how to reach the internet via the v-net-0 interface.
The v-net-0 interface has the IP 192.168.15.5.
```
ip netns exec blue ip route add default via 192.168.15.5
ip netns exec blue ping 8.8.8.8
```

How would a completely external network connect to the network namespaces.
One option is to expose the internal network to the external network's routing table to be reachable via the host.
But that's not a scalable solution.
However, if the external network can reach the host on 192.168.1.1, then another option
is to forward traffic destined to a specific port of the host to another port on the network namespace.
This can be achieved via iptables on the host
```
iptables -t nat -A PREROUTING --dport 80 --to-destination 192.168.15.2:8080 -j DNAT
```



### Docker networking

```
# Container is not accessible from even inside the host
docker run --network none nginx
# Host networking where multiple processes cant listen to the same port on the same host
docker run --network host nginx
# Third option is the bridge network where the bridge itself gets an address of 172.17.0.0 by default
docker run nginx
```

To list the docker networks
```
docker network ls
```
What shows in the output above as the `bridge` network is what shows up as `docker0` device
when you view it on the host via the `ip link` command.
`ip link add docker0 type bridge` -> docker does something similar to this for creating the bridge network.
A bridge network is like an interface to the host but like a switch to the internal network namespaces.
This docker0 interface is usually given the IP 172.17.0.1/16
When you create a docker image you can view the network namespaces created for the container via `ip netns`.
You can find the network namespace id of the container via the `NetworkSettings` > `SandboxKey`
when you run `docker inspect <containerid>`.

Docker attaches one end of the cable that connects the container to the bridge. You can find that via
`ip link | grep 'master docker0'` --> would look like `veth<containerID_small_hash>@if8`
The other end of the cable can be found on the network namespace via
`ip -n <container_network_ns_id> link` --> would look like `eth0@if7`
The interface within the container gets an internal IP. You can see it via
`ip -n <container_network_ns_id> addr` --> would look like `172.17.0.3/16`.

Once the nginx image is up and running internally you can access it via `curl http://172.17.0.3:80`
but you cant access it from outside the host. For that you need port forwarding from the host to the container.
`docker run -p 8080:80 nginx` --> listen to 8080 on the host and forward traffic to port 80 on the container.
At this point you can curl from outside the host `curl http://192.168.1.10:8080`
This is achieved internally by creating a NAT rule in the nat tables prerouting chain for iptables.
Docker does it the same way except that it adds the rule to the DOCKER chain in nat table.
`iptables -t nat -A DOCKER -j DNAT --dport 8080 --to-destination 80`.

You can see the iptable rules that docker creates via
`iptables -nvL -t nat`.
The output for docker would look like
`DNAT tcp -- anywhere        anywhere       tcp dpt:8080 to:172.17.0.2:80`
The DNAT target is used for Destination Network Address Translation, meaning rewrite destination IP address of a packet.
This is usually used to forward packets coming from outside to a firewall onto a downstream server.




### CNI

Extracting the networking overhead into common programs from different runtime engines like rocket, docker, mesos etc have proven useful for reusable code.
One such program is `bridge` and like it's name suggests it is used to attach a containers network namespace to the bridge network.

```
bridge add <container_id> <network_namespace>
bridge add 7ekjds92ldw /var/run/netns/7ekjds92ldw
```

The CNI is a set of rules that define how container runtimes needs to solve networking challenges for interoperability across container management systems.
The bridge program above is a plugin for CNI.

CNI has multiple different plugins like
  - bridge, vlan, ipvlan etc
  - IPAM plugins like dhcp, host-local
  - weave, flannel, cilium, calico etc

However, the docker runtime does not have the same standards as CNI. It is slightly different.
Which means you cant directly use docker with the CNI plugins mentioned above.

But there are work arounds.
Create a docker container with a completely isolated network, ie none.
```
docker run --network=none nginx
bridge add 4edkx78h /var/run/netns/4edkx78h
```
This is pretty much how kubernetes uses docker.


### Commands
These are some useful commands to help with debugging network problems

```
netstat -tupl
netstat -natupl
ifconfig -a
ifconfig ens3
route get <destination>
```
