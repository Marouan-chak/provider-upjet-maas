#!/usr/bin/env python3

import json
import os
import re
import sys

def _resources_from_provider_metadata_yaml(path: str) -> list[str]:
    """
    Extract Terraform resource names from `config/provider-metadata.yaml` without
    requiring external YAML deps. We only need the keys under the top-level
    `resources:` map, which are formatted like:

    resources:
        maas_machine:
            ...
    """
    resources: list[str] = []
    in_resources = False
    key_re = re.compile(r"^\s{4}([A-Za-z0-9_]+):\s*$")
    with open(path, "r", encoding="utf-8") as f:
        for line in f:
            if not in_resources:
                if line.strip() == "resources:":
                    in_resources = True
                continue

            # stop once we leave the `resources:` section
            if line.strip() and not line.startswith("    "):
                break

            m = key_re.match(line)
            if m:
                resources.append(m.group(1))
    return resources


def _load_resource_list(path: str) -> list[str]:
    """
    Accept either:
    - JSON list file (legacy: `config/generated.lst`)
    - `config/provider-metadata.yaml` (preferred, checked in)
    """
    if not os.path.exists(path):
        raise FileNotFoundError(path)

    lower = path.lower()
    if lower.endswith((".yaml", ".yml")):
        return _resources_from_provider_metadata_yaml(path)

    with open(path, encoding="utf-8") as f:
        return json.load(f)


# usage: version_diff.py <resource list file> <base JSON schema path> <bumped JSON schema path>
# resource list file can be either:
# - config/provider-metadata.yaml (preferred)
# - config/generated.lst (legacy JSON list)
# example usage: version_diff.py config/provider-metadata.yaml .work/schema.json.3.38.0 config/schema.json
if __name__ == "__main__":
    base_path = sys.argv[2]
    bumped_path = sys.argv[3]
    print(f'Reporting schema changes between "{base_path}" as base version and "{bumped_path}" as bumped version')
    resources = _load_resource_list(sys.argv[1])
    with open(base_path) as f:
        base = json.load(f)
    with open(bumped_path) as f:
        bump = json.load(f)

    provider_name = None
    for k in base["provider_schemas"]:
        # the first key is the provider name
        provider_name = k
        break
    if provider_name is None:
        print(f"Cannot extract the provider name from the base schema: {base_path}")
        sys.exit(-1)
    base_schemas = base["provider_schemas"][provider_name]["resource_schemas"]
    bumped_schemas = bump["provider_schemas"][provider_name]["resource_schemas"]

    for name in resources:
        try:
            if base_schemas[name]["version"] != bumped_schemas[name]["version"]:
                print(f'{name}:{base_schemas[name]["version"]}-{bumped_schemas[name]["version"]}')
        except KeyError as ke:
            print(f'{name} is not found in schema: {ke}')
            continue
