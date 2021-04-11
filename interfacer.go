package easyquery

type Preloader interface {
	Preload() []string
}

type Joinser interface {
	Joins() []interface{}
}

type Paginater interface {
	GetSize() int
	GetPage() int
	GetOffset() int
	GetPagable() bool
	GetTotal() int64
	SetTotal(total int64)
}

type QueryParamer interface {
	GetFields() []*QueryField
	GetPagination() Paginater
	GetJoin() bool
}
