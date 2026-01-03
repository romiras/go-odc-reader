// Package textmodel - StdTextModel implementation
package textmodel

import (
	"fmt"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

// StdTextModel is the standard implementation of a TextModel.
// It consists of a series of TextPieces.
type StdTextModel struct {
	TextModel
	pieces []TextPiece
}

// NewStdTextModel creates a new StdTextModel instance.
func NewStdTextModel(id oberon.Integer) *StdTextModel {
	return &StdTextModel{
		TextModel: *NewTextModel(id),
		pieces:    make([]TextPiece, 0),
	}
}

// GetTypeName returns the type name for StdTextModel.
func (stm *StdTextModel) GetTypeName() string {
	return TypeNameStdTextModel
}

// Internalize reads StdTextModel data from the reader.
// Format: version, metadata_length, then pieces with attributes until ano==-1
func (stm *StdTextModel) Internalize(reader store.Reader) error {
	// Call parent internalization
	if err := stm.TextModel.Internalize(reader); err != nil {
		return err
	}

	// Read version (0..1 supported)
	version, err := reader.ReadVersion(0, 1)
	if err != nil {
		return err
	}

	// Read metadata section length (not used, but must be read)
	_, err = reader.ReadInt()
	if err != nil {
		return fmt.Errorf("failed to read metadata length: %w", err)
	}

	// Attribute dictionary for pieces
	dict := make([]store.Store, 0)

	// Read pieces in a loop until ano == -1
	stm.pieces = make([]TextPiece, 0)

	ano, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("failed to read first ano: %w", err)
	}

	for ano != -1 {
		// Read or reuse attribute from dictionary
		if int(ano) == len(dict) {
			// New attribute - read it
			attr, err := reader.ReadStore()
			if err != nil {
				return fmt.Errorf("failed to read attribute store: %w", err)
			}
			dict = append(dict, attr)
		}
		// Note: We don't use the attribute for now, but it's stored for future use

		// Read piece length
		pieceLen, err := reader.ReadInt()
		if err != nil {
			return fmt.Errorf("failed to read piece length: %w", err)
		}

		var piece TextPiece
		if pieceLen > 0 {
			// ShortPiece (8-bit characters)
			piece = NewShortPiece(uint(pieceLen))
		} else if pieceLen < 0 {
			// LongPiece (16-bit characters)
			// pieceLen is negative and in bytes, so divide by 2 for character count
			piece = NewLongPiece(uint(-pieceLen / 2))
		} else {
			// ViewPiece (embedded view, pieceLen == 0)
			// Read view width and height (ignored)
			_, err := reader.ReadInt()
			if err != nil {
				return fmt.Errorf("failed to read view width: %w", err)
			}
			_, err = reader.ReadInt()
			if err != nil {
				return fmt.Errorf("failed to read view height: %w", err)
			}

			// Read the embedded view
			view, err := reader.ReadStore()
			if err != nil {
				return fmt.Errorf("failed to read embedded view: %w", err)
			}
			piece = NewViewPiece(view)
		}

		stm.pieces = append(stm.pieces, piece)

		// Read next ano
		ano, err = reader.ReadByte()
		if err != nil {
			return fmt.Errorf("failed to read next ano: %w", err)
		}
	}

	// Now read the actual piece content
	for i, piece := range stm.pieces {
		if err := piece.Read(reader); err != nil {
			return fmt.Errorf("failed to read piece %d content: %w", i, err)
		}
	}

	// Handle version 1 specific data if needed
	if version >= 1 {
		// Version 1 may have additional data - for now we skip it
		// In the C++ version, this might read additional formatting info
	}

	return nil
}

// String returns a debug representation of the StdTextModel.
func (stm *StdTextModel) String() string {
	return fmt.Sprintf("StdTextModel{id: %d, pieces: %d}", stm.GetID(), len(stm.pieces))
}

// Accept implements the visitor pattern for StdTextModel.
func (stm *StdTextModel) Accept(visitor store.Visitor) {
	visitor.PartStart()
	for _, piece := range stm.pieces {
		piece.Accept(visitor)
	}
	visitor.PartEnd()
}

// GetPieces returns the text pieces (for testing/debugging).
func (stm *StdTextModel) GetPieces() []TextPiece {
	return stm.pieces
}
