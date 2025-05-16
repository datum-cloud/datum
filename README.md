# Datum: Connectivity infrastructure to power your unique advantage

<p align="center">
  <a href="https://cloud.datum.net">Cloud Platform</a> -
  <a href="https://docs.datum.net">Docs</a> -
  <a href="https://slack.datum.net">Community Slack</a> -
  <a href="https://www.datum.net/blog/">Blog</a> -
  <a href="https://www.linkedin.com/company/datum-cloud/">LinkedIn</a>
</p>

## 🤝 Datum is an Open Network Cloud

We believe that AI is changing everything — not just how we compute, but how
ecosystems form and interact. Digital leaders today must orchestrate a complex,
fragmented web of clouds, specialty providers, customers, and data. That's why
we're building Datum...to act as a **“meet-me room” for the internet’s next
era** — a neutral, programmable middle layer where companies can
programmatically connect without needing to build and operate the entire stack
themselves.

### 🧠 Built for an AI-Native World

- Developer and agent-friendly interfaces and workflows
- Enables autonomous and programmatic interconnection between providers and
  services

### 🌍 Neutral by Design

- No allegiance to a single cloud, vendor, or region
- Operates as a trusted, independent layer for alt clouds, incumbents, and
  digital-first enterprises

### ⚙️ Fully Programmable

- Designed for developers, operators, and modern service providers
- Built using Kubernetes API patterns for operator happiness and ecosystem
  tooling (`kubectl`, Helm, etc.)

### 🛰 Bring Your Own Infra

- Use Datum’s cloud control plane along with its global network and distributed
  compute
- Or run components in your own cloud or infra (e.g., GCP, AWS, NetActuate,
  Vultr, etc.)

### 🔍 Observability & Determinism

- Feed full telemetry to your preferred tools (Grafana Cloud, Honeycomb,
  Datadog, etc.)
- Support for policy enforcement via SRv6

---

## 🚀 Some of our Favorite Features

### Declarative Management

Define your desired infrastructure state using Kubernetes manifests. Datum
controllers continuously work to reconcile the actual state with your declared
configuration. This enables infrastructure-as-code practices and GitOps
workflows.

### Kubernetes Native

The Datum control plane is natively compatible with tooling from the Kubernetes
ecosystem. Datum APIs are defined as [Custom Resources][k8s-custom-resources],
and resources are managed by operators which can be deployed into any Kubernetes
cluster.

Use the tools you're familiar with - `kubectl` for interacting with API
resources via the CLI, `kustomize` or `terraform` for configuration management
via GitOps practices, or any other tool compatible with the Kubernetes API.

Expect the same behaviors from the Datum control plane as you would from
Kubernetes. Resources are reconciled to ensure intended state has been met,
failures are automatically addressed, and transparent status information is made
available.

[k8s-custom-resources]:
    https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/

### Workloads

The `Workload` resource provides a provider-agnostic way to manage groups of
compute instances (VMs or containers). Define instance templates, placement
rules (where instances should run across locations/providers), scaling behavior,
network attachments, and volume mounts. The responsible infrastructure provider
operator handles the provisioning.

### Gateways

Leveraging the standard Kubernetes Gateway API (`GatewayClass`,
`Gateway`,`HTTPRoute`, etc.), Datum allows you to define how external or
internal traffic should connect to your services. Manage TLS certificates,
configure routing logic, and control network ingress/egress across the
infrastructure managed by Datum.

### Pluggable Infrastructure Providers

Datum uses a provider model to interact with different underlying infrastructure
environments (e.g., GCP, AWS, bare metal). Specific provider operators
(like`infra-provider-gcp`) translate the abstract Datum API resources
(`Workload`,`Gateway`) into concrete resources managed by the target provider.
This allows for consistent management across heterogeneous environments.

## Components

### Datum API Server

Built using Kubernetes API server libraries for compatibility with ecosystem
tools (`kubectl`, Helm, etc.), but focused on Datum-specific resources, not
standard Kubernetes workload types (like Pods or Deployments). This approach
takes advantage of recent developments in the Kubernetes project to build a
[generic control plane (KEP-4080)][kep-4080].

### [Network Services Operator](https://github.com/datum-cloud/network-services-operator)

Manages networking primitives like Datum VPC Networks
(`Network`,`NetworkContext`), Subnets (`SubnetClaim`, `Subnet`), IP Address
Management(IPAM), and network policy concepts (`NetworkBinding`,
`NetworkPolicy`).

### [Workload Operator](https://github.com/datum-cloud/workload-operator)

Manages the lifecycle of `Workload` resources, handling placement logic and the
creation of compute instances (`WorkloadDeployment`, `Instance`) via
infrastructure providers. See the [Workloads
RFC](https://github.com/datum-cloud/enhancements/tree/main/enhancements/compute/workloads)
for design details.

### Plugins

Datum Plugins interpret resource definitions such as Workloads and Networks to
drive the management of provider specific resources such as Virtual Machines and
VPC Networks to meet the declared expectations. Our first example is for [Google
Cloud Platform (GCP)](https://github.com/datum-cloud/infra-provider-gcp).
Supported features include:

- Deploying Virtual Machine based workload instances with OS images provided via
  an image library.
- Deploying sandboxed container based workload instances with any OCI compliant
  container image.
- VPC connectivity and IPAM.
- Attaching instances to one or more networks.

---

## 🔗 Get Started

The easiest way to leverage our value is with [Datum
Cloud](https://cloud.datum.net). Sign up and follow the [Getting Started
Guide](https://docs.datum.net/docs/get-started/) to begin connecting and
managing your infrastructure.

There are also other ways to get involved:

- **Development Setup:** See the [Development
  Guide](https://docs.datum.net/docs/tasks/developer-guide/).
- **Roadmap & Enhancements:** Visit our [enhancements
  repo](https://link.datum.net/enhancements).
- **General Discussion:** Join us on the [Datum Community
  Slack](https://link.datum.net/datumslack).

## License

Datum is primarily licensed under the [AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
Specific components mayhave different licenses; please check individual
repositories for details.
