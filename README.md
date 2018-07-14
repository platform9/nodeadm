# nodeadm

Kubernetes node administration tool

## Usage

### Init
```
nodeadm init --cfg=/tmp/nodeadm.yaml
```

### Join
```
nodeadm join --cfg /tmp/nodeadm.yaml --master 192.168.96.75:6443 --token bootstrap.token --cahash sha256:digest
```

## Example Configuration

### Init
```
networking:
    podSubnet: 10.1.0.0/16
    serviceSubnet: 172.1.0.0/24
    dnsDomain: testcluster.local
vipConfiguration:
  IP: 192.168.96.75
  RouterID: 42
  NetworkInterface: eth0
masterConfiguration:
  api:
    advertiseAddress: 192.168.96.75
    bindPort: 443
  apiServerCertSANs:
  - 192.168.96.75
  clusterName: test
  etcd:
    caFile: /etc/etcd/pki/ca.crt
    certFile: /etc/etcd/pki/apiserver-etcd-client.crt
    keyFile: /etc/etcd/pki/apiserver-etcd-client.key
    endpoints:
    - https://127.0.0.1:2379
```

### Join
```
networking:
    podSubnet: 10.1.0.0/16
    serviceSubnet: 172.1.0.0/24
    dnsDomain: testcluster.local
```
