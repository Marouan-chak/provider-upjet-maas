# MAAS Lab Setup Guide

Complete guide to setting up a MAAS lab with Raspberry Pi as the controller and Proxmox VMs as target machines, managed via Crossplane.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Prerequisites](#prerequisites)
3. [Part 1: Hardware Setup](#part-1-hardware-setup)
4. [Part 2: MAAS Installation on Raspberry Pi](#part-2-maas-installation-on-raspberry-pi)
5. [Part 3: Proxmox Configuration](#part-3-proxmox-configuration)
6. [Part 4: MAAS Network Configuration](#part-4-maas-network-configuration)
7. [Part 5: Crossplane Provider Installation](#part-5-crossplane-provider-installation)
8. [Part 6: Applying MAAS Resources](#part-6-applying-maas-resources)
9. [Troubleshooting](#troubleshooting)

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                           Home Network                               │
│                         192.168.0.0/24                               │
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────────────┐ │
│  │   Router     │     │ Raspberry Pi │     │    Proxmox Host      │ │
│  │ 192.168.0.1  │     │ 192.168.0.5  │     │    192.168.0.10      │ │
│  └──────────────┘     │   (MAAS)     │     │                      │ │
│                       │              │     │  ┌────────────────┐  │ │
│                       │  eth0 ───────┼─────┤  │    vmbr0       │  │ │
│                       │              │     │  │  (Home LAN)    │  │ │
│                       └──────────────┘     │  └────────────────┘  │ │
│                              │             │                      │ │
└──────────────────────────────┼─────────────┼──────────────────────┘ │
                               │             │                        │
┌──────────────────────────────┼─────────────┼──────────────────────┐ │
│              Provisioning Network (Isolated)                       │ │
│                      192.168.50.0/24                               │ │
│                               │             │                      │ │
│                       ┌───────┴──────┐     │  ┌────────────────┐  │ │
│                       │ USB Ethernet │     │  │    vmbr1       │  │ │
│                       │ 192.168.50.1 │─────┼──│ (Provisioning) │  │ │
│                       │ (DHCP/PXE)   │     │  └───────┬────────┘  │ │
│                       └──────────────┘     │          │           │ │
│                                            │  ┌───────┴────────┐  │ │
│                                            │  │  Proxmox VM    │  │ │
│                                            │  │  (PXE Boot)    │  │ │
│                                            │  │ 192.168.50.x   │  │ │
│                                            │  └────────────────┘  │ │
│                                            └──────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

**Key Points:**

- **Home LAN (192.168.0.0/24)**: For accessing MAAS UI and internet
- **Provisioning Network (192.168.50.0/24)**: Isolated network for PXE/DHCP, no conflict with home router
- **Raspberry Pi**: Runs MAAS controller, bridges both networks
- **Proxmox VMs**: PXE boot on provisioning network, managed by MAAS

---

## Prerequisites

### Hardware

- Raspberry Pi 4 (4GB+ RAM recommended)
- USB Ethernet adapter for the Pi (provisioning network)
- Proxmox server with spare NIC (or use existing NIC for dedicated bridge)
- Ethernet cable connecting Pi USB adapter to Proxmox spare NIC

### Software

- Ubuntu 24.04 on Raspberry Pi
- Proxmox VE 8.x on server
- Kubernetes cluster with Crossplane installed (for provider)

---

## Part 1: Hardware Setup

### 1.1 Raspberry Pi Interfaces

Your Pi will have two network interfaces:

- `eth0`: Connected to home LAN (for MAAS UI access and internet)
- `enxXXXXXXXXXXXX`: USB Ethernet adapter (for provisioning network)

Find your USB adapter name:

```bash
ip link show
```

### 1.2 Physical Connection

Connect the USB Ethernet adapter on the Pi directly to a spare NIC on the Proxmox host.

```
[Raspberry Pi USB Ethernet] ←──── Ethernet Cable ────→ [Proxmox Spare NIC]
```

---

## Part 2: MAAS Installation on Raspberry Pi

### 2.1 Install MAAS

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install MAAS (snap)
sudo snap install maas

# Initialize MAAS (region+rack controller)
sudo maas init region+rack --database-uri maas-test-db:///

# Create admin user
sudo maas createadmin --username admin --password admin --email admin@example.com
```

### 2.2 Configure Provisioning Interface

Create netplan config for the USB Ethernet adapter:

```bash
sudo nano /etc/netplan/99-provisioning.yaml
```

```yaml
network:
  version: 2
  ethernets:
    enxXXXXXXXXXXXX:  # Replace with your USB adapter name
      addresses:
        - 192.168.50.1/24
      dhcp4: false
```

Apply the configuration:

```bash
sudo netplan apply
```

Verify:

```bash
ip addr show enxXXXXXXXXXXXX
# Should show 192.168.50.1/24
```

### 2.3 Access MAAS UI

Open in browser: `http://192.168.0.5:5240/MAAS`

Login with the admin credentials you created.

---

## Part 3: Proxmox Configuration

### 3.1 Identify the Spare NIC

On Proxmox, find the NIC connected to the Pi:

```bash
ip link show
# Look for the interface that's UP when connected to Pi
```

If the NIC is bound to vfio (for passthrough), unbind it:

```bash
# Check current driver
readlink -f /sys/bus/pci/devices/0000:XX:00.0/driver

# If bound to vfio-pci, rebind to normal driver (e.g., igc for Intel I226)
echo "0000:XX:00.0" | sudo tee /sys/bus/pci/drivers/vfio-pci/unbind
echo "0000:XX:00.0" | sudo tee /sys/bus/pci/drivers/igc/bind
```

### 3.2 Create Provisioning Bridge (vmbr1)

Edit Proxmox network config:

```bash
nano /etc/network/interfaces
```

Add:

```
auto vmbr1
iface vmbr1 inet manual
    bridge-ports enp87s0  # Replace with your NIC name
    bridge-stp off
    bridge-fd 0
```

Apply:

```bash
ifreload -a
```

Verify:

```bash
ip -br link | grep vmbr1
# Should show: vmbr1 UP
```

### 3.3 Verify Link Between Pi and Proxmox

On Proxmox:

```bash
ethtool enp87s0
# Should show: Link detected: yes
```

On Pi:

```bash
ethtool enxXXXXXXXXXXXX
# Should show: Link detected: yes
```

---

## Part 4: MAAS Network Configuration

### 4.1 Sync Boot Images

In MAAS UI:

1. Go to **Images**
2. Select Ubuntu images for **amd64** architecture
3. Click **Update selection**
4. Wait for images to sync (can take 10-30 minutes)

### 4.2 Configure Provisioning Subnet

In MAAS UI:

1. Go to **Subnets**
2. Find or create subnet `192.168.50.0/24`
3. Configure:
   - **Gateway IP**: `192.168.50.1`
   - **DNS**: `8.8.8.8, 8.8.4.4`

### 4.3 Enable DHCP

1. Click on the subnet `192.168.50.0/24`
2. Go to **Reserved ranges** → Add dynamic range:
   - **Start IP**: `192.168.50.100`
   - **End IP**: `192.168.50.200`
   - **Type**: Dynamic
3. Enable DHCP:
   - **Rack controller**: Select your Pi
   - Check **MAAS provides DHCP**

### 4.4 Get API Key

In MAAS UI:

1. Click your username (top right) → **API keys**
2. Click **Generate MAAS API key**
3. Copy the key (format: `consumer:token:secret`)

---

## Part 5: Crossplane Provider Installation

### 5.1 Prerequisites

Ensure you have a Kubernetes cluster with Crossplane installed:

```bash
# Install Crossplane (if not already installed)
helm repo add crossplane-stable https://charts.crossplane.io/stable
helm repo update
helm install crossplane crossplane-stable/crossplane \
  --namespace crossplane-system --create-namespace
```

### 5.2 Install MAAS Provider

```bash
kubectl apply -f - <<EOF
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-maas
spec:
  package: ghcr.io/marouan-chak/provider-upjet-maas:latest
EOF
```

Wait for provider to be ready:

```bash
kubectl get providers -w
# Wait until HEALTHY=True
```

### 5.3 Configure Provider Credentials

Create the secret with your MAAS API key:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: maas-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    { "apiKey": "YOUR_MAAS_API_KEY" }
EOF
```

Create the ProviderConfig:

```bash
kubectl apply -f - <<EOF
apiVersion: maas.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  apiURL: "http://192.168.0.5:5240/MAAS"
  apiVersion: "2.0"
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-system
      name: maas-creds
      key: credentials
EOF
```

---

## Part 6: Applying MAAS Resources

### 6.1 Create a Proxmox VM for MAAS

In Proxmox UI:

1. **Create VM**:
   - **OS**: Do not use any media
   - **System**: Default (BIOS or UEFI)
   - **Disks**: Add a disk (e.g., 32GB)
   - **CPU**: 2 cores
   - **Memory**: 4096 MB
   - **Network**: Select `vmbr1` (provisioning bridge)

2. **Configure Boot Order**:
   - Go to VM → Options → Boot Order
   - Set **Network** first, then **Disk**

3. **Get MAC Address**:
   - Go to VM → Hardware → Network Device
   - Note the MAC address (e.g., `BC:24:11:C8:EE:50`)

4. **Start the VM**:
   - The VM will PXE boot
   - It will get an IP from MAAS DHCP
   - It will appear in MAAS UI under **Machines**

### 6.2 Apply Infrastructure Resources

```bash
# Create resource pool
kubectl apply -f - <<EOF
apiVersion: infrastructure.maas.crossplane.io/v1alpha1
kind: ResourcePool
metadata:
  name: proxmox-pool
spec:
  forProvider:
    name: proxmox-pool
    description: "Pool for Proxmox VMs managed by MAAS"
  providerConfigRef:
    name: default
EOF

# Create tag
kubectl apply -f - <<EOF
apiVersion: infrastructure.maas.crossplane.io/v1alpha1
kind: Tag
metadata:
  name: proxmox-vm-tag
spec:
  forProvider:
    name: proxmox-vm
    comment: "Tag for Proxmox VMs"
  providerConfigRef:
    name: default
EOF
```

### 6.3 Apply Network Resources (Fresh MAAS Only)

> **Note**: Skip this section if you already configured the network in MAAS UI.

```bash
# Create fabric
kubectl apply -f - <<EOF
apiVersion: network.maas.crossplane.io/v1alpha1
kind: Fabric
metadata:
  name: provisioning-fabric
spec:
  forProvider:
    name: provisioning-fabric
  providerConfigRef:
    name: default
EOF

# Create space
kubectl apply -f - <<EOF
apiVersion: network.maas.crossplane.io/v1alpha1
kind: Space
metadata:
  name: lab-space
spec:
  forProvider:
    name: lab-space
  providerConfigRef:
    name: default
EOF

# Create subnet
kubectl apply -f - <<EOF
apiVersion: network.maas.crossplane.io/v1alpha1
kind: Subnet
metadata:
  name: provisioning-subnet
spec:
  forProvider:
    cidr: "192.168.50.0/24"
    name: provisioning-subnet
    fabric: provisioning-fabric
    gatewayIp: "192.168.50.1"
    dnsServers:
      - "8.8.8.8"
      - "8.8.4.4"
    allowDns: true
    allowProxy: true
    ipRanges:
      - type: dynamic
        startIp: "192.168.50.100"
        endIp: "192.168.50.200"
        comment: "DHCP range for PXE boot"
  providerConfigRef:
    name: default
EOF
```

### 6.4 Apply Machine Resource

First, create the power parameters secret:

```bash
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
metadata:
  name: proxmox-vm-power-params
  namespace: crossplane-system
type: Opaque
stringData:
  params: "{}"
EOF
```

Then create the machine resource:

```bash
kubectl apply -f - <<EOF
apiVersion: machine.maas.crossplane.io/v1alpha1
kind: Machine
metadata:
  name: proxmox-vm
spec:
  forProvider:
    # Replace with your Proxmox VM's MAC address
    pxeMacAddress: "BC:24:11:C8:EE:50"
    powerType: manual
    powerParametersSecretRef:
      name: proxmox-vm-power-params
      namespace: crossplane-system
      key: params
    hostname: proxmox-vm
    architecture: amd64/generic
  providerConfigRef:
    name: default
EOF
```

### 6.5 Verify Resources

```bash
# Check all managed resources
kubectl get managed

# Check specific resource types
kubectl get fabrics.network.maas.crossplane.io
kubectl get subnets.network.maas.crossplane.io
kubectl get machines.machine.maas.crossplane.io
kubectl get resourcepools.infrastructure.maas.crossplane.io
kubectl get tags.infrastructure.maas.crossplane.io
```

---

## Troubleshooting

### VM Not Appearing in MAAS

1. **Check DHCP is working**:

   ```bash
   # On Pi, watch for DHCP requests
   sudo tcpdump -i enxXXXXXXXXXXXX port 67 or port 68
   ```

2. **Check PXE/TFTP traffic**:

   ```bash
   sudo tcpdump -i enxXXXXXXXXXXXX port 69
   ```

3. **Verify VM network**:
   - Ensure VM is connected to `vmbr1`
   - Ensure boot order has Network first

### Provider Not Connecting to MAAS

1. **Check API key**:

   ```bash
   kubectl get secret maas-creds -n crossplane-system -o yaml
   ```

2. **Check ProviderConfig**:

   ```bash
   kubectl describe providerconfig default
   ```

3. **Check provider logs**:

   ```bash
   kubectl logs -n crossplane-system -l pkg.crossplane.io/provider=provider-upjet-maas
   ```

### Resource Stuck in "Creating"

1. **Check resource status**:

   ```bash
   kubectl describe <resource-type> <resource-name>
   ```

2. **Look at Conditions and Events** for error messages.

### Reference Resolution Errors

If you see `referenced field was empty`:

- The referenced resource doesn't exist yet
- The referenced resource hasn't been assigned an external-name by MAAS
- Apply resources in order: infrastructure → network → machine

---

## Resource Dependency Order

When applying all resources, follow this order:

```
1. Secret (credentials)
2. ProviderConfig
3. ResourcePool, Tag (no dependencies)
4. Space (no dependencies)
5. Fabric (no dependencies)
6. VLAN (depends on Fabric)
7. Subnet (depends on Fabric, optionally VLAN)
8. SubnetIPRange (depends on Subnet)
9. Machine (depends on ResourcePool optionally)
10. InterfacePhysical (depends on Machine)
11. InterfaceLink (depends on Machine, Subnet)
12. InterfaceBond/Bridge/VLAN (depends on Machine)
13. BlockDevice (depends on Machine)
```

---

## Quick Reference

| Component | IP/Value |
|-----------|----------|
| Pi Home LAN IP | 192.168.0.5 |
| Pi Provisioning IP | 192.168.50.1 |
| Provisioning Subnet | 192.168.50.0/24 |
| DHCP Range | 192.168.50.100-200 |
| MAAS UI | <http://192.168.0.5:5240/MAAS> |
| Proxmox Host | 192.168.0.10 |
| Proxmox Bridge | vmbr1 |

---

## Next Steps

After completing this guide, you can:

1. **Commission the machine**: In MAAS UI, select the machine and click "Commission"
2. **Deploy an OS**: After commissioning, deploy Ubuntu or another OS
3. **Add more VMs**: Create additional Proxmox VMs and manage them via Crossplane
4. **Use compositions**: Create Crossplane Compositions to provision complete environments
