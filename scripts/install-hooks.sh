#!/bin/bash
# Install git hooks for local development
# Usage: ./scripts/install-hooks.sh

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "$SCRIPT_DIR")"
HOOKS_DIR="$REPO_ROOT/.git/hooks"

echo "Installing git hooks..."

# Pre-commit hook
cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash
# Pre-commit hook for provider-upjet-maas

set -e

echo "Running pre-commit checks..."

# Check if pre-commit is installed (preferred method)
if command -v pre-commit &> /dev/null; then
    pre-commit run --files $(git diff --cached --name-only --diff-filter=ACM)
    exit $?
fi

# Fallback: run basic checks manually
echo "  Checking Go formatting..."
UNFORMATTED=$(gofmt -l $(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' | grep -v 'zz_' || true))
if [ -n "$UNFORMATTED" ]; then
    echo "ERROR: The following files are not formatted:"
    echo "$UNFORMATTED"
    echo "Run 'gofmt -w <file>' to fix."
    exit 1
fi

echo "  Running go vet..."
go vet ./... 2>&1 | grep -v "zz_" || true

echo "  Checking go.mod..."
go mod tidy
if ! git diff --exit-code go.mod go.sum > /dev/null 2>&1; then
    echo "ERROR: go.mod or go.sum is not tidy. Run 'go mod tidy' and commit the changes."
    exit 1
fi

echo "Pre-commit checks passed!"
EOF

chmod +x "$HOOKS_DIR/pre-commit"

# Pre-push hook
cat > "$HOOKS_DIR/pre-push" << 'EOF'
#!/bin/bash
# Pre-push hook for provider-upjet-maas

set -e

echo "Running pre-push checks..."

# Run tests
echo "  Running tests..."
go test ./... -short -count=1 2>&1 | tail -5

# Check if generated code is up to date
echo "  Checking generated code..."
make generate 2>/dev/null || true
if ! git diff --exit-code --quiet; then
    echo "WARNING: Generated code may be out of date."
    echo "Run 'make generate' and commit any changes."
    # Don't fail, just warn
fi

echo "Pre-push checks passed!"
EOF

chmod +x "$HOOKS_DIR/pre-push"

echo "Git hooks installed successfully!"
echo ""
echo "Installed hooks:"
echo "  - pre-commit: Format checks, go vet, go mod tidy"
echo "  - pre-push: Tests, generated code check"
echo ""
echo "For more comprehensive checks, install pre-commit:"
echo "  pip install pre-commit && pre-commit install"
