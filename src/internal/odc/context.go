package odc

import (
	"fmt"
	"strings"
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
