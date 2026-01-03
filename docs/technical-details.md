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

## Binary Format Insights

During development, several critical details about the BlackBox component format were discovered and implemented:

### StdTextModel Structure
The `TextModels.StdModel` format includes a `metaLen` field that caused significant confusion during porting.
- **Structure**: `[version (1b)] [metaLen (4b)] [metadata block...] [pieces block...]`
- **Critical Finding**: The `metaLen` field indicates the length of the *metadata section* (which contains the attribute dictionary and piece headers). It is **not** a separate block to be consumed or skipped. The parser naturally consumes these bytes while reading the piece headers.
- **Resolution**: The reader reads this value but relies on the natural flow of reading piece headers (until `ano == -1`) to consume the data.

### Recursive Store Parsing & Alien State
The format allows for "Alien" stores (unknown types) that may contain other embedded stores.
- **Challenge**: Deeply nested alien stores (e.g., >3 levels) caused position mismatches because the reader state's `End` position was not being correctly preserved across recursive calls.
- **Solution**: The reader now explicitly captures the `storeEnd` position *before* creating a new reader state for a nested store. This ensures that even if the new state is initialized empty, the bound checking uses the correct absolute file position.

## Features & Implementation

- **Binary Parsing**: Comprehensive support for the `.odc` format, including nested stores and NIL stores.
- **Text Models**: Handles `ShortPiece`, `LongPiece`, and `ViewPiece` components effectively.
- **Fold Support**: Correctly handles collapsible sections (folds) within documents.
- **Alien Types**: Robust handling of unknown or unsupported types (Alien stores) to prevent parsing failures, even when nested.
- **Position Tracking**: Strict position tracking to validate parsing integrity.

## Implementation Highlights

- **`pkg/reader/reader.go`**: 
  - Implements the core state machine for parsing.
  - Handles the 8-byte header for NIL stores (a subtle format consistency).
  - Manages recursive state for nested Alien stores.
- **`pkg/textmodel/stdtextmodel.go`**: 
  - Implements the corrected `StdTextModel` parsing logic.
  - Handles the attribute dictionary loop used for piece compression.
- **`pkg/fold/fold.go`**: 
  - Implements `Fold.Internalize` for collapsible document sections.

## Known Limitations

- **Link Stores**: The format specification includes `LinkStore` (markers `0x34`, `0x35`) and `NewLinkStore`. Neither the original C++ implementation nor this Go port currently supports these types.
  - **Impact**: Approximately 0.25% of real-world .odc files (e.g., `Sys-Map.odc`) fail to parse due to this limitation. This is consistent with the reference implementation.

## Debugging

To debug binary format issues, use the provided Makefile commands or hex dump tools:

```bash
# Mass check all files
make check

# Re-run only failed files
make check-failed

# View hex dump for specific offsets
xxd _tests/mini1.odc | head -50

# Check specific offsets
xxd -s 0x2B0 -l 100 _tests/mini4.odc
```

## Test Results Detail

- **mini1.odc**: Simple text document (✅ PASS)
- **mini2.odc**: Menu definitions (✅ PASS)
- **mini4.odc**: Complex document with folds (✅ PASS)
- **Bulk Test**: 398/399 real-world files parse successfully (99.75% success rate). The single failure matches the C++ reference implementation's limitation.
