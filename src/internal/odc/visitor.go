package odc

import (
	"fmt"
	"os"

	"odcread/pkg/encoding"
	"odcread/pkg/store"
	"odcread/pkg/textmodel"
)

// MyVisitor - concrete visitor implementation for text extraction
type MyVisitor struct {
	contextStack []Context
	visited      map[store.Store]bool // Track visited stores by pointer to prevent cycles
}

func NewMyVisitor() *MyVisitor {
	return &MyVisitor{
		contextStack: make([]Context, 0),
		visited:      make(map[store.Store]bool),
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

func (mv *MyVisitor) ShouldVisit(s store.Store) bool {
	if mv.visited[s] {
		return false
	}
	mv.visited[s] = true
	return true
}
