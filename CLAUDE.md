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

# Build the provider
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
