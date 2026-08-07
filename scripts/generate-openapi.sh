#!/usr/bin/env bash
set -euo pipefail

# Regenerates docs/openapi.yaml from swaggo annotations in the Go sources.
# Runs inside the golang image; see the api-spec target in the Makefile.

SWAG_VERSION="${SWAG_VERSION:-v2.0.0-rc5}"
REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BACKEND_DIR="${BACKEND_DIR:-$REPO_ROOT/src/backend}"
DOCS_DIR="${DOCS_DIR:-$REPO_ROOT/docs}"
WORK_DIR="$(mktemp -d)"

trap 'rm -rf "$WORK_DIR"' EXIT

if [ ! -x "$(go env GOPATH)/bin/swag" ]; then
  go install "github.com/swaggo/swag/v2/cmd/swag@${SWAG_VERSION}"
fi

SWAG="$(go env GOPATH)/bin/swag"

cd "$BACKEND_DIR"

# swag logs a "TypeSpecDef is nil" line per stdlib type it cannot resolve and a
# warning per non-decimal const; neither affects the output, so both are hidden.
"$SWAG" init \
  --v3.1 \
  -g cmd/api/main.go \
  -d ./ \
  --parseInternal \
  --outputTypes yaml \
  -o "$WORK_DIR" 2>&1 |
  grep -viE 'TypeSpecDef is nil|failed to evaluate const|no Go files in' || true

if [ ! -f "$WORK_DIR/swagger.yaml" ]; then
  echo "swag did not produce a spec, see the output above" >&2
  exit 1
fi

go run -tags tools ./tools/openapigen \
  "$WORK_DIR/swagger.yaml" \
  "$DOCS_DIR/openapi.overlay.yaml" \
  "$DOCS_DIR/openapi.yaml"

echo "docs/openapi.yaml regenerated"
