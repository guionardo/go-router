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

// func (p *Path[T]) Attributions() []string {
// 	attributions := make([]string, 0)

// 	for field, tagValue := range p.Fields() {
// 		attributions = append(attributions, createAttribuition(field, "h", `r.PathValue("%s")`, tagValue))
// 	}
// 	if len(attributions) > 0 {
// 		attributions = append([]string{"", "// path"}, attributions...)

// 	}
// 	return attributions
// }
// TODO: Remover código não utilizado
