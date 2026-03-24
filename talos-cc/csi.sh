#!/bin/bash

set -euxo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

git clone -b v1.35.0 https://github.com/kubernetes/cloud-provider-openstack/

curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-4
chmod 700 $DIR/get_helm.sh
$DIR/get_helm.sh

NS="csi-system"

kubectl create namespace "$NS"
kubectl label namespace "$NS" \
  pod-security.kubernetes.io/enforce=privileged \
  --overwrite

kubectl create secret generic cloud-config --from-file=cloud.conf=$DIR/cloud.conf -n "$NS"


helm repo add cpo https://kubernetes.github.io/cloud-provider-openstack
helm repo update

CHART_VERSION="2.35.0"
helm install cinder-csi cpo/openstack-cinder-csi   --namespace "$NS"   --version "$CHART_VERSION"   -f $DIR/values.yaml
