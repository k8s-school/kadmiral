. ~/.novacreds/fink-openrc.sh
. ~/openstack_cli/bin/activate

# Create bastion

openstack keypair create --public-key ~/.ssh/id_rsa.pub fjammes-key
openstack server create --image "official-ubuntu-24.04-x86_64"   --flavor m1.small   --network fink-public   --network fink  --security-group talos   --key-name fjammes-key talos-bastion
BASTION_IP=$(openstack server show talos-bastion -f json -c addresses | jq -r '.addresses."fink-public"[0]')
ssh ubuntu@$BASTION_IP "mkdir /home/ubuntu/.novacreds /home/ubuntu/.talos/"
scp $HOME/.novacreds/fink-openrc.sh ubuntu@$BASTION_IP:/home/ubuntu/.novacreds/
scp $HOME/.talos/config  ubuntu@$BASTION_IP:/home/ubuntu/.talos/
sudo apt update && sudo apt install -y python3-openstackclient
git clone https://github.com/k8s-school/kadmiral.git

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
chmod +x talosctl
mv talosctl ./bin/

curl -LO https://dl.k8s.io/release/v1.35.0/bin/linux/amd64/kubectl
chmod +x kubectl
mv kubectl ./bin/

talosctl gen config fink-cluster https://10.180.15.250:6443

talosctl machineconfig patch controlplane.yaml --patch @patch-vip.yaml -o controlplane-vip.yaml
talosctl machineconfig patch controlplane-vip.yaml --patch @patch-ptp.yaml -o controlplane-ptp.yaml

mv controlplane-ptp.yaml ontrolplane-final.yaml

talosctl config merge ./talosconfig
talosctl config endpoint 10.180.15.250
talosctl config node 10.180.15.250
talosctl config context fink-cluster
cat ~/.talos/config


## Security group

openstack security group create talos --description "Security group for Talos cluster (VIP 10.180.15.250)"

# Talos API: Required for talosctl (bootstrap, upgrade, dashboard)
openstack security group rule create --protocol tcp --dst-port 50000 --remote-ip 0.0.0.0/0 talos

# Kubernetes API: Required for kubectl and cluster interaction
openstack security group rule create --protocol tcp --dst-port 6443 --remote-ip 0.0.0.0/0 talos

# Allow all TCP traffic within the fink subnet
openstack security group rule create --protocol tcp --remote-ip 10.180.15.0/24 talos

# Allow all UDP traffic within the fink subnet (required for DNS and VXLAN/CNI)
openstack security group rule create --protocol udp --remote-ip 10.180.15.0/24 talos

# Allow ICMP (Ping) for network diagnostics and troubleshooting
openstack security group rule create --protocol icmp --remote-ip 10.180.15.0/24 talos

# VRRP Protocol (112): Essential for the Virtual IP (10.180.15.250) failover mechanism
openstack security group rule create --protocol 112 --remote-ip 10.180.15.0/24 talos

# ICMP
openstack security group rule create --proto icmp talos

# SSH
openstack security group rule create --proto tcp --dst-port 22 talos

openstack security group rule list talos

## Create vm

for i in $( seq 1 3 ); do
  openstack server create talos-control-plane-$i --flavor m1.small --image talos --network fink --security-group talos --user-data /path/to/controlplane.yaml
done

openstack console log show talos-master-1

CONTROL_PLANE_IP=$(openstack server show talos-master-1 -f json -c addresses | jq -r '.addresses.fink[0]')
talosctl config endpoint $CONTROL_PLANE_IP
talosctl config node $CONTROL_PLANE_IP