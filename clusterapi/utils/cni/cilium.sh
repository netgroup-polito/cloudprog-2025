#!/usr/bin/env bash

function install_cilium() {
    local kubeconfig=$1
    local POD_CIDR=$2

    cat <<EOF > cilium-values.yaml
ipam:
  operator:
    clusterPoolIPv4PodCIDRList: ${POD_CIDR}

affinity:
  nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: liqo.io/type
            operator: DoesNotExist
encryption:
  enabled: true
  type: wireguard

EOF

    KUBECONFIG="$kubeconfig" cilium install --values "cilium-values.yaml" --wait
    rm cilium-values.yaml
}

function wait_cilium() {
    local kubeconfig=$1
    KUBECONFIG="$kubeconfig" cilium status --wait
}