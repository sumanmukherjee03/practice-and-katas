## networking basics

### Route tables

Get the interfaces of the machine
`ip link`

If the switch is on network `192.168.1.0` and you want to connect a machine to it.
First assign eth0 on your machine an IP address from that same network.
`ip addr add 192.168.1.10 dev eth0`

Say you did this for another machine
`ip addr add 192.168.1.11 dev eth0`

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
