// Package store provides the core data model for persistent, extensible objects.
package store

import (
	"fmt"
	"strings"

	"odcread/pkg/oberon"
)

// Store type markers (used in binary format)
const (
	NEWBASE oberon.ShortChar = 0xF0 // new base type (level = 0), not yet in dict
	NEWEXT  oberon.ShortChar = 0xF1 // new extension type (level = 1), not yet in dict
	OLDTYPE oberon.ShortChar = 0xF2 // old type, already in dict
	NIL     oberon.ShortChar = 0x80 // nil store
	LINK    oberon.ShortChar = 0x81 // link to another elem in same file
	STORE   oberon.ShortChar = 0x82 // general store
	ELEM    oberon.ShortChar = 0x83 // elem store
	NEWLINK oberon.ShortChar = 0x84 // link to another non-elem store in same file
)

// TypePath represents the inheritance path of a type.
type TypePath []string

// String returns a string representation of the type path.
func (tp TypePath) String() string {
	if len(tp) == 0 {
		return "<empty path>"
	}
	return strings.Join(tp, " -> ")
}

// Store is the interface for all storable, extensible data types.
// Stores are used as base types for all objects that must be both extensible and persistent.
type Store interface {
	// GetID returns the unique identifier for this store.
	GetID() oberon.Integer

	// GetTypeName returns the full type name (including module).
	GetTypeName() string

	// GetTypePath returns the full inheritance path for this type.
	GetTypePath() TypePath

	// Internalize reads the store's contents from the reader.
	Internalize(reader Reader) error

	// Accept implements the Visitor pattern for traversing the store tree.
	Accept(visitor Visitor)

	// String returns a debug representation of the store.
	String() string
}

// Reader interface defines methods needed to read stores from binary format.
// This is a forward declaration - the actual implementation is in the reader package.
type Reader interface {
	ReadVersion(min, max oberon.Integer) (oberon.Integer, error)
	ReadStore() (Store, error)
	ReadInt() (oberon.Integer, error)
	ReadSInt() (oberon.ShortInt, error)
	ReadByte() (oberon.Byte, error)
	ReadSChar() (oberon.ShortChar, error)
	ReadLChar() (oberon.Char, error)
	ReadSString() (string, error)
	IsCancelled() bool
	TurnIntoAlien(cause int) error
}

// Visitor interface for the visitor pattern.
// This is a forward declaration - the actual implementation is in the visitor package.
type Visitor interface {
	PartStart()
	PartEnd()
	FoldLeft(collapsed bool)
	FoldRight()
}

// BaseStore provides common functionality for all Store implementations.
type BaseStore struct {
	id oberon.Integer
}

// NewBaseStore creates a new BaseStore with the given ID.
func NewBaseStore(id oberon.Integer) BaseStore {
	return BaseStore{id: id}
}

// GetID returns the store's unique identifier.
func (bs *BaseStore) GetID() oberon.Integer {
	return bs.id
}

// GetTypePath calculates the full type path for this store.
func (bs *BaseStore) GetTypePath() TypePath {
	// This will be implemented by calling calcTypePath recursively
	// For now, return a simple path with just the type name
	return TypePath{}
}

// String returns a basic string representation.
func (bs *BaseStore) String() string {
	return fmt.Sprintf("Store{id: %d}", bs.id)
}

// Internalize reads the base store data (just validates version 0).
func (bs *BaseStore) Internalize(reader Reader) error {
	_, err := reader.ReadVersion(0, 0)
	return err
}

// Accept is a default implementation that does nothing.
func (bs *BaseStore) Accept(visitor Visitor) {
	// Default: no-op
}
