# Upjet-based Crossplane Provider for MAAS

<div align="center">

![CI](https://github.com/Marouan-chak/provider-upjet-maas/workflows/CI/badge.svg)
[![GitHub release](https://img.shields.io/github/release/Marouan-chak/provider-upjet-maas/all.svg)](https://github.com/Marouan-chak/provider-upjet-maas/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/Marouan-chak/provider-upjet-maas)](https://goreportcard.com/report/github.com/Marouan-chak/provider-upjet-maas)
[![Contributors](https://img.shields.io/github/contributors/Marouan-chak/provider-upjet-maas)](https://github.com/Marouan-chak/provider-upjet-maas/graphs/contributors)

</div>

Provider Upjet-MAAS is a [Crossplane](https://crossplane.io/) provider that is
built using [Upjet](https://github.com/crossplane/upjet) code generation tools
and exposes XRM-conformant managed resources for
[Canonical MAAS](https://maas.io/) (Metal as a Service).

Manage your bare-metal infrastructure as Kubernetes resources.

## Getting Started

### Prerequisites

- Kubernetes cluster with [Crossplane](https://crossplane.io/) installed
- MAAS server (v3.0+)
- MAAS API key

### Install the Provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-maas
spec:
  package: ghcr.io/marouan-chak/provider-upjet-maas:latest
```

### Create Credentials Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: maas-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    { "apiKey": "YOUR_CONSUMER_KEY:YOUR_TOKEN_KEY:YOUR_TOKEN_SECRET" }
```

Get your API key from MAAS UI: **User menu > API keys > Generate MAAS API key**

### Create ProviderConfig

```yaml
apiVersion: maas.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  apiURL: "http://YOUR_MAAS_SERVER:5240/MAAS"
  apiVersion: "2.0"
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: maas-creds
      key: credentials
```

### Create Resources

```yaml
# Create a resource pool
apiVersion: infrastructure.maas.crossplane.io/v1alpha1
kind: ResourcePool
metadata:
  name: my-pool
spec:
  forProvider:
    name: my-pool
    description: "Managed by Crossplane"
  providerConfigRef:
    name: default
---
# Create a fabric
apiVersion: network.maas.crossplane.io/v1alpha1
kind: Fabric
metadata:
  name: my-fabric
spec:
  forProvider:
    name: my-fabric
  providerConfigRef:
    name: default
```

## Supported Resources

| Category | Resources |
|----------|-----------|
| **Network** | Fabric, VLAN, Subnet, SubnetIPRange, Space |
| **Network Interfaces** | InterfacePhysical, InterfaceBond, InterfaceBridge, InterfaceVLAN, InterfaceLink |
| **Machine** | Machine, VMHost, VMHostMachine |
| **Infrastructure** | ResourcePool, Tag, User, Device |
| **DNS** | Domain, Record |
| **Storage** | BlockDevice |

## Documentation

- **[Lab Setup Guide](docs/lab-setup-guide.md)**: Complete guide to setting up a MAAS lab with Raspberry Pi and Proxmox
- **[Examples](examples/)**: Ready-to-use YAML examples for all resources
- **[Contributing](CONTRIBUTING.md)**: How to contribute to this project

## Resource Dependencies

Apply resources in this order:

```
1. ProviderConfig + Credentials
2. ResourcePool, Tag, Space (no dependencies)
3. Fabric
4. VLAN (depends on Fabric)
5. Subnet (depends on Fabric)
6. Machine
7. Interface* resources (depend on Machine)
8. BlockDevice (depends on Machine)
```

## Contributing

For the general contribution guide, see [CONTRIBUTING.md](CONTRIBUTING.md).

If you'd like to learn how to use Upjet, see [Upjet Usage Guide](https://github.com/crossplane/upjet/tree/main/docs).

### Build Locally

```bash
# Initialize submodules (first time only)
git submodule sync && git submodule update --init --recursive

# Run code generation
go run cmd/generator/main.go "$PWD"

# Build provider
make build
```

### Local Development with Kind

```bash
# Setup Kind cluster with Crossplane and provider
./scripts/setup-kind.sh

# Cleanup
./scripts/setup-kind.sh --cleanup
```

## Getting Help

For filing bugs, suggesting improvements, or requesting new resources, please
open an [issue](https://github.com/Marouan-chak/provider-upjet-maas/issues/new/choose).

## External Name Behavior

MAAS uses numeric database IDs as resource identifiers. The `external-name` annotation will contain:

| Resource | External Name Format |
|----------|---------------------|
| Fabric, Space, Subnet | Numeric ID (e.g., `2`) |
| Machine | System ID (e.g., `abc123`) |
| Device | System ID (e.g., `xyz789`) |
| Tag | Tag name (e.g., `my-tag`) |
| VLAN, SubnetIPRange | Numeric ID |

## Troubleshooting

### Check Provider Status

```bash
kubectl get providers
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-upjet-maas
```

### Check Resource Status

```bash
kubectl get managed
kubectl describe <resource-type> <resource-name>
```

### Common Issues

1. **Reference resolution errors**: Ensure dependent resources exist and are Ready
2. **Authentication errors**: Verify API key and MAAS URL in ProviderConfig
3. **Resource stuck in Creating**: Check MAAS UI for the resource state

## License

The provider is released under the [Apache 2.0 license](LICENSE) with [notice](NOTICE).
