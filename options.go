package sensitive

const (
	StoreMemory = iota
)

const (
	FilterDfa = iota
)

type StoreOption struct {
	Type uint32
	Dsn  string
}

type FilterOption struct {
	Type uint32
}
