// Package alien provides graceful handling of unknown or unregistered types.
package alien

import (
	"fmt"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

// AlienComponent represents a component of an alien store.
type AlienComponent interface {
	String() string
	Accept(visitor store.Visitor)
}

// AlienPiece represents raw binary data from an unrecognized part.
type AlienPiece struct {
	data []byte
}

// NewAlienPiece creates a new AlienPiece with the given data.
func NewAlienPiece(data []byte) *AlienPiece {
	return &AlienPiece{
		data: data,
	}
}

// String returns a debug representation.
func (ap *AlienPiece) String() string {
	return fmt.Sprintf("AlienPiece{%d bytes}", len(ap.data))
}

// Accept implements the visitor pattern (no-op for alien pieces).
func (ap *AlienPiece) Accept(visitor store.Visitor) {
	// No-op: alien pieces are not processed
}

// GetData returns the raw binary data.
func (ap *AlienPiece) GetData() []byte {
	return ap.data
}

// AlienPart represents a recognized sub-store within an alien.
type AlienPart struct {
	store store.Store
}

// NewAlienPart creates a new AlienPart wrapping a store.
func NewAlienPart(s store.Store) *AlienPart {
	return &AlienPart{
		store: s,
	}
}

// String returns a debug representation.
func (ap *AlienPart) String() string {
	if ap.store != nil {
		return fmt.Sprintf("AlienPart{%s}", ap.store.String())
	}
	return "AlienPart{nil}"
}

// Accept implements the visitor pattern by delegating to the store.
func (ap *AlienPart) Accept(visitor store.Visitor) {
	if ap.store != nil {
		ap.store.Accept(visitor)
	}
}

// GetStore returns the wrapped store.
func (ap *AlienPart) GetStore() store.Store {
	return ap.store
}

// Alien represents an unregistered or incompatible type.
// It allows reading files even when they contain unknown types.
type Alien struct {
	store.BaseStore
	path  store.TypePath
	comps []AlienComponent
}

// NewAlien creates a new Alien with the given ID and type path.
func NewAlien(id oberon.Integer, path store.TypePath) *Alien {
	return &Alien{
		BaseStore: store.NewBaseStore(id),
		path:      path,
		comps:     make([]AlienComponent, 0),
	}
}

// GetTypeName returns the type name (from the path).
func (a *Alien) GetTypeName() string {
	if len(a.path) > 0 {
		return a.path[len(a.path)-1]
	}
	return "Alien"
}

// GetTypePath returns the original type path.
func (a *Alien) GetTypePath() store.TypePath {
	return a.path
}

// String returns a debug representation.
func (a *Alien) String() string {
	return fmt.Sprintf("Alien{id: %d, path: %s, components: %d}",
		a.GetID(), a.path.String(), len(a.comps))
}

// Accept implements the visitor pattern for Alien.
func (a *Alien) Accept(visitor store.Visitor) {
	// Check if this alien has already been visited (prevents cycles from LINK/NEWLINK)
	if !visitor.ShouldVisit(a) {
		return
	}

	for _, comp := range a.comps {
		comp.Accept(visitor)
	}
}

// Internalize is a no-op for Alien - aliens are internalized via internalizeAlien in the reader.
func (a *Alien) Internalize(reader store.Reader) error {
	// Aliens don't read their content via the normal internalize path
	// The reader's internalizeAlien method handles alien content
	return nil
}

// AddComponent adds a component to the alien.
func (a *Alien) AddComponent(comp AlienComponent) {
	a.comps = append(a.comps, comp)
}

// GetComponents returns all components.
func (a *Alien) GetComponents() []AlienComponent {
	return a.comps
}
