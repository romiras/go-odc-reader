// Package typeregister provides runtime type registration for dynamic Store instantiation.
package typeregister

import (
	"sync"

	"odcread/pkg/oberon"
	"odcread/pkg/store"
)

// TypeProxyBase is the interface for type proxies.
type TypeProxyBase interface {
	// GetName returns the full type name (including module).
	GetName() string

	// GetSuper returns the supertype name, or nil if this is a top-level type.
	GetSuper() *string

	// NewInstance creates a new instance of this type with the given ID.
	NewInstance(id oberon.Integer) store.Store
}

// TypeRegister is a singleton registry of Oberon/BlackBox types.
type TypeRegister struct {
	registry map[string]TypeProxyBase
	mu       sync.RWMutex
}

var (
	instance *TypeRegister
	once     sync.Once
)

// GetInstance returns the singleton TypeRegister instance.
func GetInstance() *TypeRegister {
	once.Do(func() {
		instance = &TypeRegister{
			registry: make(map[string]TypeProxyBase),
		}
	})
	return instance
}

// Add registers a new type proxy.
func (tr *TypeRegister) Add(name string, proxy TypeProxyBase) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.registry[name] = proxy
}

// Get retrieves a type proxy by name.
// Returns nil if the type is not registered.
func (tr *TypeRegister) Get(name string) TypeProxyBase {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	return tr.registry[name]
}

// Has checks if a type is registered.
func (tr *TypeRegister) Has(name string) bool {
	tr.mu.RLock()
	defer tr.mu.RUnlock()
	_, exists := tr.registry[name]
	return exists
}

// StoreProxy is a concrete implementation of TypeProxyBase.
type StoreProxy struct {
	name      string
	superName *string
	factory   func(oberon.Integer) store.Store
}

// NewStoreProxy creates a new StoreProxy with the given name and factory function.
func NewStoreProxy(name string, factory func(oberon.Integer) store.Store) *StoreProxy {
	return &StoreProxy{
		name:    name,
		factory: factory,
	}
}

// NewStoreProxyWithSuper creates a new StoreProxy with a supertype.
func NewStoreProxyWithSuper(name string, superName string, factory func(oberon.Integer) store.Store) *StoreProxy {
	return &StoreProxy{
		name:      name,
		superName: &superName,
		factory:   factory,
	}
}

// GetName returns the type name.
func (sp *StoreProxy) GetName() string {
	return sp.name
}

// GetSuper returns the supertype name.
func (sp *StoreProxy) GetSuper() *string {
	return sp.superName
}

// NewInstance creates a new instance using the factory function.
func (sp *StoreProxy) NewInstance(id oberon.Integer) store.Store {
	return sp.factory(id)
}

// Register is a helper function to register a type.
func Register(name string, factory func(oberon.Integer) store.Store) {
	proxy := NewStoreProxy(name, factory)
	GetInstance().Add(name, proxy)
}

// RegisterWithSuper is a helper function to register a type with a supertype.
func RegisterWithSuper(name string, superName string, factory func(oberon.Integer) store.Store) {
	proxy := NewStoreProxyWithSuper(name, superName, factory)
	GetInstance().Add(name, proxy)
}
