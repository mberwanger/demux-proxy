# Demultiplexer Proxy
Demux Proxy is a MITM proxy server that round-robins request over multiple public ip address. This proxy enables you to aggregate data from an external source where that downstream service is rate limiting based on IP address. For example, the number of IP addresses per network interface for a m5.xlarge is the following

| Instance type | Maximum network interfaces | Private IPv4 addresses per interface |
|---------------|----------------------------|--------------------------------------|
| m5.xlarge | 4 | 15 |

Therefore running demux-proxy on a m5.xlarge with 60 public addresses will increase the rate limit by that number. Since HTTP proxies are very common many libraries already have support for it and it is easy to enable.

## Running Demux Proxy

This proxy is designed to run on an EC2 instance with multiple public addresses. To assign multiple IP addresses to an instance, please refer to the [AWS Multiple IP addresses documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/MultipleIP.html).

Grab the latest binary for your platform from the [release](https://github.com/mberwanger/demux-proxy/releases) page.

```
$ demux-proxy start
{"fields":{},"level":"info","timestamp":"2021-03-08T10:56:40.086535-05:00","message":"Starting Proxy"}
```

Now send requests through the proxy and you will see that it will round-robin through the public interfaces on the host

```
$ curl --proxy http://127.0.0.1:8080 -k https://checkip.amazonaws.com
184.73.14.124
$ curl --proxy http://127.0.0.1:8080 -k https://checkip.amazonaws.com
3.94.19.116
.
.
.
```

Proxy will log request and response information 

```json
{"fields":{},"level":"info","timestamp":"2021-03-08T16:19:08.819879612Z","message":"Starting Proxy"}
{"fields":{"interface":"172.31.20.216","method":"GET","session_id":2,"type":"request","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:13.807451213Z","message":"outbound request"}
{"fields":{"method":"GET","session_id":2,"status_code":200,"type":"response","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:13.825450683Z","message":"inbound response"}
{"fields":{"interface":"172.31.16.129","method":"GET","session_id":4,"type":"request","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:16.501391662Z","message":"outbound request"}
{"fields":{"method":"GET","session_id":4,"status_code":200,"type":"response","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:16.506607494Z","message":"inbound response"}
{"fields":{"interface":"172.31.28.219","method":"GET","session_id":8,"type":"request","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.131299168Z","message":"outbound request"}
{"fields":{"interface":"172.31.25.135","method":"GET","session_id":10,"type":"request","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.425001254Z","message":"outbound request"}
{"fields":{"method":"GET","session_id":8,"status_code":200,"type":"response","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.46526242Z","message":"inbound response"}
{"fields":{"interface":"172.31.20.216","method":"GET","session_id":11,"type":"request","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.597561633Z","message":"outbound request"}
{"fields":{"method":"GET","session_id":10,"status_code":200,"type":"response","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.637055792Z","message":"inbound response"}
{"fields":{"method":"GET","session_id":11,"status_code":200,"type":"response","url":"https://checkip.amazonaws.com:443/"},"level":"info","timestamp":"2021-03-08T16:19:17.702836882Z","message":"inbound response"}
```

Note that the interface is the private ip address on the ec2 instance. For example, 

```
$ ip add
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
    inet6 ::1/128 scope host
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 0a:4a:6e:5e:20:85 brd ff:ff:ff:ff:ff:ff
    inet 172.31.20.216/20 brd 172.31.31.255 scope global dynamic eth0
       valid_lft 2700sec preferred_lft 2700sec
    inet 172.31.16.129/20 brd 172.31.31.255 scope global secondary eth0
       valid_lft forever preferred_lft forever
    inet6 fe80::84a:6eff:fe5e:2085/64 scope link
       valid_lft forever preferred_lft forever
3: eth1: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 9001 qdisc pfifo_fast state UP group default qlen 1000
    link/ether 0a:35:a2:5c:85:e9 brd ff:ff:ff:ff:ff:ff
    inet 172.31.28.219/20 brd 172.31.31.255 scope global dynamic eth1
       valid_lft 3259sec preferred_lft 3259sec
    inet 172.31.25.135/20 brd 172.31.31.255 scope global secondary eth1
       valid_lft forever preferred_lft forever
    inet6 fe80::835:a2ff:fe5c:85e9/64 scope link
       valid_lft forever preferred_lft forever
```

AWS will route the request through the public EIP IPv4 address that is associated with that private address. For context the Elastic IP addresses, allocations, and associations for this example are the following: 

| Name | Allocated IPv4 address | Type | Allocation ID | Associated instance ID | Private IP address | Association ID | Network interface owner account ID  | Network Border Group |
|------|------------------------|------|---------------|------------------------|--------------------|----------------|-------------------------------------|----------------------|
| – | 184.73.14.124 | Public IP | eipalloc-***************** | i-***************** | 172.31.20.216 | eipassoc-****************1 | ************ | us-east-1
| – | 3.94.19.116   | Public IP | eipalloc-***************** | i-***************** | 172.31.16.129 | eipassoc-****************2 | ************ | us-east-1
| – | 54.166.147.85 | Public IP | eipalloc-***************** | i-***************** | 172.31.28.219 | eipassoc-****************3 | ************ | us-east-1
| – | 67.202.4.52   | Public IP | eipalloc-***************** | i-***************** | 172.31.25.135 | eipassoc-****************4 | ************ | us-east-1
