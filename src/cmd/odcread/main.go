// Package main - odcread command-line application
package main

import (
	"fmt"
	"os"
	"strings"

	"odcread/pkg/encoding"
	"odcread/pkg/oberon"
	"odcread/pkg/reader"
	"odcread/pkg/store"
	"odcread/pkg/textmodel"
	_ "odcread/pkg/typeregister" // Import for side-effect (type registration)
)

const (
	docTag     = oberon.Integer(0x6F4F4443)
	docVersion = oberon.Integer(0)
)

// Context interface for text accumulation
type Context interface {
	AddPiece(piece string)
	GetPlainText() string
}

// PartContext - simple text accumulation
type PartContext struct {
	text strings.Builder
}

func (pc *PartContext) AddPiece(piece string) {
	pc.text.WriteString(piece)
}

func (pc *PartContext) GetPlainText() string {
	return pc.text.String()
}

// FoldContext - handles collapsed/expanded folds
type FoldContext struct {
	collapsed bool
	haveFirst bool
	firstPart strings.Builder
	remainder strings.Builder
}

func NewFoldContext(collapsed bool) *FoldContext {
	return &FoldContext{collapsed: collapsed}
}

func (fc *FoldContext) AddPiece(piece string) {
	if !fc.haveFirst {
		fc.haveFirst = true
		fc.firstPart.WriteString(piece)
	} else {
		fc.remainder.WriteString(piece)
	}
}

func (fc *FoldContext) GetPlainText() string {
	if fc.collapsed {
		return fmt.Sprintf("##=>%s\n%s##<=", fc.remainder.String(), fc.firstPart.String())
	}
	return fmt.Sprintf("##=>%s\n%s##<=", fc.firstPart.String(), fc.remainder.String())
}

// MyVisitor - concrete visitor implementation for text extraction
type MyVisitor struct {
	contextStack []Context
}

func NewMyVisitor() *MyVisitor {
	return &MyVisitor{
		contextStack: make([]Context, 0),
	}
}

func (mv *MyVisitor) PartStart() {
	mv.contextStack = append(mv.contextStack, &PartContext{})
}

func (mv *MyVisitor) PartEnd() {
	mv.terminateContext()
}

func (mv *MyVisitor) FoldLeft(collapsed bool) {
	mv.contextStack = append(mv.contextStack, NewFoldContext(collapsed))
}

func (mv *MyVisitor) FoldRight() {
	mv.terminateContext()
}

func (mv *MyVisitor) terminateContext() {
	if len(mv.contextStack) == 0 {
		return
	}

	top := len(mv.contextStack) - 1
	ctx := mv.contextStack[top]
	mv.contextStack = mv.contextStack[:top]

	if len(mv.contextStack) == 0 {
		// Top-level context - print to stdout
		fmt.Println(ctx.GetPlainText())
	} else {
		// Nested context - add to parent
		mv.contextStack[len(mv.contextStack)-1].AddPiece(ctx.GetPlainText())
	}
}

func (mv *MyVisitor) TextShortPiece(piece interface{}) {
	if sp, ok := piece.(*textmodel.ShortPiece); ok {
		str, err := encoding.ConvertLatin1(sp.GetBuffer())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to convert short piece: %v\n", err)
			return
		}
		if len(mv.contextStack) > 0 {
			mv.contextStack[len(mv.contextStack)-1].AddPiece(str)
		}
	}
}

func (mv *MyVisitor) TextLongPiece(piece interface{}) {
	if lp, ok := piece.(*textmodel.LongPiece); ok {
		str, err := encoding.ConvertUCS2(lp.GetBuffer())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to convert long piece: %v\n", err)
			return
		}
		if len(mv.contextStack) > 0 {
			mv.contextStack[len(mv.contextStack)-1].AddPiece(str)
		}
	}
}

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
	visitor := NewMyVisitor()
	s.Accept(visitor)
}
