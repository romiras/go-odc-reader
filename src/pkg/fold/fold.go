// Package fold provides collapsible fold views for documents.
package fold

import (
	"fmt"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

const (
	TypeNameView = "Views.View^"
	TypeNameFold = "StdFolds.Fold^"
)

// View is the supertype for views in the MVC framework.
type View struct {
	store.BaseStore
}

// NewView creates a new View instance.
func NewView(id oberon.Integer) *View {
	return &View{
		BaseStore: store.NewBaseStore(id),
	}
}

// GetTypeName returns the type name for View.
func (v *View) GetTypeName() string {
	return TypeNameView
}

// Internalize reads View data from the reader.
func (v *View) Internalize(reader store.Reader) error {
	if err := v.BaseStore.Internalize(reader); err != nil {
		return err
	}
	_, err := reader.ReadVersion(0, 0)
	return err
}

// String returns a string representation of the View.
func (v *View) String() string {
	return fmt.Sprintf("View{id: %d}", v.GetID())
}

// Fold represents a collapsible section in a document.
type Fold struct {
	View
	hidden    store.Store
	label     []oberon.ShortChar
	collapsed bool
}

// NewFold creates a new Fold instance.
func NewFold(id oberon.Integer) *Fold {
	return &Fold{
		View: *NewView(id),
	}
}

// GetTypeName returns the type name for Fold.
func (f *Fold) GetTypeName() string {
	return TypeNameFold
}

// Internalize reads Fold data from the reader.
func (f *Fold) Internalize(reader store.Reader) error {
	// Call parent internalization
	if err := f.View.Internalize(reader); err != nil {
		return err
	}

	// Read version (0..0 supported, matching C++ line 33)
	_, err := reader.ReadVersion(0, 0)
	if err != nil {
		return err
	}

	// Read leftSide (C++ line 36 - we don't use it but must consume it)
	_, err = reader.ReadSInt()
	if err != nil {
		return fmt.Errorf("failed to read leftSide: %w", err)
	}

	// Read collapsed state as SInt (C++ line 38)
	collapsedInt, err := reader.ReadSInt()
	if err != nil {
		return fmt.Errorf("failed to read collapsed state: %w", err)
	}
	f.collapsed = collapsedInt == 0

	// Read label as null-terminated string (C++ line 42)
	labelStr, err := reader.ReadSString()
	if err != nil {
		return fmt.Errorf("failed to read label: %w", err)
	}
	f.label = []oberon.ShortChar(labelStr)

	// Read the hidden part (a Store)
	hidden, err := reader.ReadStore()
	if err != nil {
		return fmt.Errorf("failed to read hidden store: %w", err)
	}
	f.hidden = hidden

	return nil
}

// String returns a debug representation of the Fold.
func (f *Fold) String() string {
	labelStr := string(f.label)
	return fmt.Sprintf("Fold{id: %d, collapsed: %v, label: %q}", f.GetID(), f.collapsed, labelStr)
}

// Accept implements the visitor pattern for Fold.
func (f *Fold) Accept(visitor store.Visitor) {
	visitor.FoldLeft(f.collapsed)

	// Visit the hidden part
	if f.hidden != nil {
		f.hidden.Accept(visitor)
	}

	visitor.FoldRight()
}

// IsCollapsed returns whether the fold is collapsed.
func (f *Fold) IsCollapsed() bool {
	return f.collapsed
}

// GetLabel returns the fold label.
func (f *Fold) GetLabel() string {
	return string(f.label)
}

// GetHidden returns the hidden store.
func (f *Fold) GetHidden() store.Store {
	return f.hidden
}
