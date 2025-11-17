# Datum is building the internet for AI.

<p align="left">
  <a href="https://cloud.datum.net">Datum Cloud</a> -
  <a href="https://docs.datum.net">Docs</a> -
  <a href="https://link.datum.net/discord">Community Discord</a> -
  <a href="https://www.datum.net/blog/">Blog</a> -
  <a href="https://www.linkedin.com/company/datum-cloud/">LinkedIn</a>
</p>

## 🤝 Overview

Datum was founded to help 1k new clouds thrive in the age of AI. Unlike existing alternatives, Datum’s open network cloud is built specifically for modern service providers, can be shipped anywhere, and is backed by an AGPLv3 open source license. 

With Datum, cloud and AI-native builders can use the tools they love (like Cursor or a Kubernetes native CLI) to access the internet superpowers that today’s tech giants leverage at scale: authoritative DNS, distributed proxies, global backbones, deterministic routing, cloud on-ramps, and private interconnection. 

## Our Purpose

We believe that most people devote their time, energy, families, reputations and money to something not because of what it does, but why it exists and what it believes about the world. When we introduce Datum to prospective users, customers, investors, partners or employees, here is what we share.

- We are connectors — of people, businesses, apps and networks. ABCD!
- We are operators at heart who know how to get stuff done.
- We build for scale with thoughtful abstractions. "Our future selves will thank us."
- We believe that “open” is better, software is the customer, and partners have real value.
- We value grit, humility and hunger.

## So What is Datum?  

### 🌍 Neutral by design

- No allegiance to a single cloud, vendor, or region
- Operates as a trusted, independent layer for alt clouds, incumbents, and
  digital-first enterprises

### ⚙️ Fully programmable

- Designed for developers, operators, and modern service providers
- Built using Kubernetes API patterns for operator happiness and ecosystem
  tooling (`kubectl`, Helm, etc.)

### 🛰 Bring your own infra

- Use Datum’s cloud control plane along with its global network and distributed
  compute
- Or run components in your own cloud or infra (e.g., GCP, AWS, NetActuate,
  Vultr, etc.)

### 🔍 Observability & determinism

- Feed full telemetry to your preferred tools (Grafana Cloud, Honeycomb,
  Datadog, etc.)
- Support for policy enforcement via SRv6

---

## 🚀 Some of our Favorite Features

### Declarative management

Define your desired infrastructure state using Kubernetes manifests. Datum
controllers continuously work to reconcile the actual state with your declared
configuration. This enables infrastructure-as-code practices and GitOps
workflows.

### Kubernetes native

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

## Components

### Datum API server

Built using Kubernetes API server libraries for compatibility with ecosystem
tools (`kubectl`, Helm, etc.), but focused on Datum-specific resources, not
standard Kubernetes workload types (like Pods or Deployments). This approach
takes advantage of recent developments in the Kubernetes project to build a
[generic control plane (KEP-4080)][kep-4080].

### [Network services operator](https://github.com/datum-cloud/network-services-operator)

Manages networking primitives like Datum VPC Networks
(`Network`,`NetworkContext`), Subnets (`SubnetClaim`, `Subnet`), IP Address
Management(IPAM), and network policy concepts (`NetworkBinding`,
`NetworkPolicy`).

### [Workload operator](https://github.com/datum-cloud/workload-operator)

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
Guide](https://www.datum.net/docs/quickstart/) to begin connecting and
managing your infrastructure.

There are also other ways to get involved:

- **Development Setup:** See the [Development
  Guide](https://docs.datum.net/docs/developer-guide/).
- **Roadmap & Enhancements:** Visit our [enhancements
  repo](https://link.datum.net/enhancements).
- **General Discussion:** Join us on the [Datum Community
  Discord](https://discord.com/invite/AeA9XZu4Py).

## License

Datum is primarily licensed under the [AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
Specific components mayhave different licenses; please check individual
repositories for details.
