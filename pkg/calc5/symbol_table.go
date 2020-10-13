package calc5

import (
	"fmt"
	"strings"
)

type SymbolTable struct {
	symbols map[string]Symbol
}

func (s *SymbolTable) String() string {
	values := make([]string, len(s.symbols))
	cnt := 0
	for _, symbol := range s.symbols {
		values[cnt] = symbol.(fmt.Stringer).String()
		cnt++
	}
	return fmt.Sprintf("Symbols: [%s]", strings.Join(values, ","))
}

func (s *SymbolTable) define(symbol Symbol) {
	fmt.Printf("Define Symbol: %s\n", symbol)
	s.symbols[symbol.Name()] = symbol
}

func (s *SymbolTable) lookup(name string) Symbol {
	fmt.Printf("Lookup Symbol: %s\n", name)
	return s.symbols[name]
}

func (s *SymbolTable) initBuiltins() {
	s.define(&builtinTypeSymbol{name: "INTEGER"})
	s.define(&builtinTypeSymbol{name: "REAL"})
}

func NewSymbolTable() *SymbolTable {
	st := new(SymbolTable)
	st.initBuiltins()
	return st
}

type SymbolTableBuilder struct {
	*SymbolTable
}

func (sb *SymbolTableBuilder) VisitBlock(node *block) {
	for _, declaration := range node.declarations {
		sb.VisitNode(declaration)
	}
	sb.VisitNode(node.compoundStatement)
}

func (sb *SymbolTableBuilder) visitProgram(node *program) interface{} {
	sb.VisitBlock(node.block)

	return nil
}

func (sb *SymbolTableBuilder) visitBinOp(node *BinOp) interface{} {
	sb.VisitNode(node.left)
	sb.VisitNode(node.right)
	return nil
}

func (sb *SymbolTableBuilder) visitNum(_ *Num) interface{} { return nil }

func (sb *SymbolTableBuilder) VisitUnaryOp(node *UnaryOp) interface{} {
	sb.VisitNode(node.expr)
	return nil
}

func (sb *SymbolTableBuilder) VisitCompound(node *Compound) interface{} {
	for _, child := range node.children {
		sb.VisitNode(child)
	}
	return nil
}

func (sb *SymbolTableBuilder) VisitNoOp(_ *NoOp) {}

func (sb *SymbolTableBuilder) VisitVarDecl(node *varDecl) {
	typeName, _ := node.typeNode.Value()
	typeSymbol := sb.lookup(typeName.(string))
	varName, _ := node.varNode.Value()
	varSymbol := &varSymbol{
		name: varName.(string),
		typ:  typeSymbol,
	}
	sb.define(varSymbol)
}

func (sb *SymbolTableBuilder) visitAssign(node *assign) {
	varName, _ := node.left.Value()
	varSymbol := sb.lookup(varName.(string))
	if varSymbol == nil {
		panic("reference before assignment")
	}

	sb.VisitNode(node.right)
}

func (sb SymbolTableBuilder) visitVar(node *Var) interface{} {
	varName, _ := node.Value()
	varSymbol := sb.lookup(varName.(string))

	if varSymbol == nil {
		panic("reference before assignment")
	}
	return nil
}

func (sb *SymbolTableBuilder) VisitNode(node Node) interface{} {
	switch v := node.(type) {
	case *BinOp:
		return sb.visitBinOp(v)
	case *Num:
		return sb.visitNum(v)
	case *UnaryOp:
		return sb.VisitUnaryOp(v)
	case *Compound:
		sb.VisitCompound(v)
	case *assign:
		sb.visitAssign(v)
	case *NoOp:
		sb.VisitNoOp(v)
	case *Var:
		return sb.visitVar(v)
	case *block:
		sb.VisitBlock(v)
	case *varDecl:
		sb.VisitVarDecl(v)
	case *typeNode:
		//sb.VisitType(v)
		return nil
	case *program:
		return sb.visitProgram(v)
	default:
		panic(fmt.Sprintf("unexpected type occurrence %T", v))
	}
	return nil
}
