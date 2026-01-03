#!/usr/bin/env bash
set -euo pipefail

# Re-check failed .odc files
# Usage: ./check-failed.sh <binary_path>

BINARY="$1"
FAILED_LOG=".failed_tests"

if [ ! -s "$FAILED_LOG" ]; then
    echo "No failed tests to run."
    exit 0
fi

echo "Re-checking failed files..."
mv "$FAILED_LOG" "${FAILED_LOG}.old"
failed_count=0
total_count=0
ERR_TEMP=$(mktemp)

while read -r f; do
    [ -z "$f" ] && continue
    total_count=$((total_count + 1))
    if ! "$BINARY" "$f" > /dev/null 2> "$ERR_TEMP"; then
        echo "❌ STILL FAILING: $f"
        sed 's/^/  /' "$ERR_TEMP"
        echo "$f" >> "$FAILED_LOG"
        failed_count=$((failed_count + 1))
    else
        echo "✅ FIXED: $f"
    fi
done < "${FAILED_LOG}.old"

rm -f "${FAILED_LOG}.old" "$ERR_TEMP"
echo "Summary: $((total_count - failed_count))/$total_count fixed."
