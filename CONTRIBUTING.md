# Contributing to Provider MAAS

Thank you for your interest in contributing to provider-upjet-maas! This document provides guidelines for contributing to the project.

## Getting Started

### Prerequisites

- Go 1.24+
- Docker
- kubectl
- A Kubernetes cluster (or Kind for local development)
- Crossplane installed in your cluster

### Setting Up Development Environment

1. **Clone the repository**

   ```bash
   git clone https://github.com/Marouan-chak/provider-upjet-maas.git
   cd provider-upjet-maas
   ```

2. **Initialize submodules**

   ```bash
   git submodule sync && git submodule update --init --recursive
   ```

3. **Install dependencies**

   ```bash
   go mod download
   ```

4. **Run code generation**

   ```bash
   go run cmd/generator/main.go "$PWD"
   ```

5. **Build the provider**

   ```bash
   make build
   ```

## Development Workflow

### Setting Up Pre-commit Hooks

We recommend setting up pre-commit hooks to catch issues early:

**Option 1: Using pre-commit (recommended)**

```bash
pip install pre-commit
pre-commit install
```

**Option 2: Using our install script**

```bash
./scripts/install-hooks.sh
```

### Making Changes

1. **Create a feature branch**

   ```bash
   git checkout -b feature/my-feature
   ```

2. **Make your changes**
   - Follow the existing code style
   - Add tests for new functionality
   - Update documentation as needed

3. **Run quick local checks before committing**

   ```bash
   make check
   ```

   This runs:
   - `make fmt` - Format Go code
   - `make vet` - Run go vet
   - `make tidy` - Tidy go modules
   - `make lint-fast` - Quick lint checks

4. **Run tests**

   ```bash
   go test ./...
   ```

5. **Regenerate if you modified config/**

   ```bash
   go run cmd/generator/main.go "$PWD"
   ```

6. **Verify generation is deterministic**

   ```bash
   make check-generate
   ```

### Adding New Resources

1. Configure the resource in `config/` package using ResourceConfigurators
2. Run code generation: `go run cmd/generator/main.go "$PWD"`
3. Add example YAML in `examples/resources/`
4. Update README.md with the new resource

### Testing Locally

Use the provided script to set up a local Kind cluster:

```bash
./scripts/setup-kind.sh
```

To clean up:

```bash
./scripts/setup-kind.sh --cleanup
```

## Pull Request Process

1. **Ensure your PR**:
   - Has a clear description of the changes
   - Includes relevant tests
   - Has updated documentation
   - Passes all CI checks

2. **PR Title Format**:
   - `feat: Add support for X resource`
   - `fix: Correct reference resolution for Y`
   - `docs: Update installation instructions`
   - `chore: Update dependencies`

3. **Code Review**:
   - All PRs require at least one maintainer approval
   - Address review feedback promptly
   - Keep PRs focused and reasonably sized

## Reporting Issues

When reporting issues, please include:

- Provider version
- Crossplane version
- MAAS version
- Kubernetes version
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs and resource YAML

## Code of Conduct

Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## Questions?

- Open an issue for bugs or feature requests
- Check existing issues before creating new ones

## License

By contributing to this project, you agree that your contributions will be licensed under the Apache 2.0 License.
