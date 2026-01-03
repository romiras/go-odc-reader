#!/usr/bin/env bash
set -euo pipefail

# Mass-check .odc files
# Usage: ./check.sh <binary_path> <test_dir> [paths_file]

BINARY="$1"
TEST_DIR="$2"
PATHS_FILE="${3:-paths.txt}"
FAILED_LOG=".failed_tests"

rm -f "$FAILED_LOG"
FILES_LIST=$(mktemp)

# Collect files
if [ -f "$PATHS_FILE" ]; then
    while read -r p; do
        # Skip comments and empty lines
        [[ "$p" =~ ^#.*$ ]] || [ -z "$p" ] && continue
        if [ -d "$p" ]; then
            find "$p" -name "*.odc" >> "$FILES_LIST"
        else
            echo "⚠️  Skipping missing directory: $p"
        fi
    done < "$PATHS_FILE"
else
    if [ -d "$TEST_DIR" ]; then
        find "$TEST_DIR" -name "*.odc" >> "$FILES_LIST"
    else
        echo "❌ Error: Test directory '$TEST_DIR' not found."
        rm -f "$FILES_LIST"
        exit 1
    fi
fi

# Run checks
failed_count=0
total_count=0
ERR_TEMP=$(mktemp)

while read -r f; do
    [ -z "$f" ] && continue
    total_count=$((total_count + 1))
    if ! "$BINARY" "$f" > /dev/null 2> "$ERR_TEMP"; then
        echo "❌ FAIL: $f"
        sed 's/^/  /' "$ERR_TEMP"
        echo "$f" >> "$FAILED_LOG"
        failed_count=$((failed_count + 1))
    fi
done < "$FILES_LIST"

rm -f "$FILES_LIST" "$ERR_TEMP"

if [ "$failed_count" -eq 0 ]; then
    echo "✅ All $total_count files passed!"
else
    echo "⚠️  $failed_count/$total_count files failed. List saved to $FAILED_LOG"
fi
