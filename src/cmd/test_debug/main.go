package main

import (
	"fmt"
	"odcread/pkg/alien"
	"odcread/pkg/reader"
	"os"
)

func main() {
	file, _ := os.Open("../../_tests/mini1.odc")
	defer file.Close()

	r := reader.NewReader(file)

	// Read document header
	tag, _ := r.ReadInt()
	fmt.Fprintf(os.Stderr, "Tag: 0x%X\n", tag)

	version, _ := r.ReadInt()
	fmt.Fprintf(os.Stderr, "Version: %d\n", version)

	// Read position before store
	pos, _ := file.Seek(0, 1)
	fmt.Fprintf(os.Stderr, "Position before ReadStore: %d (0x%X)\n", pos, pos)

	// Try to read the root store
	s, err := r.ReadStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, "Successfully read: %s (ID: %d)\n", s.GetTypeName(), s.GetID())
	fmt.Fprintf(os.Stderr, "Type path: %v\n", s.GetTypePath())
	fmt.Fprintf(os.Stderr, "Actual type: %T\n", s)

	// Check if it's an Alien with direct type assertion
	if a, ok := s.(*alien.Alien); ok {
		comps := a.GetComponents()
		fmt.Fprintf(os.Stderr, "✓ IS ALIEN with %d components:\n", len(comps))
		for i, comp := range comps {
			fmt.Fprintf(os.Stderr, "  [%d] %s\n", i, comp.String())
		}
	} else {
		fmt.Fprintf(os.Stderr, "✗ NOT AN ALIEN - unexpected!\n")
	}
}
