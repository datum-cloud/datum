# Legacy organizations overlay

Apply this kustomize component when `UnifiedOrganizations=false` on the datum
controller-manager. It installs legacy personal/standard project quota grant
policies and the personal organization name validation policy.

Do **not** apply together with the default unified
`organization-project-quota-policy` — only one grant policy set should be active.

Policy manifests live under this overlay directory so kustomize can build the
component without referencing paths outside the overlay root.

## Usage

Include as a component in your environment kustomization:

```yaml
components:
  - ../path/to/datum/config/overlays/legacy-organizations
```

Ensure milo and datum controller-managers omit `UnifiedOrganizations=true`.
