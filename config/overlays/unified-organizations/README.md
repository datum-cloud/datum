# Unified organizations overlay

The default datum service configuration already applies
`organization-project-quota-policy` (10 projects per organization).

Apply this component together with `UnifiedOrganizations=true` on milo and datum
controller-managers. No additional quota resources are required beyond the
default service kustomization.

For legacy environments, use `config/overlays/legacy-organizations` instead and
keep `UnifiedOrganizations=false`.
