# MAAS Provider Examples

## Quick Start

### 1. Setup Kind Cluster with Crossplane

```bash
./scripts/setup-kind.sh
```

Or manually:

```bash
# Create cluster
kind create cluster --name crossplane-test

# Install Crossplane
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm install crossplane crossplane-stable/crossplane \
    --namespace crossplane-system \
    --create-namespace

# Build and load provider image
make build
kind load docker-image provider-upjet-maas-amd64 --name crossplane-test

# Apply CRDs
kubectl apply -f package/crds/

# Apply provider
kubectl apply -f examples/provider/provider.yaml
```

### 2. Configure Credentials

Update `examples/providerconfig/creds.yaml` with your MAAS API key:

```yaml
stringData:
  credentials: |
    {
      "apiKey": "YOUR_API_KEY_HERE"
    }
```

Update `examples/providerconfig/providerconfig.yaml` with your MAAS URL:

```yaml
spec:
  apiURL: "http://your-maas-server:5240/MAAS"
```

Apply:

```bash
kubectl apply -f examples/providerconfig/
```

### 3. Create Resources

#### Network Resources (Core)

```bash
# Create a fabric
kubectl apply -f examples/resources/network/fabric.yaml

# Create VLAN (depends on fabric)
kubectl apply -f examples/resources/network/vlan.yaml

# Create subnet (depends on fabric and vlan)
kubectl apply -f examples/resources/network/subnet.yaml

# Create space
kubectl apply -f examples/resources/network/space.yaml

# Create standalone IP range
kubectl apply -f examples/resources/network/subnet-ip-range.yaml
```

#### DNS Resources

```bash
# Create DNS domain
kubectl apply -f examples/resources/dns/domain.yaml

# Create DNS record (depends on domain)
kubectl apply -f examples/resources/dns/record.yaml
```

#### Infrastructure Resources

```bash
# Create resource pool
kubectl apply -f examples/resources/infrastructure/resource-pool.yaml

# Create tag
kubectl apply -f examples/resources/infrastructure/tag.yaml

# Create user (requires password secret)
kubectl apply -f examples/resources/infrastructure/user.yaml

# Create device
kubectl apply -f examples/resources/infrastructure/device.yaml
```

#### Machine Resources

```bash
# Create machine (requires power params secret)
kubectl apply -f examples/resources/machine/machine.yaml

# Create VM host
kubectl apply -f examples/resources/machine/vm-host.yaml

# Create VM on host (depends on vm-host)
kubectl apply -f examples/resources/machine/vm-host-machine.yaml
```

#### Storage Resources

```bash
# Create block device (depends on machine)
kubectl apply -f examples/resources/storage/block-device.yaml
```

#### Network Interface Resources

```bash
# Configure physical interface (depends on machine)
kubectl apply -f examples/resources/network/interface-physical.yaml

# Create bond interface (depends on machine)
kubectl apply -f examples/resources/network/interface-bond.yaml

# Create bridge interface (depends on machine)
kubectl apply -f examples/resources/network/interface-bridge.yaml

# Create VLAN interface (depends on machine, fabric)
kubectl apply -f examples/resources/network/interface-vlan.yaml

# Link interface to subnet (depends on machine, subnet)
kubectl apply -f examples/resources/network/interface-link.yaml
```

### 4. Verify

```bash
# Check all MAAS resources
kubectl get managed

# List resources by type
kubectl get fabrics,vlans,subnets,spaces
kubectl get domains,records
kubectl get resourcepools,tags,users,devices
kubectl get machines,vmhosts,vmhostmachines
kubectl get blockdevices
kubectl get interfacephysicals,interfacebonds,interfacebridges

# Describe a specific resource
kubectl describe fabric example-fabric

# Check provider logs
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-upjet-maas
```

### 5. Cleanup

```bash
# Delete resources (in reverse dependency order)
kubectl delete -f examples/resources/storage/
kubectl delete -f examples/resources/machine/
kubectl delete -f examples/resources/infrastructure/
kubectl delete -f examples/resources/dns/
kubectl delete -f examples/resources/network/

# Delete cluster
kind delete cluster --name crossplane-test
```

## Resource Dependencies

```
                    ResourcePool
                         │
                         ▼
Fabric ──────────► Machine ◄──────── Tag
   │                  │
   ▼                  ▼
 VLAN            BlockDevice
   │
   ▼
Subnet ──────► InterfaceLink
   │                │
   ▼                ▼
SubnetIPRange   InterfacePhysical
                InterfaceBond
                InterfaceBridge
                InterfaceVLAN

VMHost
   │
   ▼
VMHostMachine

Domain
   │
   ▼
Record

Space (standalone)
Device (standalone)
User (standalone)
```

## Available Resources

### Network (network.maas.crossplane.io)
| Kind | Description |
|------|-------------|
| `Fabric` | Network fabric |
| `VLAN` | VLAN on a fabric |
| `Subnet` | Subnet on a VLAN |
| `SubnetIPRange` | IP range within a subnet |
| `Space` | Network space |
| `InterfacePhysical` | Physical network interface |
| `InterfaceBond` | Bonded network interfaces |
| `InterfaceBridge` | Bridge interface |
| `InterfaceVLAN` | VLAN interface |
| `InterfaceLink` | Interface-to-subnet link |

### DNS (dns.maas.crossplane.io)
| Kind | Description |
|------|-------------|
| `Domain` | DNS domain |
| `Record` | DNS record |

### Infrastructure (infrastructure.maas.crossplane.io)
| Kind | Description |
|------|-------------|
| `ResourcePool` | Machine resource pool |
| `Tag` | Tag for machines |
| `User` | MAAS user |
| `Device` | Non-deployable device |

### Machine (machine.maas.crossplane.io)
| Kind | Description |
|------|-------------|
| `Machine` | Physical/virtual machine |
| `VMHost` | VM host (LXD/virsh) |
| `VMHostMachine` | VM on a host |

### Storage (storage.maas.crossplane.io)
| Kind | Description |
|------|-------------|
| `BlockDevice` | Block storage device |

## Troubleshooting

### Provider not healthy

```bash
kubectl describe provider provider-upjet-maas
kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-upjet-maas
```

### Resource stuck in "Creating"

```bash
kubectl describe <resource-type> <resource-name>
# Check Events section for errors
```

### API connection errors

1. Verify MAAS URL is accessible from the cluster
2. Check API key is correct
3. For local MAAS, you may need to expose it to Kind:

```bash
# Get host IP accessible from Kind
docker network inspect kind | grep Gateway
```
