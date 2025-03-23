#!/usr/bin/env bash

function install_flannel() {
    local kubeconfig=$1
    local POD_CIDR=$2
    kubectl create ns kube-flannel --kubeconfig "$kubeconfig"
    kubectl label --overwrite ns kube-flannel pod-security.kubernetes.io/enforce=privileged --kubeconfig "$kubeconfig"
    helm repo add flannel https://flannel-io.github.io/flannel/
    helm install flannel --set podCidr="${POD_CIDR}" --namespace kube-flannel flannel/flannel --kubeconfig "$kubeconfig"
}

function wait_flannel() {
    local kubeconfig=$1
    if ! waitandretry 5s 12 "kubectl wait --for condition=Ready=true -n kube-flannel pod --all --timeout=-1s --kubeconfig $kubeconfig";
    then
      echo "Failed to wait for flannel pods to be ready"
      exit 1
    fi
}