// Package store - Store hierarchy implementations
package store

import "fmt"

const (
	// TypeNames for the store hierarchy
	TypeNameStore          = "Stores.Store^"
	TypeNameElem           = "Stores.Elem^"
	TypeNameModel          = "Models.Model^"
	TypeNameContainerModel = "Containers.Model^"
)

// StoreImpl is a concrete implementation of the base Store.
type StoreImpl struct {
	BaseStore
}

// NewStore creates a new Store instance.
func NewStore(id int32) *StoreImpl {
	return &StoreImpl{
		BaseStore: NewBaseStore(id),
	}
}

// GetTypeName returns the type name for Store.
func (s *StoreImpl) GetTypeName() string {
	return TypeNameStore
}

// Elem represents an "Elem" store - a legacy BlackBox type.
type Elem struct {
	BaseStore
}

// NewElem creates a new Elem instance.
func NewElem(id int32) *Elem {
	return &Elem{
		BaseStore: NewBaseStore(id),
	}
}

// GetTypeName returns the type name for Elem.
func (e *Elem) GetTypeName() string {
	return TypeNameElem
}

// Internalize reads Elem data from the reader.
func (e *Elem) Internalize(reader Reader) error {
	// Call base internalization
	if err := e.BaseStore.Internalize(reader); err != nil {
		return err
	}
	// Elem version check (0..0)
	_, err := reader.ReadVersion(0, 0)
	return err
}

// String returns a string representation of the Elem.
func (e *Elem) String() string {
	return fmt.Sprintf("Elem{id: %d}", e.id)
}

// Model represents a Model store - the basis for all model objects (MVC framework).
type Model struct {
	Elem
}

// NewModel creates a new Model instance.
func NewModel(id int32) *Model {
	return &Model{
		Elem: *NewElem(id),
	}
}

// GetTypeName returns the type name for Model.
func (m *Model) GetTypeName() string {
	return TypeNameModel
}

// Internalize reads Model data from the reader.
func (m *Model) Internalize(reader Reader) error {
	// Call parent internalization
	if err := m.Elem.Internalize(reader); err != nil {
		return err
	}
	// Model version check (0..0)
	_, err := reader.ReadVersion(0, 0)
	return err
}

// String returns a string representation of the Model.
func (m *Model) String() string {
	return fmt.Sprintf("Model{id: %d}", m.id)
}

// ContainerModel is the supertype for models that contain other stuff (e.g., TextModel).
type ContainerModel struct {
	Model
}

// NewContainerModel creates a new ContainerModel instance.
func NewContainerModel(id int32) *ContainerModel {
	return &ContainerModel{
		Model: *NewModel(id),
	}
}

// GetTypeName returns the type name for ContainerModel.
func (cm *ContainerModel) GetTypeName() string {
	return TypeNameContainerModel
}

// Internalize reads ContainerModel data from the reader.
func (cm *ContainerModel) Internalize(reader Reader) error {
	// Call parent internalization
	if err := cm.Model.Internalize(reader); err != nil {
		return err
	}
	// ContainerModel version check (0..0)
	_, err := reader.ReadVersion(0, 0)
	return err
}

// String returns a string representation of the ContainerModel.
func (cm *ContainerModel) String() string {
	return fmt.Sprintf("ContainerModel{id: %d}", cm.id)
}
