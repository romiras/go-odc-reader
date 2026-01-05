// Package main - odcread command-line application
package main

import (
	"fmt"
	"os"

	"odcread/internal/odc"
	"odcread/pkg/oberon"
	"odcread/pkg/reader"
	"odcread/pkg/store"
	_ "odcread/pkg/typeregister" // Import for side-effect (type registration)
)

const (
	docTag     = oberon.Integer(0x6F4F4443)
	docVersion = oberon.Integer(0)
)

// importDocument reads and validates an .odc document.
func importDocument(file *os.File) (store.Store, error) {
	r := reader.NewReader(file)

	// Read and validate document tag
	tag, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read document tag: %w", err)
	}

	if tag != docTag {
		return nil, fmt.Errorf("invalid document tag: 0x%X (expected 0x%X)", tag, docTag)
	}

	// Read and validate document version
	version, err := r.ReadInt()
	if err != nil {
		return nil, fmt.Errorf("failed to read document version: %w", err)
	}

	if version != docVersion {
		return nil, fmt.Errorf("unsupported document version: %d (expected %d)", version, docVersion)
	}

	// Read the root store
	s, err := r.ReadStore()
	if err != nil {
		return nil, fmt.Errorf("failed to read root store: %w", err)
	}

	return s, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <file.odc>\n", os.Args[0])
		os.Exit(1)
	}

	// Open the input file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(2)
	}
	defer file.Close()

	// Import the document
	s, err := importDocument(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing document: %v\n", err)
		os.Exit(2)
	}

	if s == nil {
		fmt.Fprintf(os.Stderr, "Error: document root is nil\n")
		os.Exit(2)
	}

	// Process the document with the visitor
	visitor := odc.NewMyVisitor()
	s.Accept(visitor)
}
