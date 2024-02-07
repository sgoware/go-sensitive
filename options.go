package sensitive

import (
	"github.com/sgoware/go-sensitive/store"
)

const (
	StoreMemory = iota
	StoreMysql
	StoreMongo
)

const (
	FilterDfa = iota
	FilterAc
)

type StoreOption struct {
	Type        uint32
	MysqlConfig *store.MysqlConfig
	MongoConfig *store.MongoConfig
}

type FilterOption struct {
	Type uint32
}
