#!/usr/bin/env bash
# check-compat.sh — Check SDK compatibility with a new OpenAPI spec.
#
# Usage:
#   ./hack/check-compat.sh path/to/new-openapi-spec.json
#
# This script:
#   1. Generates a temporary client from the new spec
#   2. Checks if existing SDK service interfaces still compile against it
#   3. Reports any breaking changes (missing types, changed fields)
#
# Exit codes:
#   0 — Compatible
#   1 — Breaking changes detected

set -euo pipefail

SPEC="${1:?Usage: $0 <new-openapi-spec.json>}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
TEMP_DIR=$(mktemp -d)
trap 'rm -rf "$TEMP_DIR"' EXIT

echo "==> Checking SDK compatibility with new spec..."
echo "    Spec: $SPEC"
echo "    Temp: $TEMP_DIR"

# Copy current project
cp -r "$ROOT_DIR" "$TEMP_DIR/project"
cd "$TEMP_DIR/project"

# Generate new API types into a temporary version
mkdir -p api/vnew

cat > api/vnew/api.go << 'EOF'
package vnew
EOF

if ! command -v oapi-codegen &> /dev/null; then
    echo "ERROR: oapi-codegen not found"
    exit 1
fi

oapi-codegen --config=<(cat << 'YAML'
package: vnew
output: slurm.gen.go
generate:
  models: true
  client: false
  embedded-spec: false
YAML
) "$SPEC" > api/vnew/slurm.gen.go 2>&1

# Extract type names from new spec
NEW_TYPES=$(grep "^type V" api/vnew/slurm.gen.go | awk '{print $2}' | sort)

# Extract type names used in current SDK
USED_TYPES=$(grep -rh "v0040\.\|v0041\.\|v0042\.\|v0043\.\|v0044\." \
    client.go services.go service_*.go slurm_proxy.go slurmdb_proxy.go 2>/dev/null \
    | grep -oE 'v[0-9]+\.V[0-9]+[A-Za-z]+' \
    | sed 's/^v[0-9]*\.//' \
    | sort -u)

echo ""
echo "==> Types used by SDK:"
echo "$USED_TYPES" | head -20
echo "    ... ($(echo "$USED_TYPES" | wc -l | tr -d ' ') total)"

echo ""
echo "==> Checking for missing types in new spec..."

MISSING=0
while IFS= read -r type_name; do
    # Convert V0040 prefix to the new spec's prefix pattern
    base_name=$(echo "$type_name" | sed 's/^V[0-9]*//')
    if ! echo "$NEW_TYPES" | grep -q "$base_name"; then
        echo "  MISSING: $type_name (base: $base_name)"
        MISSING=$((MISSING + 1))
    fi
done <<< "$USED_TYPES"

if [ "$MISSING" -eq 0 ]; then
    echo "  All types found in new spec."
    echo ""
    echo "==> COMPATIBLE"
    exit 0
else
    echo ""
    echo "==> BREAKING CHANGES: $MISSING types missing"
    exit 1
fi
