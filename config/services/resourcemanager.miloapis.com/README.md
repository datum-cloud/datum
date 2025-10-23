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

- **Personal Organizations**: 2 project maximum
- **Standard Organizations**: 10 projects maximum

## Deployment

```bash
# Deploy entire service
kubectl apply -k config/services/resourcemanager.miloapis.com

# Deploy specific components
kubectl apply -k config/services/resourcemanager.miloapis.com/quota
kubectl apply -k config/services/resourcemanager.miloapis.com/validation
```
