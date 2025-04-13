# Gateway API Deployment with Kubernetes

This project demonstrates how to deploy an application using Kubernetes Gateway API.

## Installing Contour

Contour is a high-performance ingress controller that supports the Gateway API. Follow these steps to install Contour:

1. **Create a KinD cluster**:
   ```bash
   kind create cluster
   ```

2. **Install Contour Gateway Provisioner**:
   Apply the Contour Gateway Provisioner API manifests:
   ```bash
   kubectl apply -f https://projectcontour.io/quickstart/contour-gateway-provisioner.yaml
   ```
