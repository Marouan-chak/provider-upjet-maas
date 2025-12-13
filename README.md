# Provider MAAS

`provider-upjet-maas` is a [Crossplane](https://crossplane.io/) provider for
[Canonical MAAS](https://maas.io/) (Metal as a Service) built using [Upjet](https://github.com/crossplane/upjet).

Manage your bare-metal infrastructure as Kubernetes resources.

## Supported Resources

| Category | Resources |
|----------|-----------|
| **Network** | Fabric, VLAN, Subnet, SubnetIPRange, Space |
| **Network Interfaces** | InterfacePhysical, InterfaceBond, InterfaceBridge, InterfaceVLAN, InterfaceLink |
| **Machine** | Machine, VMHost, VMHostMachine |
| **Infrastructure** | ResourcePool, Tag, User, Device |
| **DNS** | Domain, Record |
| **Storage** | BlockDevice |

## Quick Start

### 1. Install the Provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-maas
spec:
  package: docker.io/marouandock/provider-maas:latest
```

### 2. Create Credentials Secret

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

Get your API key from MAAS UI: **User menu → API keys → Generate MAAS API key**

### 3. Create ProviderConfig

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

### 4. Create Resources

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

## Documentation

- **[Lab Setup Guide](docs/lab-setup-guide.md)**: Complete guide to setting up a MAAS lab with Raspberry Pi and Proxmox
- **[Examples](examples/)**: Ready-to-use YAML examples for all resources

## Examples Directory Structure

```
examples/
├── providerconfig/
│   ├── creds.yaml           # Credentials secret
│   └── providerconfig.yaml  # ProviderConfig
├── resources/
│   ├── infrastructure/
│   │   ├── resource-pool.yaml
│   │   └── tag.yaml
│   ├── network/
│   │   ├── fabric.yaml
│   │   ├── vlan.yaml
│   │   ├── subnet.yaml
│   │   ├── subnet-ip-range.yaml
│   │   ├── space.yaml
│   │   ├── interface-physical.yaml
│   │   ├── interface-bond.yaml
│   │   ├── interface-bridge.yaml
│   │   ├── interface-vlan.yaml
│   │   └── interface-link.yaml
│   ├── machine/
│   │   └── machine.yaml
│   ├── dns/
│   │   ├── domain.yaml
│   │   └── record.yaml
│   └── storage/
│       └── block-device.yaml
```

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

## Developing

### Prerequisites

- Go 1.24+
- Docker
- Crossplane CLI

### Build

```bash
# Initialize submodules (first time only)
git submodule sync && git submodule update --init --recursive

# Run code generation
go run cmd/generator/main.go "$PWD"

# Build provider
make build
```

### Push Package

```bash
# Push to registry
crossplane xpkg push \
  --package-files="_output/xpkg/linux_amd64/provider-upjet-maas-*.xpkg" \
  docker.io/YOUR_USERNAME/provider-maas:latest
```

### Local Development

```bash
# Setup Kind cluster with Crossplane and provider
./scripts/setup-kind.sh

# Cleanup
./scripts/setup-kind.sh --cleanup
```

## External Name Behavior

MAAS uses numeric database IDs as resource identifiers. The `external-name` annotation will contain:

| Resource | External Name Format |
|----------|---------------------|
| Fabric, Space, Subnet | Numeric ID (e.g., `2`) |
| Machine | System ID (e.g., `abc123`) |
| Device | System ID (e.g., `xyz789`) |
| Tag | Tag name (e.g., `my-tag`) |
| VLAN, SubnetIPRange | Numeric ID |

This is expected behavior - MAAS internally identifies resources by database primary keys.

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

## Contributing

Issues and PRs welcome at [GitHub](https://github.com/Marouan-chak/provider-upjet-maas).

## License

Apache 2.0
