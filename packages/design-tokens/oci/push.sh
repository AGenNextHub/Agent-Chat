#!/usr/bin/env bash
# Push the design tokens as a single, lightweight OCI artifact via ORAS.
# A tool/artifact only — no Penpot app, no database. Works with any OCI registry.
#
# Usage: oci/push.sh <registry>/<repo>:<tag>
#   e.g. oci/push.sh ghcr.io/agennext/design-tokens:0.1.0
set -euo pipefail

REF="${1:?usage: push.sh <registry>/<repo>:<tag>}"
HERE="$(cd "$(dirname "$0")" && pwd)"
DIST="$HERE/../dist"

node "$HERE/../build.mjs" # refresh dist/ from the canonical source

oras push "$REF" \
  --artifact-type application/vnd.agennext.design-tokens.v1 \
  "$DIST/tokens.json:application/vnd.agennext.tokens+json" \
  "$DIST/tokens.css:text/css" \
  "$DIST/logo.svg:image/svg+xml"

echo "pushed OCI artifact: $REF   (pull with: oras pull $REF)"
