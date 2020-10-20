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
	s.define(&builtinTypeSymbol{name: "integer"})
	s.define(&builtinTypeSymbol{name: "real"})
}

func NewSymbolTable() *SymbolTable {
	st := &SymbolTable{
		symbols: make(map[string]Symbol),
	}
	st.initBuiltins()
	return st
}

type SemanticAnalyzer struct {
	*SymbolTable
}

func (sb *SemanticAnalyzer) VisitBlock(node *block) {
	for _, declaration := range node.declarations {
		sb.VisitNode(declaration)
	}
	sb.VisitNode(node.compoundStatement)
}

func (sb *SemanticAnalyzer) visitProgram(node *program) interface{} {
	sb.VisitBlock(node.block)

	return nil
}

func (sb *SemanticAnalyzer) visitBinOp(node *BinOp) interface{} {
	sb.VisitNode(node.left)
	sb.VisitNode(node.right)
	return nil
}

func (sb *SemanticAnalyzer) visitNum(_ *Num) interface{} { return nil }

func (sb *SemanticAnalyzer) VisitUnaryOp(node *UnaryOp) interface{} {
	sb.VisitNode(node.expr)
	return nil
}

func (sb *SemanticAnalyzer) VisitCompound(node *Compound) interface{} {
	for _, child := range node.children {
		sb.VisitNode(child)
	}
	return nil
}

func (sb *SemanticAnalyzer) VisitNoOp(_ *NoOp) {}

func (sb *SemanticAnalyzer) VisitVarDecl(node *varDecl) {
	typeName, _ := node.typeNode.Value()
	typeSymbol := sb.lookup(typeName.(string))
	varName, _ := node.varNode.Value()
	varSymbol := &varSymbol{
		name: varName.(string),
		typ:  typeSymbol,
	}
	sb.define(varSymbol)
}

func (sb *SemanticAnalyzer) visitAssign(node *assign) {
	varName, _ := node.left.Token().value.(string)
	varSymbol := sb.lookup(varName)
	if varSymbol == nil {
		panic("reference before assignment")
	}

	sb.VisitNode(node.right)
}

func (sb SemanticAnalyzer) visitVar(node *Var) interface{} {
	varName, _ := node.Token().value.(string)
	varSymbol := sb.lookup(varName)

	if varSymbol == nil {
		panic("reference before assignment")
	}
	return nil
}

func (sb *SemanticAnalyzer) VisitNode(node Node) interface{} {
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
	case *procDecl:
		sb.VisitProcedureDec(v)
	case *typeNode:
		sb.VisitType(v)
	case *program:
		return sb.visitProgram(v)
	default:
		panic(fmt.Sprintf("unexpected type occurrence %T", v))
	}
	return nil
}

func (sb *SemanticAnalyzer) VisitProcedureDec(_ *procDecl) {}

func (sb *SemanticAnalyzer) VisitType(_ *typeNode) {}
