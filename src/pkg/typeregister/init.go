// Package typeregister - type registration initialization
package typeregister

import (
	"odcread/pkg/fold"
	"odcread/pkg/store"
	"odcread/pkg/textmodel"
)

// init registers all known types with the TypeRegister.
func init() {
	// Register Store hierarchy
	Register(store.TypeNameStore, func(id int32) store.Store {
		return store.NewStore(id)
	})

	Register(store.TypeNameElem, func(id int32) store.Store {
		return store.NewElem(id)
	})

	Register(store.TypeNameModel, func(id int32) store.Store {
		return store.NewModel(id)
	})

	Register(store.TypeNameContainerModel, func(id int32) store.Store {
		return store.NewContainerModel(id)
	})

	// Register TextModel hierarchy
	Register(textmodel.TypeNameTextModel, func(id int32) store.Store {
		return textmodel.NewTextModel(id)
	})

	Register(textmodel.TypeNameStdTextModel, func(id int32) store.Store {
		return textmodel.NewStdTextModel(id)
	})

	// Register View/Fold hierarchy
	Register(fold.TypeNameView, func(id int32) store.Store {
		return fold.NewView(id)
	})

	Register(fold.TypeNameFold, func(id int32) store.Store {
		return fold.NewFold(id)
	})
}
