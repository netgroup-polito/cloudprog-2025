# Kyverno Installation Guide

Kyverno is a Kubernetes-native policy engine that enables you to validate, mutate, and generate configurations dynamically.

## Installation Steps

1. **Add the Kyverno Helm repository**:

   ```bash
   helm repo add kyverno https://kyverno.github.io/kyverno/
   helm repo update
   ```

2. **Install Kyverno using Helm**:

   ```bash
   helm install kyverno kyverno/kyverno --namespace kyverno --create-namespace
   ```

3. **Verify the installation**:

   ```bash
   kubectl get pods -n kyverno
   ```

   Ensure all Kyverno pods are running.
