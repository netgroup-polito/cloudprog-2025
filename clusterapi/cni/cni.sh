#!/usr/bin/env bash

WORKDIR=$(dirname "$(realpath "${BASH_SOURCE[0]}")")

# shellcheck disable=SC1091
# shellcheck source=./calico.sh
source "$WORKDIR/calico.sh"

# shellcheck disable=SC1091
# shellcheck source=./cilium.sh
source "$WORKDIR/cilium.sh"

# shellcheck disable=SC1091
# shellcheck source=./flannel.sh
source "$WORKDIR/flannel.sh"

function installcni() {
    kubeconfig=$1
    cni=$2
    podcidr=$3

    case "${cni}" in
        "calico")
            install_calico "${kubeconfig}" "${podcidr}"
            wait_calico "${kubeconfig}"
            ;;
        "cilium")
            install_cilium "${kubeconfig}" "${podcidr}"
            wait_cilium "${kubeconfig}"
            ;;
        "flannel")
            install_flannel "${kubeconfig}" "${podcidr}"
            wait_flannel "${kubeconfig}"
            ;;
    esac
}