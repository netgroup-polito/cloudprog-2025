#!/usr/bin/env bash

if [ -z "$1" ]; then
    echo "Usage: $0 <tot-vm-number>"
    exit 1
fi

TOT_VM_NUM=$1
SSH_PUB_KEY=$(cat ~/.ssh/id_rsa.pub)
export SSH_PUB_KEY

for((i=1;i<=TOT_VM_NUM;i++)); do
    export VM_NUM="$i"
    echo "Creating VM ${VM_NUM}"
    export VM_NUM="$i"
    envsubst < ./vm.yml > tmp
    kubectl apply -f tmp
done

rm tmp