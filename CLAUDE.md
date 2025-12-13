# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an Upjet-based Crossplane provider for MAAS (Canonical's Metal-as-a-Service). Upjet generates Crossplane providers from Terraform providers, reducing hand-written controller code.

The provider wraps Canonical's Terraform provider for MAAS: https://github.com/canonical/terraform-provider-maas

## Build Commands

```bash
# Initialize submodules (required first time)
git submodule sync && git submodule update --init --recursive

# Run the code generation pipeline
go run cmd/generator/main.go "$PWD"

# Build the provider binary and container image
make build

# Run the controller locally against a cluster
make run

# Run all standard targets
make all

# Run unit tests
go test ./...

# Regenerate after config or upstream changes
make generate
```

## Building and Pushing the Crossplane Package

Crossplane providers are distributed as OCI packages (xpkg), not regular Docker images. The xpkg contains the provider metadata (`crossplane.yaml`), CRDs, and references the controller image.

### Prerequisites

- Docker logged in to your registry: `docker login`
- Crossplane CLI installed: `curl -sL https://raw.githubusercontent.com/crossplane/crossplane/master/install.sh | sh`

### Build Steps

```bash
# 1. Initialize submodules (first time only)
git submodule sync && git submodule update --init --recursive

# 2. Generate code and build everything (binary + container image + xpkg)
make build

# 3. The xpkg file will be created at:
#    _output/xpkg/linux_amd64/provider-upjet-maas-<version>.xpkg
ls -la _output/xpkg/linux_amd64/
```

### Push to Registry

```bash
# Option 1: Using crossplane CLI directly
crossplane xpkg push \
  --package-files="_output/xpkg/linux_amd64/provider-upjet-maas-<version>.xpkg" \
  docker.io/<your-username>/provider-maas:<tag>

# Option 2: Using make target (pushes to configured registries)
make xpkg.push XPKG_REG_ORGS=docker.io/<your-username>

# Example with specific version:
crossplane xpkg push \
  --package-files="_output/xpkg/linux_amd64/provider-upjet-maas-v0.0.0-1.g1e3a221.dirty.xpkg" \
  docker.io/marouandock/provider-maas:v0.1.0
```

### Install in Cluster

```bash
# 1. Apply CRDs
kubectl apply -f package/crds/

# 2. Apply provider configuration
kubectl apply -f examples/provider/provider.yaml

# 3. Check provider status
kubectl get providers.pkg.crossplane.io

# 4. Configure credentials
kubectl apply -f examples/providerconfig/secret.yaml
kubectl apply -f examples/providerconfig/providerconfig.yaml
```

### Local Development with Kind

For local testing without pushing to a registry, use the setup script:

```bash
# Full setup (builds and deploys to Kind)
./scripts/setup-kind.sh

# Skip build if image already exists
./scripts/setup-kind.sh --skip-build

# Cleanup
./scripts/setup-kind.sh --cleanup
```

## Architecture

### Upjet Provider Structure

- `cmd/generator/` - Code generation entry point
- `cmd/provider/` - Provider main entry point
- `config/` - Resource configurators that control API group naming, kind naming, and generation details per Terraform resource
- `apis/` - Generated API types including ProviderConfig
- `internal/clients/` - Provider credential handling (maas.go)
- `internal/controller/` - Generated controllers

### ProviderConfig Design

Authentication uses explicit ProviderConfig fields for non-sensitive settings with the API key stored in a Kubernetes Secret:

**ProviderConfig fields:**
- `apiURL` - MAAS API endpoint (required)
- `apiVersion` - MAAS API version (default: "2.0")
- `installationMethod` - optional, "snap" or other

**Secret format:** JSON object stored under one key:
```yaml
stringData:
  credentials: |
    { "apiKey": "YOUR_MAAS_API_KEY" }
```

Maps to Terraform provider config:
```hcl
provider "maas" {
  api_url = "..."
  api_version = "..."
  api_key = "..."
  installation_method = "..."
}
```

### API Groups

- Cluster-scoped resources: `maas.crossplane.io`
- Namespaced resources: `maas.m.crossplane.io`

## Key Workflows

### Adding/Modifying Resources

1. Configure resources in the `config/` package using ResourceConfigurators
2. Run `go run cmd/generator/main.go "$PWD"` to regenerate
3. Verify regeneration is deterministic (running twice produces no diff)

### Upgrading Terraform Provider

1. Bump `TERRAFORM_PROVIDER_VERSION` in Makefile
2. Regenerate with `make generate` and `go run cmd/generator/main.go "$PWD"`
3. Fix any compilation or schema edge cases
4. Cut a new provider release

### Generation Check

CI should fail if generated code differs from committed code. Generation must be deterministic.

## References

- Upjet provider template: https://github.com/crossplane/upjet-provider-template
- Upjet project: https://github.com/crossplane/upjet
- MAAS Terraform provider: https://registry.terraform.io/providers/canonical/maas/latest/docs
- Crossplane provider packages: https://docs.crossplane.io/latest/packages/providers/
- Safe start guide: https://docs.crossplane.io/latest/guides/implementing-safe-start/
