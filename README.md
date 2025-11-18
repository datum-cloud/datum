# Datum is building the internet for AI.

<p align="left">
  <a href="https://cloud.datum.net">Datum Cloud</a> -
  <a href="https://docs.datum.net">Docs</a> -
  <a href="https://link.datum.net/discord">Community Discord</a> -
  <a href="https://www.datum.net/blog/">Blog</a> -
  <a href="https://www.linkedin.com/company/datum-cloud/">LinkedIn</a>
</p>

## Why Datum?

We believe that AI is changing everything — not just how we work and create, but how
quickly new applications, agents, and even clouds are being built. Digital
leaders today must orchestrate a complex, fragmented web of clouds, specialty
providers, customers, and data.

The Internet is built on data center [meet-me
rooms](https://en.wikipedia.org/wiki/Meet-me_room), where telco providers and
hyperscaler clouds talk to each other over real physical cables, called
cross-connects. A new connection takes days or weeks of humans moving things
around to set up.

We believe the next era of the Internet is already here, and it's growing
fast. In the [alt-cloud](https://github.com/datum-cloud/awesome-alt-clouds)
universe, you don't think about virtual machines and VPCs, you think about
*services*. You connect your Vercel app with your Supabase instance all wired up
with your Kestra workflow, monitored by your Resolve SRE agent. There isn't a
switch or routing table in sight. It's just virtual plumbing to make your
business go. 

With Datum, cloud and AI-native builders can use the tools they love (like
Cursor or a Kubernetes native CLI) to access the internet superpowers that
today’s tech giants leverage at scale: authoritative DNS, distributed proxies,
global backbones, deterministic routing, cloud on-ramps, and private
interconnection.

That's why we're building Datum: to help build 1k clouds in the age of AI.

## So What is Datum?

### Fully programmable and AI-native

- Developer- and agent-friendly protocols, interfaces, and workflows
- Programmatic interconnection between providers and services
- Security through network-level encryption that's built-in and impossible to break or disable
- Built using the "operating system for AI" Kubernetes API patterns for operator
  tooling and familiarity (`kubectl`, Helm, etc.)

### Neutral by design

- No allegiance to a single cloud, vendor, or region
- Operates as a trusted, independent layer for alt clouds, incumbents, and
  digital-first enterprises

### Bring your own infra

- Use Datum’s cloud control plane along with its global network and distributed
  compute
- Or run components in your own cloud or infra (e.g., GCP, AWS, NetActuate,
  Vultr, etc.)

### Maximum flexibility

- Feed full telemetry to your preferred tools (Grafana Cloud, Honeycomb,
  Datadog, etc.)
- Support for policy enforcement via SRv6

---

## Some of our Favorite Features

### Declarative management

Datum works just like Kubernetes, because it *is* Kubernetes. Define your desired infrastructure state and our components reconcile the living system to match. No more syncing or drift.

The Datum control plane is natively compatible with tooling from the Kubernetes
ecosystem. Datum APIs are defined as [Custom Resources][k8s-custom-resources],
and resources are managed by operators which can be deployed into any Kubernetes
cluster.

Use the tools you're familiar with - `kubectl` for interacting with API
resources via the CLI, `kustomize` or `terraform` for configuration management,
or any other tool compatible with the Kubernetes API.

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
Datum infrastructure.

## Components

### Datum API server

We deploy a Datum variant of the Kubernetes api-server in the style of the [generic control plane (KEP-4080)][kep-4080] so that we can leverage the vast ecosystem of libraries and tooling. There is no need to design a bespoke, infrastructure-focused distributed system for you to learn; Kubernetes has the primitives to support it.  While the standard api-server operates normally for the cluster itself (think Pods and Deployments), Datum's api-server handles Datum-specific resources like `Network` and `Workload`.

[kep-4080]: https://github.com/kubernetes/enhancements/tree/master/keps/sig-api-machinery/4080-generic-controlplane

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

## Get Started

The easiest way to understand Datum is to try it! Head over to [Datum
Cloud](https://cloud.datum.net), sign up, and follow the [Quickstart
Guide](https://www.datum.net/docs/quickstart/) to begin your journey to a reimagined world of interconnection.

We hope that you will then come and build with us:

- **General Discussion:** Join us on the [Datum Community
  Discord](https://link.datum.net/discord).
- **Development Setup:** See the [Development
  Guide](https://docs.datum.net/docs/developer-guide/).
- **Roadmap & Enhancements:** Visit our [enhancements
  repo](https://link.datum.net/enhancements).

## License

Datum is primarily licensed under the [AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
Specific components mayhave different licenses; please check individual
repositories for details.







