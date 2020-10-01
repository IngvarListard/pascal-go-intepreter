package calc5

import "fmt"

type Interpreter struct {
	parser *Parser
}

func (i *Interpreter) interpret() interface{} {
	node := i.parser.parse()
	v, _ := node.Value()
	return v
}

func (i *Interpreter) visitBinOp(binary *BinOp) interface{} {
	switch binary.op.typ {
	case Plus:
		vl := i.VisitNode(binary.left)
		val := vl.(int)
		vr := i.VisitNode(binary.right)
		vll := vr.(int)
		return val + vll
	case Minus:
		vl := i.VisitNode(binary.left)
		val := vl.(int)
		vr := i.VisitNode(binary.right)
		vll := vr.(int)
		return val - vll
	case Mul:
		vl := i.VisitNode(binary.left)
		val := vl.(int)
		vr := i.VisitNode(binary.right)
		vll := vr.(int)
		return val * vll
	case Div:
		vl := i.VisitNode(binary.left)
		val := vl.(int)
		vr := i.VisitNode(binary.right)
		vll := vr.(int)
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
	default:
		panic(fmt.Sprintf("unexpected type occurrence %T", v))
	}
}
