# `stripe.billing.miloapis.com` — Datum platform overlay

Datum-specific wiring for the upstream `milo-os/stripe-provider`. The
upstream `milo-integration` Component (in
`milo-os/stripe-provider/config/components/milo-integration/`)
deliberately leaves all volume sources unspecified so each operator can
plug in their own trust bundle, cert-manager Issuer, and webhook cert
issuance mechanism. This directory is Datum's choice.

## Layout

- `datum-control-plane-wiring/` — Kustomize Component that adds the four
  `volumes:` entries the upstream `volumeMounts` reference:
  - `trust-bundle` — ConfigMap `datum-control-plane-trust-bundle`,
    populated by trust-manager when the namespace carries the
    `infra.datum.net/inject-datum-control-plane-trust-bundle: ""` label.
  - `webhook-tls` — cert-manager CSI volume issued by the
    `datum-control-plane` ClusterIssuer.
  - `client-cert` — cert-manager CSI volume issuing
    `system:control@stripe.billing.miloapis.com` so Milo recognises the
    controller as the stripe-provider service identity.
  - `discovery-kubeconfig` — ConfigMap
    `stripe-provider-discovery-kubeconfig`, supplied by the deployment
    overlay alongside this Component.

## Consumption

This Component lives here so Datum has one canonical place to describe
how the stripe-provider controller plugs into the Datum control plane.
It is **not** wired into `config/services/kustomization.yaml` —
`services/kustomization.yaml` is applied milo-side by the
`datum-milo-customization` Flux Kustomization, and this Component is a
host-cluster Deployment patch, not a milo-side resource.

The active consumer today is the deployment Kustomization in
[datum-cloud/infra](https://github.com/datum-cloud/infra) at
`apps/stripe-provider/base/manager.yaml`. Infra carries an inline copy
of the same patch because Flux Kustomization can only resolve one
`sourceRef` and we cannot compose this Component alongside the upstream
`milo-os/stripe-provider` OCI bundle from the same Flux apply. Until
that composition is solved, keep both copies in sync — this file is the
source of truth.
