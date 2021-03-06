package calc5

import (
	"bytes"
	"fmt"
	"strings"
)

type ScopedSymbolTable struct {
	symbols        map[string]Symbol
	scopeName      string
	scopeLevel     int
	enclosingScope *ScopedSymbolTable
}

func (s *ScopedSymbolTable) String() string {
	h1 := "SCOPE (SCOPED SYMBOL TABLE)"
	var sep bytes.Buffer
	for _ = range h1 {
		sep.WriteRune('=')
	}
	lines := []string{"\n", h1, sep.String()}

	scopeName := "nil"
	if s.enclosingScope != nil {
		scopeName = s.enclosingScope.scopeName
	}

	for k, v := range map[string]interface{}{
		"Scope name":      s.scopeName,
		"Scope level":     s.scopeLevel,
		"Enclosing scope": scopeName,
	} {
		lines = append(lines, fmt.Sprintf("%s: %v", k, v))
	}

	h2 := "Scope (Scoped symbol table) contents"
	sep = bytes.Buffer{}
	for _ = range h2 {
		sep.WriteRune('-')
	}
	lines = append(lines, h2, sep.String())

	for k, v := range s.symbols {
		lines = append(lines, fmt.Sprintf("%s: %s", k, v))
	}
	lines = append(lines, "\n")

	values := make([]string, len(s.symbols))
	cnt := 0
	for _, symbol := range s.symbols {
		values[cnt] = symbol.(fmt.Stringer).String()
		cnt++
	}
	return strings.Join(lines, "\n")
}

func (s *ScopedSymbolTable) define(symbol Symbol) {
	fmt.Printf("Define Symbol: %s\n", symbol)
	s.symbols[symbol.Name()] = symbol
}

func (s *ScopedSymbolTable) lookup(name string, currentScopeOnly bool) Symbol {
	fmt.Printf("Lookup Symbol: %s in %s Scope\n", name, s.scopeName)

	v, ok := s.symbols[name]
	if ok {
		return v
	}

	if currentScopeOnly {
		return nil
	}

	if s.enclosingScope != nil {
		return s.enclosingScope.lookup(name, false)
	}
	return nil
}

func (s *ScopedSymbolTable) initBuiltins() {
	s.define(&builtinTypeSymbol{name: "integer"})
	s.define(&builtinTypeSymbol{name: "real"})
}

func NewScopedSymbolTable(name string, level int, enclosingScope *ScopedSymbolTable) *ScopedSymbolTable {
	st := &ScopedSymbolTable{
		symbols:        make(map[string]Symbol),
		scopeName:      name,
		scopeLevel:     level,
		enclosingScope: enclosingScope,
	}
	st.initBuiltins()
	return st
}

type SemanticAnalyzer struct {
	*ScopedSymbolTable // currentScope?
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return new(SemanticAnalyzer)
}

func (sb *SemanticAnalyzer) VisitBlock(node *block) {
	for _, declaration := range node.declarations {
		sb.VisitNode(declaration)
	}
	sb.VisitNode(node.compoundStatement)
}

func (sb *SemanticAnalyzer) visitProgram(node *program) interface{} {
	fmt.Println("Enter scope: Global")
	globalScope := NewScopedSymbolTable("global", 1, sb.ScopedSymbolTable)
	sb.ScopedSymbolTable = globalScope
	sb.VisitNode(node.block)
	fmt.Println(globalScope)

	sb.ScopedSymbolTable = sb.enclosingScope
	fmt.Println("Leave scope: Global")
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
	typeSymbol := sb.lookup(typeName.(string), false)
	varName, _ := node.varNode.Value()
	varNameStr := varName.(string)
	varSymbol := &varSymbol{name: varNameStr, typ: typeSymbol}

	if sb.lookup(varNameStr, true) != nil {
		panic(fmt.Sprintf("Error: Duplicate identifier '%s' found", varNameStr))
	}

	sb.define(varSymbol)
}

func (sb *SemanticAnalyzer) visitAssign(node *assign) {
	varName, _ := node.left.Token().value.(string)
	varSymbol := sb.lookup(varName, false)
	if varSymbol == nil {
		panic("reference before assignment")
	}

	sb.VisitNode(node.right)
}

func (sb *SemanticAnalyzer) visitVar(node *Var) interface{} {
	varName, _ := node.Token().value.(string)
	varSymbol := sb.lookup(varName, false)

	if varSymbol == nil {
		panic("reference before assignment")
	}
	return nil
}

func (sb *SemanticAnalyzer) visitVarDecl(node *varDecl) {
	typeName := node.typeNode.(*typeNode).value
	typeSymbol := sb.lookup(typeName.(string), false)

	varName, _ := node.varNode.Value()
	varSymbol := varSymbol{name: varName.(string), typ: typeSymbol}

	if sb.lookup(varName.(string), true) != nil {
		panic(fmt.Sprintf("Error: Duplicate identifier '%s' found", varName))
	}

	sb.define(&varSymbol)
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

func (sb *SemanticAnalyzer) VisitProcedureDec(node *procDecl) {
	procName := node.procName
	procSymbol := &procedureSymbol{
		name: procName,
	}
	sb.define(procSymbol)
	fmt.Printf("Entering scope: %s", procName)
	procedureScope := NewScopedSymbolTable(procName, sb.scopeLevel+1, sb.ScopedSymbolTable)
	sb.ScopedSymbolTable = procedureScope

	for _, p := range node.params {
		paramType := sb.lookup(p.typeNode.value.(string), false)
		paramName := p.varNode.value

		varSymbol := &varSymbol{
			name: paramName.(string),
			typ:  paramType,
		}
		sb.define(varSymbol)
		procSymbol.params = append(procSymbol.params, varSymbol)
	}
	sb.VisitNode(node.block)

	fmt.Println(procedureScope)

	sb.ScopedSymbolTable = procedureScope
	fmt.Printf("Leave scope: %s", procName)
}

func (sb *SemanticAnalyzer) VisitType(_ *typeNode) {}
