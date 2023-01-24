package sensitive

import (
	"go-sensitive/filter"
	"go-sensitive/store"
	"sync"
)

type Manager struct {
	store  store.Model
	filter filter.Model

	filterMux sync.RWMutex
}

func NewFilter(storeOption StoreOption, filterOption FilterOption) *Manager {
	var filterStore store.Model
	var myFilter filter.Model

	switch storeOption.Type {
	case StoreMemory:
		filterStore = store.NewMemoryModel()
	}

	switch filterOption.Type {
	case FilterDfa:
		dfaModel := filter.NewDfaModel()

		go dfaModel.Listen(filterStore.GetAddChan(), filterStore.GetDelChan())

		myFilter = dfaModel
	}

	return &Manager{
		store:  filterStore,
		filter: myFilter,
	}
}

func (m *Manager) GetStore() store.Model {
	return m.store
}

func (m *Manager) GetFilter() filter.Model {
	m.filterMux.RLock()
	myFilter := m.filter
	m.filterMux.RUnlock()
	return myFilter
}
