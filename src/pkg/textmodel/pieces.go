// Package textmodel provides text document components and text pieces.
package textmodel

import (
	"fmt"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

const (
	TypeNameTextModel    = "TextModels.Model^"
	TypeNameStdTextModel = "TextModels.StdModel^"
)

// TextModel is the abstract base for text models.
type TextModel struct {
	store.ContainerModel
}

// NewTextModel creates a new TextModel instance.
func NewTextModel(id oberon.Integer) *TextModel {
	return &TextModel{
		ContainerModel: *store.NewContainerModel(id),
	}
}

// GetTypeName returns the type name for TextModel.
func (tm *TextModel) GetTypeName() string {
	return TypeNameTextModel
}

// Internalize reads TextModel data from the reader.
func (tm *TextModel) Internalize(reader store.Reader) error {
	if err := tm.ContainerModel.Internalize(reader); err != nil {
		return err
	}
	_, err := reader.ReadVersion(0, 0)
	return err
}

// TextPiece is the interface for all text piece types.
type TextPiece interface {
	// Read reads the piece content from the reader.
	Read(reader store.Reader) error

	// String returns a debug representation.
	String() string

	// Accept implements the visitor pattern.
	Accept(visitor store.Visitor)

	// Size returns the size in bytes (excluding null terminator).
	Size() uint
}

// basePiece provides common functionality for text pieces.
type basePiece struct {
	length uint
}

// Size returns the piece size in bytes.
func (bp *basePiece) Size() uint {
	return bp.length
}

// ShortPiece represents a text piece with 8-bit Latin-1 characters.
type ShortPiece struct {
	basePiece
	buffer []oberon.ShortChar
}

// NewShortPiece creates a new ShortPiece with the given length.
func NewShortPiece(length uint) *ShortPiece {
	return &ShortPiece{
		basePiece: basePiece{length: length},
		buffer:    make([]oberon.ShortChar, length+1), // +1 for null terminator
	}
}

// Read reads the short piece content from the reader.
func (sp *ShortPiece) Read(reader store.Reader) error {
	// Read exactly 'length' characters (not length+1)
	for i := 0; i < int(sp.length); i++ {
		ch, err := reader.ReadSChar()
		if err != nil {
			return fmt.Errorf("failed to read short char at position %d: %w", i, err)
		}
		sp.buffer[i] = ch
	}
	// Null-terminate
	sp.buffer[sp.length] = 0
	return nil
}

// GetBuffer returns the raw buffer contents.
func (sp *ShortPiece) GetBuffer() []oberon.ShortChar {
	return sp.buffer
}

// String returns a debug representation.
func (sp *ShortPiece) String() string {
	return fmt.Sprintf("ShortPiece{len: %d}", sp.length)
}

// Accept implements the visitor pattern.
func (sp *ShortPiece) Accept(visitor store.Visitor) {
	// Type-assert to the full visitor interface if available
	if v, ok := visitor.(interface{ TextShortPiece(interface{}) }); ok {
		v.TextShortPiece(sp)
	}
}

// LongPiece represents a text piece with 16-bit Unicode characters.
type LongPiece struct {
	basePiece
	buffer []oberon.Char
}

// NewLongPiece creates a new LongPiece with the given length.
func NewLongPiece(length uint) *LongPiece {
	return &LongPiece{
		basePiece: basePiece{length: length},
		buffer:    make([]oberon.Char, length+1), // +1 for null terminator
	}
}

// Read reads the long piece content from the reader.
func (lp *LongPiece) Read(reader store.Reader) error {
	// Read exactly 'length' characters (length is in chars, not bytes)
	// Note: length here is already adjusted (d_len/2 in C++)
	for i := 0; i < int(lp.length); i++ {
		ch, err := reader.ReadLChar()
		if err != nil {
			return fmt.Errorf("failed to read long char at position %d: %w", i, err)
		}
		lp.buffer[i] = ch
	}
	// Null-terminate
	lp.buffer[lp.length] = 0
	return nil
}

// GetBuffer returns the raw buffer contents.
func (lp *LongPiece) GetBuffer() []oberon.Char {
	return lp.buffer
}

// String returns a debug representation.
func (lp *LongPiece) String() string {
	return fmt.Sprintf("LongPiece{len: %d}", lp.length)
}

// Accept implements the visitor pattern.
func (lp *LongPiece) Accept(visitor store.Visitor) {
	// Type-assert to the full visitor interface if available
	if v, ok := visitor.(interface{ TextLongPiece(interface{}) }); ok {
		v.TextLongPiece(lp)
	}
}

// ViewPiece represents a text piece that embeds a View.
type ViewPiece struct {
	basePiece
	view store.Store
}

// NewViewPiece creates a new ViewPiece with the given view.
func NewViewPiece(view store.Store) *ViewPiece {
	return &ViewPiece{
		basePiece: basePiece{length: 1}, // View pieces have length 1
		view:      view,
	}
}

// Read reads the view piece (reads one extra byte as per C++ implementation).
func (vp *ViewPiece) Read(reader store.Reader) error {
	// ViewPiece requires reading one extra byte
	_, err := reader.ReadByte()
	return err
}

// String returns a debug representation.
func (vp *ViewPiece) String() string {
	if vp.view != nil {
		return fmt.Sprintf("ViewPiece{view: %s}", vp.view.GetTypeName())
	}
	return "ViewPiece{view: nil}"
}

// Accept implements the visitor pattern.
func (vp *ViewPiece) Accept(visitor store.Visitor) {
	if vp.view != nil {
		vp.view.Accept(visitor)
	}
}

// GetView returns the embedded view.
func (vp *ViewPiece) GetView() store.Store {
	return vp.view
}
