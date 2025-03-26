# Kustomize Manifests for Datum APIServer and Controller Manager

## Overview
This repository provides Kustomize manifests to deploy the Datum APIServer and Datum Controller Manager components. These manifests are structured for ease of use by the community and integration with FluxCD pipelines.

## Repository Structure
```
config
├── api-server
│   ├── deployment.yaml
│   ├── httpproxy.yaml
│   ├── kustomization.yaml
│   ├── service.yaml
├── controller-manager
│   ├── deployment.yaml
│   ├── kustomization.yaml
docs
```

### API Server
The `api-server` folder contains the Kustomize manifests required to deploy the Datum APIServer, including:
- **deployment.yaml**: Defines the Kubernetes Deployment for the API Server.
- **httpproxy.yaml**: Configuration for HTTP routing (if applicable).
- **kustomization.yaml**: Kustomize configuration for managing API Server resources.
- **service.yaml**: Defines the Kubernetes Service for the API Server.

### Controller Manager
The `controller-manager` folder contains the Kustomize manifests required to deploy the Datum Controller Manager, including:
- **deployment.yaml**: Defines the Kubernetes Deployment for the Controller Manager.
- **kustomization.yaml**: Kustomize configuration for managing Controller Manager resources.

## Pushing Manifests using Flux CLI
We utilize `flux push artifact` to publish Kustomize manifests to an OCI repository.

### Example Workflow
To push the manifests to an OCI registry, use the following command:
```sh
flux push artifact oci://ghcr.io/your-org/datum-kustomize:latest \
  --path=./config --source=your-repository-url
```

## GitHub Actions Integration
A GitHub Action is set up to automatically push these manifests upon changes. The workflow is defined as follows:

```yaml
name: Publish Kustomize Manifests

on:
  push:
    branches:
      - main
  release:
    types: [published]

jobs:
  push-kustomize:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3
      
      - name: Install Flux CLI
        run: |
          curl -s https://fluxcd.io/install.sh | sudo bash
      
      - name: Push Manifests
        run: |
          flux push artifact oci://ghcr.io/your-org/datum-kustomize:latest \
            --path=./config --source=\$(git remote get-url origin)
```

This ensures that any updates to the `config` directory are automatically pushed to the OCI registry.

This setup enables both community users and internal automation (e.g., FluxCD) to deploy the Datum APIServer and Controller Manager efficiently.

