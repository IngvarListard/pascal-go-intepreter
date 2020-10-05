package calc5

import (
	"fmt"
)

type Interpreter struct {
	parser      *Parser
	GlobalScope map[string]interface{}
}

func (i *Interpreter) interpret() interface{} {
	node := i.parser.parse()
	//v, _ := node.Value()
	//return v
	return i.VisitNode(node)
}

func (i *Interpreter) visitBinOp(binary *BinOp) interface{} {
	vl := i.VisitNode(binary.left)
	val := vl.(int)
	vr := i.VisitNode(binary.right)
	vll := vr.(int)
	switch binary.op.typ {
	case Plus:
		return val + vll
	case Minus:
		return val - vll
	case Mul:
		return val * vll
	case Div:
		return val / vll
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
