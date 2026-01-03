#!/usr/bin/bash
# More detailed test with hex inspection

echo "File size:"
ls -l ../_tests/mini1.odc

echo -e "\nFirst 400 bytes (hex):"
xxd -l 400 ../_tests/mini1.odc | head -20

echo -e "\nTrying to read with Go version:"
./odcread ../_tests/mini1.odc 2>&1 | head -5

echo -e "\nC++ version output:"
../src-cpp/odcread ../_tests/mini1.odc 2>&1 | head -5
