# Unified organizations overlay

Apply this kustomize component when `UnifiedOrganizations=true` on milo and datum
controller-managers. It replaces legacy personal/standard project quota grant
policies with a single 10-project policy for all organizations and removes the
personal organization display-name validation policy.

The default datum service configuration keeps legacy behavior so releases do not
require coordinated cutover across environments. Only environments that opt in
(staging first) should apply this overlay together with the feature gate.

## Usage

Include as a component in your environment kustomization:

```yaml
components:
  - ../path/to/datum/config/overlays/unified-organizations
```

Ensure milo and datum controller-managers run with `UnifiedOrganizations=true`.
