#!/bin/bash
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

CLUSTER_NAME="${CLUSTER_NAME:-crossplane-test}"
CROSSPLANE_VERSION="${CROSSPLANE_VERSION:-2.1.3}"
PROVIDER_IMAGE_PATTERN="${PROVIDER_IMAGE_PATTERN:-provider-upjet-maas-amd64}"
PROVIDER_IMAGE=""
# Canonical image name for Crossplane (must be fully qualified)
LOCAL_IMAGE="xpkg.local/provider-upjet-maas:latest"

echo -e "${GREEN}=== MAAS Provider Setup Script ===${NC}"
echo ""

# Check prerequisites
check_prereqs() {
  echo -e "${YELLOW}Checking prerequisites...${NC}"

  for cmd in kind kubectl helm docker; do
    if ! command -v $cmd &>/dev/null; then
      echo -e "${RED}Error: $cmd is not installed${NC}"
      exit 1
    fi
    echo "  - $cmd: OK"
  done
  echo ""
}

# Create Kind cluster
create_cluster() {
  echo -e "${YELLOW}Creating Kind cluster: ${CLUSTER_NAME}...${NC}"

  if kind get clusters 2>/dev/null | grep -q "^${CLUSTER_NAME}$"; then
    echo "  Cluster already exists, skipping creation"
  else
    kind create cluster --name "${CLUSTER_NAME}" --wait 5m
  fi

  kubectl cluster-info --context "kind-${CLUSTER_NAME}"
  echo ""
}

# Install Crossplane
install_crossplane() {
  echo -e "${YELLOW}Installing Crossplane ${CROSSPLANE_VERSION}...${NC}"

  helm repo add crossplane-stable https://charts.crossplane.io/stable 2>/dev/null || true
  helm repo update

  if kubectl get namespace crossplane-system &>/dev/null; then
    echo "  Crossplane namespace exists, checking installation..."
    if helm status crossplane -n crossplane-system &>/dev/null; then
      echo "  Crossplane already installed, skipping"
    else
      helm install crossplane crossplane-stable/crossplane \
        --namespace crossplane-system \
        --create-namespace \
        --version "${CROSSPLANE_VERSION}" \
        --wait
    fi
  else
    helm install crossplane crossplane-stable/crossplane \
      --namespace crossplane-system \
      --create-namespace \
      --version "${CROSSPLANE_VERSION}" \
      --wait
  fi

  echo "  Waiting for Crossplane pods to be ready..."
  kubectl wait --for=condition=ready pod -l app=crossplane -n crossplane-system --timeout=120s
  echo ""
}

# Find provider image
find_provider_image() {
  # Find the most recent image matching the pattern
  PROVIDER_IMAGE=$(docker images --format "{{.Repository}}:{{.Tag}}" | grep "${PROVIDER_IMAGE_PATTERN}" | head -1)
  if [ -n "${PROVIDER_IMAGE}" ]; then
    echo "  Found image: ${PROVIDER_IMAGE}"
    return 0
  fi
  return 1
}

# Build provider image
build_provider() {
  echo -e "${YELLOW}Building provider image...${NC}"

  cd "$(dirname "$0")/.."

  if find_provider_image; then
    read -r -p "  Rebuild? (y/N): " rebuild
    if [[ "$rebuild" =~ ^[Yy]$ ]]; then
      make build
      find_provider_image
    fi
  else
    echo "  No existing image found, running make build..."
    make build
    find_provider_image
  fi
  echo ""
}

# Load image into Kind
load_image() {
  echo -e "${YELLOW}Loading provider image into Kind...${NC}"

  if [ -z "${PROVIDER_IMAGE}" ]; then
    echo -e "${RED}Error: PROVIDER_IMAGE is not set${NC}"
    exit 1
  fi

  # Retag to a fully qualified name for Crossplane
  echo "  Retagging ${PROVIDER_IMAGE} -> ${LOCAL_IMAGE}..."
  docker tag "${PROVIDER_IMAGE}" "${LOCAL_IMAGE}"

  echo "  Loading ${LOCAL_IMAGE}..."
  kind load docker-image "${LOCAL_IMAGE}" --name "${CLUSTER_NAME}"
  echo ""
}

# Install provider
install_provider() {
  echo -e "${YELLOW}Installing MAAS provider...${NC}"

  cd "$(dirname "$0")/.."

  # Apply CRDs first
  echo "  Applying CRDs..."
  kubectl apply -f package/crds/

  # Apply provider configuration
  echo "  Applying provider..."
  kubectl apply -f examples/provider/provider.yaml

  echo "  Waiting for provider to be healthy..."
  sleep 5
  kubectl wait --for=condition=healthy providers.pkg.crossplane.io/provider-upjet-maas --timeout=120s || true
  echo ""
}

# Apply credentials
apply_credentials() {
  echo -e "${YELLOW}Applying credentials and ProviderConfig...${NC}"

  cd "$(dirname "$0")/.."

  kubectl apply -f examples/providerconfig/creds.yaml
  kubectl apply -f examples/providerconfig/providerconfig.yaml
  echo ""
}

# Show status
show_status() {
  echo -e "${GREEN}=== Setup Complete ===${NC}"
  echo ""
  echo "Cluster: ${CLUSTER_NAME}"
  echo ""
  echo "Provider status:"
  kubectl get providers.pkg.crossplane.io
  echo ""
  echo "ProviderConfig:"
  kubectl get providerconfigs.maas.crossplane.io 2>/dev/null || echo "  (none configured yet)"
  echo ""
  echo -e "${YELLOW}Next steps:${NC}"
  echo "  1. Update examples/providerconfig/secret.yaml with your MAAS API key"
  echo "  2. Update examples/providerconfig/providerconfig.yaml with your MAAS API URL"
  echo "  3. Apply: kubectl apply -f examples/providerconfig/"
  echo "  4. Test with: kubectl apply -f examples/resources/network/fabric.yaml"
  echo "  5. Watch: kubectl get fabrics.network.maas.crossplane.io -w"
  echo ""
}

# Cleanup function
cleanup() {
  echo -e "${YELLOW}Cleaning up Kind cluster...${NC}"
  kind delete cluster --name "${CLUSTER_NAME}"
}

# Main
main() {
  case "${1:-}" in
  --cleanup)
    cleanup
    ;;
  --skip-build)
    check_prereqs
    create_cluster
    install_crossplane
    echo -e "${YELLOW}Finding provider image...${NC}"
    if ! find_provider_image; then
      echo -e "${RED}Error: No provider image found matching '${PROVIDER_IMAGE_PATTERN}'${NC}"
      echo "  Run 'make build' first or omit --skip-build"
      exit 1
    fi
    load_image
    install_provider
    apply_credentials
    show_status
    ;;
  *)
    check_prereqs
    create_cluster
    install_crossplane
    build_provider
    load_image
    install_provider
    apply_credentials
    show_status
    ;;
  esac
}

main "$@"
