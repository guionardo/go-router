package parsers

type Query[T any] struct {
	Base[T]
}

func NewQuery[T any]() *Query[T] {
	p := &Query[T]{}
	p.readFields("query")
	p.fillImports(true)

	return p
}
