#!/usr/bin/env bash
set -euo pipefail

LIMIT=250
FAILED=0

FILES=$(git ls-files \
  'src/backend/**/*.go' \
  'src/frontend/app/**/*.ts' 'src/frontend/app/**/*.tsx' \
  'src/frontend/app/**/*.css' 'src/frontend/app/**/*.json' \
  'src/frontend/scripts/**/*.mjs' \
  'scripts/**/*.mjs' 'scripts/**/*.sh' \
  'tests/**/*.go' 'tests/**/*.ts' \
  | grep -v -E '(node_modules/|/dist/|/build/|\.react-router/|/generated/|\.generated\.|_gen\.go|\.pb\.go|mock_|/mocks/|package-lock\.json|go\.sum)' || true)

if [ -z "$FILES" ]; then
  echo "No source files matched the file length check, verify the globs in $0"
  exit 1
fi

OVERSIZED=$(printf '%s\n' "$FILES" | tr '\n' '\0' | xargs -0 awk -v limit="$LIMIT" '
  FNR == 1 && NR > 1 && lines > limit { print prev "\t" lines }
  FNR == 1 { prev = FILENAME }
  { lines = FNR }
  END { if (NR > 0 && lines > limit) print prev "\t" lines }
')

while IFS=$'\t' read -r path lines; do
  [ -n "${path:-}" ] || continue
  echo "::error file=$path::$path has $lines lines, the limit is $LIMIT"
  FAILED=1
done <<< "$OVERSIZED"

if [ "$FAILED" -eq 0 ]; then
  count=$(printf '%s\n' "$FILES" | wc -l | tr -d ' ')
  echo "All $count files are within the $LIMIT line limit"
fi

exit $FAILED
