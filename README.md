# Provider MAAS

`provider-upjet-maas` is a [Crossplane](https://crossplane.io/) provider for
[Canonical MAAS](https://maas.io/) that is built using [Upjet](https://github.com/crossplane/upjet)
code generation tools and exposes XRM-conformant managed resources for the MAAS API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/Marouan-chak/provider-upjet-maas):
```
up ctp provider install Marouan-chak/provider-upjet-maas:v0.1.0
```

Alternatively, you can use the Provider manifest below, editing the version as necessary:
```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-upjet-maas
spec:
  package: ghcr.io/Marouan-chak/provider-upjet-maas:v0.1.0
```

## Developing

Run code-generation pipeline:
```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:
```console
make run
```

Build, push, and install:
```console
make all
```

Build binary:
```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/Marouan-chak/provider-upjet-maas/issues).
