. ~/.novacreds/fink-openrc.sh
. ~/openstack_cli/bin/activate

## Create talos image

```
curl -LO https://factory.talos.dev/image/376567988ad370138ad8b2698212367b8edcb69b5fd68c80be1f2ec7d603b4ba/v1.12.5/openstack-amd64.raw.xz
xz -d openstack-amd64.raw.xz
openstack image create --public --disk-format raw --file openstack-amd64.raw talos
```

## Create VIP

### Fail to create custom network

# openstack network create talos-net

### Us fink network

openstack port create --network fink --fixed-ip ip-address=10.180.15.250 talos-vip-port

## Install talosctl

curl -Lo talosctl https://github.com/siderolabs/talos/releases/download/v1.12.6/talosctl-linux-amd64
