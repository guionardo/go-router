package parsers

type Path[T any] struct {
	Base[T]
}

func NewPath[T any]() *Path[T] {
	p := &Path[T]{}
	p.readFields("path")
	p.fillImports(true)

	return p
}
