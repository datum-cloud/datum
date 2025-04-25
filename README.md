# Datum: Connectivity infrastructure to power your unique advantage

<p align="center">
  <a href="https://cloud.datum.net">Cloud Platform</a> -
  <a href="https://docs.datum.net">Docs</a> -
  <a href="https://slack.datum.net">Community Slack</a> -
  <a href="https://www.datum.net/blog/">Blog</a>
</p>

- [Simplify Distributed Infrastructure Management](#simplify-distributed-infrastructure-management)
- [Future Directions: Enhanced Connectivity](#future-directions-enhanced-connectivity)
- [Kubernetes Native](#kubernetes-native)
- [Getting Started](#getting-started)
  - [Datum Cloud (Recommended)](#datum-cloud-recommended)
  - [Self-hosting Datum Operators (Advanced)](#self-hosting-datum-operators-advanced)
- [Core Concepts](#core-concepts)
  - [Declarative Management](#declarative-management)
  - [Pluggable Infrastructure Providers](#pluggable-infrastructure-providers)
  - [Workloads](#workloads)
  - [Gateways](#gateways)
- [Key Features \& Components](#key-features--components)
- [Contributing](#contributing)
- [License](#license)

## Simplify Distributed Infrastructure Management

Datum is an open platform designed to **unify and simplify** how you connect and
manage distributed infrastructure, wherever it runs. Whether you're using public
clouds, private data centers, or edge locations, Datum aims to provide a
**consistent, declarative control plane** built on Kubernetes principles.

Instead of managing each environment separately, Datum allows you to define your
desired infrastructure state using familiar Kubernetes APIs. This declarative
approach means you focus on *what* you want, and Datum works with its pluggable
providers to make it happen.

## Future Directions: Enhanced Connectivity

Datum aims to further simplify global infrastructure management by building out
a robust edge network and expanding connectivity options. Future goals include:

- **Global Edge Network**: Providing optimized ingress and egress points closer
  to users worldwide using Anycast IPs, advanced load balancing, and traffic
  management capabilities.
- **Seamless Interconnect**: Enabling flexible and secure connections between
  Datum VPC Networks and external environments (public clouds, private data
  centers) using various standard and modern tunneling technologies.
- **Integrated Network Services**: Offering built-in services like global load
  balancing and DNS integrated across the entire managed infrastructure fabric.

## Kubernetes Native

Datum embraces the Kubernetes ecosystem. By using Custom Resource Definitions
(CRDs) and standard APIs:

- **Declarative Configuration**: Define your entire infrastructure topology in
  code.
- **Ecosystem Compatibility**: Integrate seamlessly with existing Kubernetes
  tools like `kubectl`, Kustomize, Helm, GitOps controllers (Argo CD, Flux), and
  more.
- **Extensibility**: Build on standard Kubernetes constructs and APIs.

Get started quickly with [Datum Cloud](https://cloud.datum.net), or run the
Datum operators directly within your own Kubernetes clusters for maximum
control.

## Getting Started

### Datum Cloud (Recommended)

The easiest way to leverage Datum is via the hosted [Datum
Cloud](https://cloud.datum.net) platform. Sign up and follow the [Getting
Started Guide](https://docs.datum.net/docs/get-started/) to begin connecting and
managing your infrastructure.

### Self-hosting Datum Operators (Advanced)

For users who prefer to manage their own control plane, Datum's core components
are available as Kubernetes operators. You can install these operators
(`workload-operator`, `network-services-operator`, relevant `infra-provider-*`
operators) into any existing Kubernetes cluster. A set of Kustomizations are
maintained for each component, which Datum builds upon internally for operating
the platform.

- **Quick Start (GCP Example):** [Set up a Datum managed Location backed by
  GCP](https://docs.datum.net/docs/tutorials/infra-provider-gcp/)
- **Development Guide:** [General Development
  Setup](https://docs.datum.net/docs/tasks/developer-guide/)

## Core Concepts

Datum extends the Kubernetes API to provide a unified infrastructure control
plane.

### Declarative Management

Define your desired infrastructure state using Kubernetes manifests. Datum
controllers continuously work to reconcile the actual state with your declared
configuration. This enables infrastructure-as-code practices and GitOps
workflows.

### Pluggable Infrastructure Providers

Datum uses a provider model to interact with different underlying infrastructure
environments (e.g., GCP, AWS, bare metal). Specific provider operators (like
`infra-provider-gcp`) translate the abstract Datum API resources (`Workload`,
`Gateway`) into concrete resources managed by the target provider. This allows
for consistent management across heterogeneous environments.

### Workloads

The `Workload` resource provides a provider-agnostic way to manage groups of
compute instances (VMs or containers). Define instance templates, placement
rules (where instances should run across locations/providers), scaling behavior,
network attachments, and volume mounts. The responsible infrastructure provider
operator handles the provisioning.

### Gateways

Leveraging the standard Kubernetes Gateway API (`GatewayClass`, `Gateway`,
`HTTPRoute`, etc.), Datum allows you to define how external or internal traffic
should connect to your services. Manage TLS certificates, configure routing
logic, and control network ingress/egress across the infrastructure managed by
Datum.

## Key Features & Components

Datum consists of several core components, primarily implemented as Kubernetes
Operators:

- **API Server**: Built using Kubernetes API server libraries for compatibility
  with ecosystem tools (`kubectl`, Helm, etc.), but focused on Datum-specific
  resources, not standard Kubernetes workload types (like Pods or Deployments).
  This approach takes advantage of recent developments in the Kubernetes project
  to build a [generic control plane
  (KEP-4080)](https://github.com/kubernetes/enhancements/tree/master/keps/sig-api-machinery/4080-generic-controlplane).
- **[Network Services
  Operator](https://github.com/datum-cloud/network-services-operator)**: Manages
  networking primitives like Datum VPC Networks (`Network`, `NetworkContext`),
  Subnets (`SubnetClaim`, `Subnet`), IP Address Management (IPAM), and network
  policy concepts (`NetworkBinding`, `NetworkPolicy`).
- **[Workload Operator](https://github.com/datum-cloud/workload-operator)**:
  Manages the lifecycle of `Workload` resources, handling placement logic and
  the creation of compute instances (`WorkloadDeployment`, `Instance`) via
  infrastructure providers. See the [Workloads RFC](https://github.com/datum-cloud/enhancements/tree/main/enhancements/compute/workloads)
  for design details.
- **Infrastructure Providers** (e.g.,
  [`infra-provider-gcp`](https://github.com/datum-cloud/infra-provider-gcp)):
  Pluggable operators that translate abstract Datum resource definitions
  (`Workload`, network configurations) into concrete actions and resources
  within a specific target environment (like GCP, AWS, etc.). They often
  integrate with provider-specific tooling (e.g., GCP Config Connector).

Leverage standard Kubernetes APIs and patterns to manage these components and
define your infrastructure.

## Contributing

We welcome contributions!

Get involved:

- **Development Setup:** See the [Development Guide](https://docs.datum.net/docs/tasks/developer-guide/).
- **Roadmap & Enhancements:** Visit our [enhancements repo](https://github.com/orgs/datum-cloud/projects/22).
- **General Discussion:** Join us on the [Datum Community Slack](https://slack.datum.net).
- Follow [us on LinkedIn](https://www.linkedin.com/company/datum-cloud/).

## License

Datum is primarily licensed under the [AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
Specific components mayhave different licenses; please check individual
repositories for details.
