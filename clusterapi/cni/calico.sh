#!/usr/bin/env bash

function install_calico() {
    local kubeconfig=$1
    local POD_CIDR=$2
    kubectl create -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/tigera-operator.yaml --kubeconfig "$kubeconfig"

    # append a slash to DOCKER_PROXY if not present
    if [[ "${DOCKER_PROXY}" != */ ]]; then
        registry="${DOCKER_PROXY}/"
    else
        registry="${DOCKER_PROXY}"
    fi

cat <<EOF | kubectl apply -f - --kubeconfig "$kubeconfig"
# This section includes base Calico installation configuration.
# For more information, see: https://projectcalico.docs.tigera.io/master/reference/installation/api#operator.tigera.io/v1.Installation
apiVersion: operator.tigera.io/v1
kind: Installation
metadata:
  name: default
spec:
  registry: $registry
  # Configures Calico networking.
  calicoNetwork:
    # Note: The ipPools section cannot be modified post-install.
    ipPools:
    - blockSize: 26
      cidr: $POD_CIDR
      encapsulation: VXLAN
      natOutgoing: Enabled
      nodeSelector: all()
    nodeAddressAutodetectionV4:
      skipInterface: liqo.*

---

# This section configures the Calico API server.
# For more information, see: https://projectcalico.docs.tigera.io/master/reference/installation/api#operator.tigera.io/v1.APIServer
apiVersion: operator.tigera.io/v1
kind: APIServer
metadata:
  name: default
spec: {}
EOF

}

function wait_calico() {
    local kubeconfig=$1
    if ! waitandretry 5s 12 "kubectl wait --for condition=Ready=true -n calico-system pod --all --kubeconfig $kubeconfig --timeout=-1s"
    then
      echo "Failed to wait for calico pods to be ready"
      exit 1
    fi
    # set felix to use different port for VXLAN
    if ! waitandretry 5s 12 "kubectl patch felixconfiguration default --type=merge -p {\"spec\":{\"vxlanPort\":6789}} --kubeconfig $kubeconfig";
    then
      echo "Failed to patch felixconfiguration"
      exit 1
    fi
}