# Technical Details

This document provides in-depth technical information about the `odcread` Go implementation.

## Project Structure

```
src/
├── cmd/
│   └── odcread/          # Main application
│       └── main.go       # Text extraction visitor
├── pkg/
│   ├── oberon/           # Primitive type definitions
│   ├── reader/           # Binary file reader
│   ├── store/            # Core data model
│   ├── textmodel/        # Text document components
│   ├── fold/             # Collapsible fold views
│   ├── alien/            # Unknown type handling
│   ├── typeregister/     # Runtime type registry
│   ├── visitor/          # Visitor pattern interface
│   └── encoding/         # Character encoding conversion
├── docs/                 # Documentation
└── tests/                # Tests
```

## Architecture

The implementation follows a modular architecture:

1.  **Reader**: Handles binary `.odc` format parsing with state management.
2.  **TypeRegister**: Manages runtime type registration and instantiation.
3.  **Store Hierarchy**: Represents document objects using Go's composition and interfaces.
4.  **Visitor Pattern**: Traverses the document tree for processing (e.g., text extraction).
5.  **Encoding**: Converts character encodings (ISO-8859-1, UCS-2) to UTF-8.

## Features & Implementation

- **Binary Parsing**: Comprehensive support for the `.odc` format, including nested stores.
- **Text Models**: Handles `ShortPiece`, `LongPiece`, and `ViewPiece` components.
- **Fold Support**: Correctly handles collapsible sections (folds) within documents.
- **Alien Types**: Graceful handling of unknown or unsupported types to prevent parsing failures.
- **Position Tracking**: Built-in position tracking for debugging binary format issues.

## Implementation Highlights

- **`pkg/reader/reader.go`**: Binary parsing with correct NIL store handling.
- **`pkg/textmodel/stdtextmodel.go`**: Attribute loop internalization.
- **`pkg/textmodel/pieces.go`**: Optimized `Read()` methods for all piece types.
- **`pkg/fold/fold.go`**: Fold internalization logic.

## Debugging

To debug binary format issues, you can use hex dump tools:

```bash
# View hex dump
xxd _tests/mini1.odc | head -50

# Check specific offsets
xxd -s 0x2B0 -l 100 _tests/mini4.odc
```

## Test Results Detail

- **mini1.odc**: Simple text document (exact match with C++ output).
- **mini2.odc**: Menu definitions (exact match with C++ output).
- **mini4.odc**: Complex document with folds (successfully parsed).
