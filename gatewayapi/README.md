# Gateway API Deployment with Kubernetes

This project demonstrates how to deploy an application using Kubernetes Gateway API.

## Installing Contour

Contour is a high-performance ingress controller that supports the Gateway API. Follow these steps to install Contour:

1. **Create a KinD cluster**:

   ```bash
   kind create cluster --name gatewayapidemo --config kind-config.yml
   ```

2. **Install a CNI**:

   ```bash
   cilium install --wait
   ```

3. **Install Contour Gateway Provisioner**:

   ```bash
   kubectl apply -f https://projectcontour.io/quickstart/contour-gateway-provisioner.yaml
   ```
