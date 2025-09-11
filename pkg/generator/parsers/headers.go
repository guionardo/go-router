package parsers

type Header[T any] struct {
	Base[T]
}

func NewHeader[T any]() *Header[T] {
	p := &Header[T]{}
	p.readFields("header")
	p.fillImports(true)

	return p
}
