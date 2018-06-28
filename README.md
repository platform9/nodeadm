# nodeadm
Kubernetes node administration tool
## Usage
```
nodeadm init --cfg=<location to config file>
```
## Example config
```
vipConfiguration:
  IP: 192.168.96.75
  RouterID: 42
  NetworkInterface: eth0
masterConfiguration:
  api:
    advertiseAddress: 192.168.96.75
    bindPort: 443
  apiServerCertSANs:
  - 192.168.96.17
  - 192.168.96.75
  clusterName: test
  etcd:
    caFile: "<location of ca crt>"
    certFile: "<location of client cert file>"
    keyFile: "<location of client key file>"
    endpoints:
    - https://127.0.0.1:2379
  networking:
    podSubnet: 10.1.0.0/16
    serviceSubnet: 10.2.0.0/16
```
