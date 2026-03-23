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

VIP_IP=$(openstack port show "$VIP_PORT_NAME" -f json -c fixed_ips | jq -r '.fixed_ips[0].ip_address')
if [ -z "$VIP_IP" ] || [ "$VIP_IP" == "null" ]; then
    echo "ERROR: Could not find IP for $VIP_PORT_NAME. check if 'jq' is installed."
    exit 1
fi

echo "Found VIP: $VIP_IP"

for i in $(seq 1 3); do
    VM_NAME="talos-control-plane-$i"
    echo "--- Creating $VM_NAME ---"

    # 1. Create the server
    # Added --config-drive true to help Talos find metadata in OpenStack
    openstack server create "$VM_NAME" \
      --flavor "$FLAVOR" \
      --image "$IMAGE" \
      --network "$NETWORK" \
      --security-group "$SEC_GROUP" \
      --user-data "$CONFIG_PATH"

    echo "Waiting a few seconds for port creation..."
    sleep 5

    # 2. Find the Port ID for this specific VM
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
