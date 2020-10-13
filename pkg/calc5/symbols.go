package calc5

import "fmt"

type Symbol interface {
	Name() string
	Type() Symbol
}

type builtinTypeSymbol struct {
	name string
}

func (b *builtinTypeSymbol) Name() string { return b.name }

func (b *builtinTypeSymbol) Type() Symbol { return nil }

func (b *builtinTypeSymbol) String() string { return b.name }

type varSymbol struct {
	name string
	typ  Symbol
}

func (v *varSymbol) Name() string { return v.name }

func (v *varSymbol) Type() Symbol { return v.typ }

func (v *varSymbol) String() string { return fmt.Sprintf("<%v:%v>", v.name, v.typ) }
