#!/usr/bin/env bash
# generate-version.sh — Generate OpenAPI client + scaffold service code for a Slurm version.
#
# Usage:
#   ./hack/generate-version.sh v0.0.45 path/to/openapi-spec.json
#
# This script:
#   1. Generates api/v0045/ from the OpenAPI spec using oapi-codegen
#   2. Verifies the generated code compiles
#
# After running, manually update:
#   - version.go: add V0045 constant
#   - services.go: add version-specific types if needed
#
# Prerequisites:
#   go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

set -euo pipefail

VERSION="${1:?Usage: $0 <version> <openapi-spec.json>}"
SPEC="${2:?Usage: $0 <version> <openapi-spec.json>}"

# Normalize: v0.0.45 → v0045
DIR_NAME=$(echo "$VERSION" | sed 's/v0\.0\./v00/')

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
API_DIR="$ROOT_DIR/api/$DIR_NAME"

if [ -d "$API_DIR" ]; then
    echo "WARNING: $API_DIR already exists. Overwriting."
fi

mkdir -p "$API_DIR"

echo "==> Generating api/$DIR_NAME from $SPEC..."

# Generate the client code
cat > "$API_DIR/api.go" << GOEOF
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=oapi-codegen-config.yaml $SPEC
package $DIR_NAME
GOEOF

# Create oapi-codegen config
cat > "$API_DIR/oapi-codegen-config.yaml" << YAMLEOF
package: $DIR_NAME
output: slurm.gen.go
generate:
  models: true
  client: false
  embedded-spec: false
YAMLEOF

# Run code generation
echo "==> Running oapi-codegen..."
cd "$API_DIR"
if command -v oapi-codegen &> /dev/null; then
    oapi-codegen --config=oapi-codegen-config.yaml "$SPEC"
    echo "==> Generated $API_DIR/slurm.gen.go"
else
    echo "ERROR: oapi-codegen not found. Install with:"
    echo "  go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest"
    exit 1
fi

cd "$ROOT_DIR"

# Verify build
echo "==> Verifying build..."
if go build ./api/"$DIR_NAME"/...; then
    echo "==> Build OK"
else
    echo "ERROR: Generated code does not compile"
    exit 1
fi

echo ""
echo "=== Done ==="
echo ""
echo "Next steps:"
echo "  1. Add V00XX constant to version.go"
echo "  2. Run: go mod tidy"
echo "  3. Run: go test ./..."
