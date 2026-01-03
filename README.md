# odcread (Go Port)

A lightweight Go utility for extracting plain text from Oberon/F compound documents (`.odc` files).

## Why use odcread?

Developed to provide a modern, portable alternative for the Oberon community.

Oberon documents (`.odc`) are a binary compound format used by the BlackBox Component Builder, WinBUGS, and OpenBUGS. While powerful, they are difficult to read without specialized software.

### Key Use Cases:

1.  **Stand-alone Extraction**: Extract text from `.odc` files on any platform supported by Go, without needing the BlackBox Component Builder environment.
2.  **Git Companion**: Use it as a `textconv` filter for Git. This allows you to `git diff` `.odc` files as plain text, making it easy to track changes in documents or source code stored in the Oberon/F format.

## Quick Start

### Installation

Ensure you have [Go](https://go.dev/doc/install) installed (version 1.20 or later).

```bash
# Clone the repository
git clone https://github.com/romiras/go-odc-reader
cd go-odc-reader

# Build the utility
make build
```

### Basic Usage

Simply pass an `.odc` file to the executable. The text content will be sent to standard output (stdout).

```bash
./bin/odcread document.odc > output.txt
```

### Using as a Git Diff tool

To see text changes when you modify `.odc` files in a Git repository:

1.  Add the following to your `~/.gitconfig` or the project's `.git/config`:
    ```ini
    [diff "odc"]
        textconv = /path/to/bin/odcread
    ```
2.  Create or update a `.gitattributes` file in your repository:
    ```
    *.odc diff=odc
    ```

## Project Status

✅ **100% Functional**: All core document types and pieces are supported, including folds and various character encodings.

## Credits

This project was inspired by the original `odcread` C++ utility developed by [Gert van Valkenhoef](https://github.com/gertvv) ([gertvv/odcread](https://github.com/gertvv/odcread)).

---

For more detailed information, see:
- [Technical Details](docs/technical-details.md)
