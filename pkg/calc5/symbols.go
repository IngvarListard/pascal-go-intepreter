package calc5

//type Symbol struct {
//	name string
//	typ interface{}
//}

type Symbol interface {
	Name() string
	Type() interface{}
}

type symbol struct {
	name string
	typ  interface{}
}

func (s *symbol) Name() string { return s.name }

func (s *symbol) Type() interface{} { return s.typ }

type builtinTypeSymbol struct {
	name string
}

func (b *builtinTypeSymbol) Name() string { return b.name }

func (b *builtinTypeSymbol) Type() interface{} { return nil }

// может быть второй тип и интерфейс не нужны, пока не понятно
