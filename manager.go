package sensitive

import (
	"github.com/sgoware/go-sensitive/filter"
	"github.com/sgoware/go-sensitive/store"
)

type Manager struct {
	store.Store
	filter.Filter
}

func NewFilter(storeOption StoreOption, filterOption FilterOption) *Manager {
	var filterStore store.Store
	var myFilter filter.Filter

	switch storeOption.Type {
	case StoreMemory:
		filterStore = store.NewMemoryModel()
	default:
		panic("invalid store type")
	}

	switch filterOption.Type {
	case FilterDfa:
		dfaModel := filter.NewDfaModel()

		go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())

		myFilter = dfaModel
	case FilterAc:
		acModel := filter.NewAcModel()

		go acModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())

		myFilter = acModel
	default:
		panic("invalid filter type")
	}

	return &Manager{
		Store:  filterStore,
		Filter: myFilter,
	}
}
