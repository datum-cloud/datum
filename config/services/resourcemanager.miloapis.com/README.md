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

The default configuration ships legacy personal/standard quota grant policies:

- **Personal Organizations**: 2 project maximum
- **Standard Organizations**: 10 projects maximum

When enabling unified organizations, apply `config/overlays/unified-organizations`
with `UnifiedOrganizations=true` on milo and datum controller-managers for a
single 10-project quota policy on all organizations.

## Deployment

```bash
# Deploy entire service (legacy quota policies by default)
kubectl apply -k config/services/resourcemanager.miloapis.com

# Unified environments: also apply unified overlay and enable feature gate
kubectl apply -k config/overlays/unified-organizations
```
