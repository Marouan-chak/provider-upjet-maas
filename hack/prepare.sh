#!/usr/bin/env bash
set -euox pipefail

prompt_if_unset() {
  local var_name="$1"
  local prompt="$2"
  local -n var_ref="${var_name}"

  if [ -z "${var_ref:-}" ]; then
    read -r -p "${prompt}" var_ref
  fi
}

prompt_if_unset PROVIDER_NAME_LOWER "Lower case provider name (ex. github): "
prompt_if_unset PROVIDER_NAME_NORMAL "Normal case provider name (ex. GitHub): "
prompt_if_unset ORGANIZATION_NAME "Organization (e.g., my-org-name): "
prompt_if_unset CRD_ROOT_GROUP "CRD rootGroup (e.g., crossplane.io): "

REPLACE_FILES='./* ./.github :!build/** :!go.* :!hack/prepare.sh'
# shellcheck disable=SC2086
git grep -l 'template' -- ${REPLACE_FILES} | xargs sed -i.bak "s/upjet-provider-template/provider-${PROVIDER_NAME_LOWER}/g"
# shellcheck disable=SC2086
git grep -l 'template' -- ${REPLACE_FILES} | xargs sed -i.bak "s/template/${PROVIDER_NAME_LOWER}/g"
# shellcheck disable=SC2086
git grep -l "crossplane/provider-${PROVIDER_NAME_LOWER}" -- ${REPLACE_FILES} | xargs sed -i.bak "s|crossplane/provider-${PROVIDER_NAME_LOWER}|${ORGANIZATION_NAME}/provider-${PROVIDER_NAME_LOWER}|g"
# shellcheck disable=SC2086
git grep -l 'Template' -- ${REPLACE_FILES} | xargs sed -i.bak "s/Template/${PROVIDER_NAME_NORMAL}/g"
# shellcheck disable=SC2086
git grep -l "crossplane.io" -- "apis/cluster/v1*" | xargs sed -i.bak "s|crossplane.io|${CRD_ROOT_GROUP}|g"
git grep -l "crossplane.io" -- "apis/namespaced/v1*" | xargs sed -i.bak "s|crossplane.io|${CRD_ROOT_GROUP}|g"
# shellcheck disable=SC2086
git grep -l "crossplane.io" -- "cluster/test/setup.sh" | xargs sed -i.bak "s|crossplane.io|${CRD_ROOT_GROUP}|g"
# shellcheck disable=SC2086
git grep -l "ujconfig\.WithRootGroup(\"${PROVIDER_NAME_LOWER}\.crossplane\.io\")" -- "config/provider.go" | xargs sed -i.bak "s|ujconfig.WithRootGroup(\"${PROVIDER_NAME_LOWER}\.crossplane\.io\")|ujconfig.WithRootGroup(\"${PROVIDER_NAME_LOWER}.${CRD_ROOT_GROUP}\")|g"
git grep -l "ujconfig\.WithRootGroup(\"${PROVIDER_NAME_LOWER}\.m\.crossplane\.io\")" -- "config/provider.go" | xargs sed -i.bak "s|ujconfig.WithRootGroup(\"${PROVIDER_NAME_LOWER}\.m.crossplane\.io\")|ujconfig.WithRootGroup(\"${PROVIDER_NAME_LOWER}.m.${CRD_ROOT_GROUP}\")|g"

# We need to be careful while replacing "template" keyword in go.mod as it could tamper
# some imported packages under require section.
sed -i.bak "s|crossplane/upjet-provider-template|${ORGANIZATION_NAME}/provider-${PROVIDER_NAME_LOWER}|g" go.mod
sed -i.bak -e "s|PROJECT_REPO ?= github.com/crossplane/|PROJECT_REPO ?= github.com/${ORGANIZATION_NAME}/|g" -e "s|\(blob/main/internal/\)${PROVIDER_NAME_LOWER}s|\1templates|g" Makefile
sed -i.bak "s/\[YEAR\]/$(date +%Y)/g" LICENSE

# Clean up the .bak files created by sed
git clean -fd

git mv "internal/clients/template.go" "internal/clients/${PROVIDER_NAME_LOWER}.go"
git mv "cluster/images/upjet-provider-template" "cluster/images/provider-${PROVIDER_NAME_LOWER}"

# We need to remove this api folder otherwise first `make generate` fails with
# the following error probably due to some optimizations in go generate with v1.17:
# generate: open /Users/hasanturken/Workspace/crossplane-contrib/upjet-provider-template/apis/null/v1alpha1/zz_generated.deepcopy.go: no such file or directory
rm -rf apis/cluster/null
rm -rf apis/namespaced/null
# remove the sample directory which was a configuration in the template
rm -rf config/cluster/null
rm -rf config/namespaced/null
# remove the sample MR example from the template
rm -rf examples/cluster/null
rm -rf examples/namespaced/null
