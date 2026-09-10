# CLAUDE.md

Guidance for working in this repository (datum-cloud/datum). This is a
Kubebuilder-based Go controller-manager for the Datum control plane.

## Common commands

- `make test` — generate manifests/code, fmt, vet, then run tests (uses envtest).
- `make lint` / `make lint-fix` — run golangci-lint.
- `make build` — build the manager binary.
- `make run` — run a controller against the cluster in `~/.kube/config`.
- `make manifests generate` — regenerate CRDs/RBAC/webhooks and deepcopy code
  after changing API types; commit the regenerated output.

## Deployment flow (read before assuming a change is "live")

CI (`.github/workflows/build-and-test.yaml`) publishes the
`ghcr.io/datum-cloud/datum-kustomize` bundle and the `datum` image on every push
(main builds tagged `v0.0.0-main-*`) and on each release tag (`vX.Y.Z`).
Deployment to clusters is driven by Flux in the `datum-cloud/infra` repo:

- **Merging a PR to `main` ships to STAGING only.** The staging OCIRepository
  tracks `-main-*` prerelease builds and auto-rolls every main build. Merging is
  safe — it does not touch prod.
- **Cutting a release tag (`vX.Y.Z`) auto-ships to PROD.** The prod image policy
  matches release semver tags (`>= 0.0.0`); Flux's image-updater then commits the
  bump to infra `main`, which the prod cluster tracks and reconciles in ~1–2 min.
  There is no further human gate after the tag. **Treat creating a release tag as
  the production deploy action**, not merging.

To check the live prod version, read the `datum-kustomize-bundles` OCIRepository
`tag:` in the `datum-system` namespace (or the `$imagepolicy` marker line in
infra's `apps/datum-control-plane-system/core-control-plane/production/datum-oci-repository-patch.yaml`).

## Conventions

- After changing kubebuilder markers (RBAC, webhooks) or generated types, run
  `make manifests generate` and commit the regenerated artifacts. Controllers
  live in `internal/controller`; API types are largely external (see `PROJECT`).
