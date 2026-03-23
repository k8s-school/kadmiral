#!/bin/bash

set -euxo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Configuration
VIP_PORT_NAME="talos-vip-port"
IMAGE="talos"
FLAVOR="m1.small"
NETWORK="fink"
SEC_GROUP="talos"
CONFIG_PATH="$DIR/controlplane-final.yaml"
CP_NODE_PREFIX="talos-control-plane"
WORKER_NODE_PREFIX="talos-worker"

VIP_IP=$(openstack port show "$VIP_PORT_NAME" -f json -c fixed_ips | jq -r '.fixed_ips[0].ip_address')
if [ -z "$VIP_IP" ] || [ "$VIP_IP" == "null" ]; then
    echo "ERROR: Could not find IP for $VIP_PORT_NAME. check if 'jq' is installed."
    exit 1
fi

echo "Found VIP: $VIP_IP"

for i in $(seq 1 3); do
    VM_NAME="$CP_NODE_PREFIX-$i"
    echo "--- Creating $VM_NAME ---"

    # Create the server
    openstack server create "$VM_NAME" \
      --flavor "$FLAVOR" \
      --image "$IMAGE" \
      --network "$NETWORK" \
      --security-group "$SEC_GROUP" \
      --user-data "$CONFIG_PATH"

    echo "Waiting a few seconds for port creation..."
    sleep 5

    # Find the Port ID for this specific VM
    PORT_ID=$(openstack port list --server "$VM_NAME" -c ID -f value)

    if [ -n "$PORT_ID" ]; then
        echo "Found Port ID: $PORT_ID. Authorizing VIP $VIP_IP..."
        # 3. Apply the allowed-address-pair
        openstack port set --allowed-address ip-address="$VIP_IP" "$PORT_ID"
        echo "SUCCESS: $VM_NAME is ready for High Availability."
    else
        echo "WARNING: Could not find port for $VM_NAME. You might need to run the port set command manually."
    fi
done

CP_NODE_1_IP=$(openstack server show "$CP_NODE_PREFIX-1" -c addresses -f value | tr -d "{}'[] " | cut -d: -f2)

if [ -z "$CP_NODE_1_IP" ]; then
    echo "ERROR: Could not find IP for $CP_NODE_PREFIX-1. check if 'jq' is installed."
    exit 1
fi

# Configure talosctl to point to the first control plane node
talosctl config endpoint "$CP_NODE_1_IP"
talosctl config node "$CP_NODE_1_IP"

# Bootstrap the cluster using the first control plane node
talosctl bootstrap

# Create kubeconfig for the cluster
talosctl kubeconfig

talosctl config endpoint "$VIP_IP"
talosctl config node "$VIP_IP"

talosctl machineconfig patch $DIR/worker.yaml --patch @ptp.yaml -o $DIR/worker-ptp.yaml
cp $DIR/worker-ptp.yaml $DIR/worker-final.yaml

for i in $(seq 1 2); do
    openstack server create "$WORKER_NODE_PREFIX-$i" --flavor "$FLAVOR" --network "$NETWORK" --image "$IMAGE" --security-group "$SEC_GROUP" --user-data $DIR/worker-final.yaml
done

