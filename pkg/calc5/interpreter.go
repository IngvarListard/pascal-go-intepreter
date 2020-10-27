package calc5

import (
	"fmt"
	"reflect"
)

type Interpreter struct {
	parser      *Parser
	GlobalScope map[string]interface{}
	Symbols     *SemanticAnalyzer
}

func (i *Interpreter) interpret() interface{} {
	node := i.parser.parse()
	i.Symbols = new(SemanticAnalyzer)
	i.Symbols.VisitNode(node)
	return i.VisitNode(node)
}

func (i *Interpreter) visitBinOp(binary *BinOp) interface{} {
	vl := i.VisitNode(binary.left)
	vr := i.VisitNode(binary.right)

	lTyp := reflect.TypeOf(vl).Kind()
	rTyp := reflect.TypeOf(vr).Kind()

	switch {
	case lTyp == reflect.Int && rTyp == reflect.Int:
		return execIntOp(vl.(int), vr.(int), binary.op.typ)
	case lTyp == reflect.Float64 || rTyp == reflect.Float64:
		left := getFloat(vl)
		right := getFloat(vr)
		return execFloatOp(left, right, binary.op.typ)
	}

	panic("unexpected types")
}

func getFloat(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case float64:
		return val
	default:
		panic(fmt.Sprintf("unexpected type %T", v))
	}
}

func execFloatOp(left, right float64, op TokenTyp) float64 {
	switch op {
	case Plus:
		return left + right
	case Minus:
		return left - right
	case Mul:
		return left * right
	case IntegerDiv:
		return left / right
	case FloatDiv:
		return left / right
	default:
		panic("AAA")
	}
}

func execIntOp(left, right int, op TokenTyp) int {
	switch op {
	case Plus:
		return left + right
	case Minus:
		return left - right
	case Mul:
		return left * right
	case IntegerDiv:
		return left / right
	case FloatDiv:
		return left / right
	default:
		panic("AAA")
	}
}

func (i *Interpreter) visitNum(num *Num) interface{} {
	return num.value
}

func (i *Interpreter) VisitNode(node Node) interface{} {
	switch v := node.(type) {
	case *BinOp:
		return i.visitBinOp(v)
	case *Num:
		return i.visitNum(v)
	case *UnaryOp:
		return i.VisitUnaryOp(v)
	case *Compound:
		i.VisitCompound(v)
	case *assign:
		i.VisitAssign(v)
	case *NoOp:
		i.VisitNoOp(v)
	case *Var:
		return i.VisitVar(v)
	case *block:
		i.VisitBlock(v)
	case *varDecl:
		i.VisitVarDecl(v)
	case *procDecl:
		i.VisitProcedureDec(v)
	case *typeNode:
		i.VisitType(v)
	case *program:
		return i.VisitProgram(v)
	default:
		panic(fmt.Sprintf("unexpected type occurrence %T", v))
	}
	return nil
}

func (i *Interpreter) VisitUnaryOp(node *UnaryOp) interface{} {
	switch node.Token().typ {
	case Plus:
		return +i.VisitNode(node.expr).(int)
	case Minus:
		return -i.VisitNode(node.expr).(int)
	default:
		panic("Unexpected")
	}
}

func (i *Interpreter) VisitCompound(node *Compound) {
	for _, child := range node.children {
		i.VisitNode(child)
	}
}

func (i *Interpreter) VisitAssign(node *assign) {
	nodeName := node.left.(*Var).token.value
	n := nodeName.(string)
	v := i.VisitNode(node.right)
	i.GlobalScope[n] = v
}

func (i *Interpreter) VisitVar(node *Var) interface{} {
	name := node.token.value
	val, ok := i.GlobalScope[name.(string)]
	if !ok {
		panic(fmt.Sprintf("no %s name in global scope", name.(string)))
	}
	return val
}

func (i *Interpreter) VisitNoOp(_ *NoOp) {}

func (i *Interpreter) VisitProgram(node *program) interface{} {
	return i.VisitNode(node.block)
}

func (i *Interpreter) VisitBlock(node *block) {
	for _, declaration := range node.declarations {
		i.VisitNode(declaration)
	}
	i.VisitNode(node.compoundStatement)
}

func (i *Interpreter) VisitVarDecl(_ *varDecl) {}

func (i *Interpreter) VisitType(_ *typeNode) {}

func (i *Interpreter) VisitProcedureDec(_ *procDecl) {}
