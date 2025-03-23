#!/usr/bin/env bash

set -e           # Fail in case of error
set -o nounset   # Fail if undefined variables are used
set -o pipefail  # Fail if one of the piped commands fails

error() {
   local sourcefile=$1
   local lineno=$2
   echo "An error occurred at $sourcefile:$lineno."
}
trap 'error "${BASH_SOURCE}" "${LINENO}"' ERR

FILEPATH=$(realpath "$0")
WORKDIR=$(dirname "$FILEPATH")

# shellcheck disable=SC1091
# shellcheck source=./utils/utils.sh
source "$WORKDIR/utils/utils.sh"

# shellcheck disable=SC1091
# shellcheck source=./utils/cni/cni.sh
source "$WORKDIR/utils/cni/cni.sh"

DOCKER_PROXY="${DOCKER_PROXY:-docker.io}"

function createcluster () {
    index=$1
    cni=$2
    podcidrtype=$3
    template=$4
    image=$5

    name=$(forgename "${index}" "${cni}" "${podcidrtype}" "${template}" "${image}")

    if kubectl get "clusters.cluster.x-k8s.io/${name}" -n cloudprog-demo &> /dev/null; then
        echo "Cluster ${name} already exists"
        clusterctl get kubeconfig -n cloudprog-demo "${name}" > "${HOME}/${name}"
        return
    fi

    export NODE_VM_IMAGE_TEMPLATE="${image}"

    podcidr=""
    if [ "${podcidrtype}" == "over" ]; then
        podcidr="10.80.0.0/16"
    else
        podcidr="10.8${index}.0.0/16"
    fi
    export POD_CIDR="${podcidr}"

    clusterctl generate cluster "${name}" \
        --kubernetes-version v1.30.3 \
        --control-plane-machine-count=1 \
        --worker-machine-count=2 \
        --from "${template}" \
        --target-namespace cloudprog-demo | kubectl apply -f -
    
    echo "Waiting for cluster ${name} to be ready"
    kubectl wait --for condition=Ready=true -n cloudprog-demo "clusters.cluster.x-k8s.io/${name}" --timeout=-1s

    echo "Getting kubeconfig for cluster ${name}"
    clusterctl get kubeconfig -n cloudprog-demo "${name}" > "${HOME}/${name}"

    echo "Installing CNI ${cni} on cluster ${name}"
    installcni "${HOME}/${name}" "${cni}" "${POD_CIDR}"

    echo "Installing local-path-provisioner on cluster ${name}"
    kubectl apply -f https://raw.githubusercontent.com/rancher/local-path-provisioner/v0.0.24/deploy/local-path-storage.yaml --kubeconfig "${HOME}/${name}"
    kubectl annotate storageclass local-path storageclass.kubernetes.io/is-default-class=true --kubeconfig "${HOME}/${name}"

    echo "Cluster ${name} ready"
}


CLUSTER_NUM=1

images=(
    harbor.crownlabs.polito.it/capk/ubuntu-2204-container-disk:v1.30.3
    #harbor.crownlabs.polito.it/capk/rockylinux-9-container-disk:v1.30.3
)

templates=(
    ./capi-templates/kubeadm.yaml
)

cnis=(
    cilium
    #calico
    #flannel
)

podcidrtypes=(
    #over
    dist
)

for podcidrtype in "${podcidrtypes[@]}"; do
    for template in "${templates[@]}"; do
        for image in "${images[@]}"; do
            for cni in "${cnis[@]}"; do
                for i in $(seq 1 $CLUSTER_NUM); do
                    createcluster "${i}" "${cni}" "${podcidrtype}" "${template}" "${image}"
                done
            done
        done
    done
done


