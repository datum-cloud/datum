# Datum is building the internet for AI.

<p align="left">
  <a href="https://cloud.datum.net">Datum Cloud</a> -
  <a href="https://www.datum.net/docs">Docs</a> -
  <a href="https://link.datum.net/discord">Community Discord</a> -
  <a href="https://www.datum.net/blog/">Blog</a> -
  <a href="https://www.linkedin.com/company/datum-cloud/">LinkedIn</a>
</p>

## Why Datum?

We believe that AI is changing everything — not just how we work and create, but how
quickly new applications, agents, and clouds are being built. 

The Internet is built on data center [meet-me
rooms](https://en.wikipedia.org/wiki/Meet-me_room), where telco providers and
hyperscaler clouds talk to each other over real physical cables, called
cross-connects. A new connection takes days or weeks of humans moving things
around to set up.

We believe the next era of the Internet is already here, and it's growing
fast. In the [alt-cloud](https://www.alt-cloud.org/)
universe, you don't think about virtual machines and VPCs, you think about
*services*. You connect your Vercel app with your Supabase instance all wired up
with your Kestra workflow, monitored by your favorite SRE agent. There isn't a
switch or routing table in sight. It's just virtual plumbing to make your
business go, powered by a fleet of agents. 

With Datum, cloud and AI-native builders can use the tools they love (like
Claude, Cursor or a Kubernetes native CLI) to access the internet superpowers that
today’s tech giants leverage at scale: authoritative DNS, edge proxies,
global backbones, deterministic routing, cloud on-ramps, and private
interconnection.

That's why we're building Datum: to help build 1k clouds in the age of AI.

## So what is Datum?

### An open network cloud built for AI

- Developer and agent-friendly protocols, interfaces, and workflows
- Backed by an AGPvL 3.0 license
- Powerful suite of infrastructure primitives, deployed at the edge
- Built using Kubernetes API patterns for operator tooling and familiarity (`datumctl`, Helm, etc.)
- Flexible deployment models (public cloud, managed cloud, BYOC, OSS)

### Neutral & flexible by design

- Ecosystem friendly partner model
- No allegiance to a single cloud, vendor, or region
- Operates as a trusted, independent layer for alt clouds, incumbents, and digital-first enterprises
- Feed full telemetry to your preferred tools (Grafana Cloud, etc)

---

## Key features

### Declarative management

Our most important feature isn't a "what" but a "how". Datum works just like Kubernetes, because it *is* Kubernetes. Define your desired infrastructure state and our components reconcile the living system to match. No more syncing or drift.

The Datum control plane is natively compatible with tooling from the Kubernetes
ecosystem. Datum APIs are defined as [Custom Resources][k8s-custom-resources],
and resources are managed by operators which can be deployed into any Kubernetes
cluster.

Use the tools you're familiar with, but especially `datumctl` for interacting with API
resources via the CLI. Read more [about datumctl here](https://www.datum.net/docs/datumctl/overview).

### AI Edge
An Envoy-based edge that provides an intelligent HTTPProxy along with a Coraza-based Web Application Firewall (WAF) to protect and route internet traffic to backend services. We support HTTP(S) 1.1, HTTP2, gRPC, and WebSockets.

### Galactic VPCs

Internet backbones weren't designed for most humans, let alone agents. Our Galactic VPC feature is built for an agentic world to provide policy-driven SRv6 virtual backbones that go anywhere.

### UFOs

We've partnered with Unikraft to build out a modern edge compute layer that is ideal for agentic and network use cases. "Unikernel Function Offloads" provide 100% isolation, millisecond cold starts, and scale to zero snapshotting.

### Connectors

We plan to support all kinds of connections, from developer-focused (e.g. Tailscale Tailnets, Wireguard VPNs) to low level L2/L3 telco (AWS Direct Connect, Equinix Fabric, Megaport Onramps, etc). We've started with QUIC-based tunnels powered by the [Iroh protocol](https://www.iroh.computer/).

### Essentials

We support a growing collection of features that help make agentic and internet scale applications "go". While these may not be the star of any show, they are necessary ingredients.  

- Authoritative DNS
- Domain resource tracking
- Fine grained roles and permissions
- Secrets & machine accounts
- Activity logs

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

---

## Get Started

The easiest way to understand Datum is to try it! Head over to [Datum
Cloud](https://cloud.datum.net), sign up, and [check out our docs](https://www.datum.net/docs/platform/setup) to get started.

We hope that you will then come and build with us:

- **General Discussion:** Join us on the [Datum Community
  Discord](https://link.datum.net/discord).
- **Enhancements:** Visit our [enhancements
  repo](https://link.datum.net/enhancements).
- **Milestones:** Visit our [planned milestones](https://link.datum.net/enhancements).
  

## License

Datum is primarily licensed under the [AGPL v3.0](https://www.gnu.org/licenses/agpl-3.0.en.html).
Specific components mayhave different licenses; please check individual
repositories for details.







