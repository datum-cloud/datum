# Milo API Groups Configuration

This directory contains Milo configurations for **all Milo API groups** in Datum Cloud.

It defines common validation admission policies and other base configurations needed for correct operation of all Milo API resources.

## Overview

- Applies to all Milo API groups
- Admission policies for controlling access
- (e.g., policies to prevent unapproved users from accessing the Datum platform)

## Validation and Security

A key admission policy included in this configuration is **approved user enforcement**. This policy ensures that only users who have been approved may interact with the Datum platform via Milo APIs.

- Requests from not-approved users are automatically denied at the admission level.

## Structure

```
├── validation/      # Admission policies for all Milo APIs
└── kustomization.yaml
```

## Deployment

```bash
# Deploy entire Milo API group configuration
kubectl apply -k config/services/miloapis.com

# Deploy just the validation policies
kubectl apply -k config/services/miloapis.com/validation
```
