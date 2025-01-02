# Datum

Datum is a network cloud you can take anywhere, backed by open source.

The Datum control plane is a collection of multiple projects developed with
Kubernetes control plane technology, most of which can be installed into native
Kubernetes clusters.

## Get started

- Quick Start with GCP - Coming Soon!
- Development Guide - Coming Soon!

## Documentation

Our documentation is available at [docs.datum.net](https://docs.datum.net/). If you want to learn more about what's under development or suggest features, please visit our [feedback site](https://feedback.datum.net).

## Key Components

### API Server

The Datum API server leverages Kubernetes API server libraries to enable
compatibility with existing Kubernetes ecosystems tooling such as kubectl, helm,
kustomize, Terraform, Pulumi, Ansible, kubebuilder, operator-sdk, and more.

While the Datum API server exposes a handful of existing Kubernetes API types
such as Secrets and ConfigMaps, you will not find definitions for Pods,
Deployments, Services, etc. This approach takes advantage of recent developments
in the Kubernetes project to build a [generic control plane][kep-4080], exposing
libraries that external software can depend on and build upon.

[kep-4080]: https://github.com/kubernetes/enhancements/tree/master/keps/sig-api-machinery/4080-generic-controlplane

### Kubernetes Operators

Datum leverages the [operator pattern][operator-pattern] to define APIs and implement
controllers via the use of [kubebuilder][kubebuilder]. Each Datum operator can
be deployed into any native Kubernetes cluster that meets minimum API version
requirements, and does not rely on specific functionality provided by the Datum
API server.

[operator-pattern]: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/
[kubebuilder]: https://github.com/kubernetes-sigs/kubebuilder

#### [Network Services Operator](https://github.com/datum-cloud/network-services-operator)

APIs:

- Networks, NetworkContexts, NetworkBindings, and NetworkPolicies
- SubnetClaims and Subnets
- IP Address Management (IPAM)

Controller responsibilities:

- Creating NetworkContexts as required by NetworkBindings.
- Allocating Subnets to SubnetClaims

#### [Workload Operator](https://github.com/datum-cloud/workload-operator)

APIs:

- Workloads, WorkloadDeployments, and Instances

Controller responsibilities:

- Creating one or more WorkloadDeployments for candidate locations based on
  Workload placement intent.
- Scheduling WorkloadDeployments onto Locations.

The [Workloads RFC][workload-rfc] is recommended reading for those interested in
the design goals of this system.

[workload-rfc]: https://github.com/datum-cloud/workload-operator/blob/integration/datum-poc/docs/compute/development/rfcs/workloads/README.md

#### [GCP Infrastructure Provider](https://github.com/datum-cloud/infra-provider-gcp)

Integrates with the Network Services and Workload Operators to provision
resources within Google Cloud Platform (GCP). This operator connects to:

- An upstream control plane hosting Datum entities
- An infrastructure control plane running [GCP Config
  Connector](https://github.com/GoogleCloudPlatform/k8s-config-connector)

Controller responsibilities:

- Maintaining instances as Virtual Machines in GCP.
- Maintaining VPC related entities such as networks and subnets.
- Discovering instances provisioned by GCP controllers, such as managed instance
  groups.

## Get involved

If you choose to contribute to any of our projects, we would love to work with you to ensure a great experience.

- Check out our [roadmap and changelog](https://feedback.datum.net).
- Read and subscribe to the [Datum blog](https://www.datum.net/blog/).
- For general discussions, join us on the [Datum Community Slack](https://slack.datum.net) team.
- Follow [us on LinkedIn](https://www.linkedin.com/company/datum-cloud/).
