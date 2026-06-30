# Resource Manager Service Configuration

This directory contains Milo configurations for the Resource Manager service.
This extends the existing configuration for Milo's resource manager service to
customize it for Datum Cloud.

## Overview

- Organization management (Personal and Standard types)
- Project lifecycle management
- Quota enforcement for resource limits
- Validation policies for data integrity

## Structure

```
├── quota/           # Quota management
├── validation/      # Admission policies
└── kustomization.yaml
```

## Quota Limits

Quota grant policies depend on the `UnifiedOrganizations` feature gate:

- **Legacy (`UnifiedOrganizations=false`)**: apply `config/overlays/legacy-organizations`
  - Personal organizations: 2 projects maximum
  - Standard organizations: 10 projects maximum
- **Unified (`UnifiedOrganizations=true`)**: default `organization-project-quota-policy`
  - All organizations: 10 projects maximum

## Deployment

```bash
# Deploy entire service (unified quota policy by default)
kubectl apply -k config/services/resourcemanager.miloapis.com

# Legacy environments: also apply legacy overlay and keep feature gate off
kubectl apply -k config/overlays/legacy-organizations
```
