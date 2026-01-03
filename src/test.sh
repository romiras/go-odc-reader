#!/usr/bin/bash
# Test script for odcread Go implementation

echo "=== Testing odcread Go implementation ==="
echo ""

GOBIN="./odcread"
CPPBIN="../src-cpp/odcread"

if [ ! -f "$GOBIN" ]; then
    echo "Error: Go binary not found at $GOBIN"
    exit 1
fi

if [ ! -f "$CPPBIN" ]; then
    echo "Warning: C++ binary not found at $CPPBIN"
    CPPBIN=""
fi

# Find test files
TEST_FILES=(../_tests/mini*.odc)

for test_file in "${TEST_FILES[@]}"; do
    if [ ! -f "$test_file" ]; then
        continue
    fi

    echo "Testing: $(basename "$test_file")"
    echo "----------------------------------------"

    # Test Go version
    echo "Go output:"
    $GOBIN "$test_file" 2>&1 | head -10
    go_exit=$?
    echo "Exit code: $go_exit"
    echo ""

    # Test C++ version if available
    if [ -n "$CPPBIN" ]; then
        echo "C++ output:"
        $CPPBIN "$test_file" 2>&1 | head -10
        cpp_exit=$?
        echo "Exit code: $cpp_exit"
        echo ""
    fi

    echo "========================================"
    echo ""
done
