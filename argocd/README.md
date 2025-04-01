# ArgoCD Setup Instructions

## Setup Cluster

### Run Kind Cluster

   ```bash
   kind create cluster --name argocddemo --config kind-config.yml
   ```

### Install a CNI

   ```bash
   cilium install --wait
   ```

## Install ArgoCD

### Install ArgoCD in your Kubernetes cluster

   ```bash
   kubectl create namespace argocd
   kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
   ```

### Wait for the ArgoCD pods to be ready

   ```bash
   kubectl get pods -n argocd
   ```

### Patch argocd-cm to avoid ciliumidentities out of sync

   Edit the configmap argocd-cm to disable cilium identity management. This is necessary to avoid issues with Cilium's identity management when using ArgoCD.

   ```bash
   kubectl edit configmap argocd-cm -n argocd
   ```

   Add the following lines to the configmap:

   ```yaml
   data:
      resource.exclusions: |
         - apiGroups:
           - cilium.io
           kinds:
           - CiliumIdentity
           clusters:
           - "*"
   ```

### Access the ArgoCD server

   ```bash
   kubectl port-forward svc/argocd-server -n argocd 8080:443
   ```

   Then open `https://localhost:8080` in your browser.

### Login to ArgoCD

   The default username is `admin`. Retrieve the initial password:

   ```bash
   kubectl get secret argocd-initial-admin-secret -n argocd -o jsonpath="{.data.password}" | base64 -d
   ```

## Apply the `main.yml` File

1. Ensure you are in the correct directory:

   ```bash
   cd /home/cheina/Documents/cloudprog-2025-fork/argocd
   ```

2. Apply the `main.yml` file:

   ```bash
   kubectl apply -f main.yml
   ```

3. Verify the application is running:

   ```bash
   kubectl get applications -n argocd
   ```

   Or check the ArgoCD UI at `https://localhost:8080` to see the status of your application.
