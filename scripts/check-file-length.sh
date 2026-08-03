#!/usr/bin/env bash
set -euo pipefail

LIMIT=250
FAILED=0

FILES=$(git ls-files \
  'src/**/*.go' 'src/**/*.ts' 'src/**/*.tsx' 'src/**/*.css' 'src/**/*.json' \
  'tests/**/*.go' 'tests/**/*.ts' \
  | grep -v -E '(_gen\.go|\.pb\.go|mock_|/mocks/|package-lock\.json|go\.sum)' || true)

for f in $FILES; do
  [ -f "$f" ] || continue
  lines=$(wc -l < "$f" | tr -d ' ')
  if [ "$lines" -gt "$LIMIT" ]; then
    echo "::error file=$f::$f has $lines lines, the limit is $LIMIT"
    FAILED=1
  fi
done

if [ "$FAILED" -eq 0 ]; then
  echo "All files are within the $LIMIT line limit"
fi

exit $FAILED
