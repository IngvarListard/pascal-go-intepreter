package calc5

import (
	"fmt"
	"strconv"
)

type Interpreter struct {
	parser *Parser
}

func (i *Interpreter) interpret() interface{} {
	node := i.parser.parse()
	//v, _ := node.Value()
	//return v
	return i.VisitNode(node)
}

func (i *Interpreter) interpretPolish() interface{} {
	node := i.parser.parse()
	//v, _ := node.Value()
	//return v
	return i.VisitNodePolish(node)
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

/*
значит нужно написать так чтобы скобки опускались, а знак

(2+3) * (4 + 5)
2 3 + 4 * 5 +

(2+3) * 7
2 3 + 7 *
*/

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

func (i *Interpreter) VisitNodePolish(node Node) interface{} {
	switch v := node.(type) {
	case *BinOp:
		return i.visitBinOpPolish(v)
	case *Num:
		return i.visitNum(v)
	default:
		panic(fmt.Sprintf("unexpected type occurrence %T", v))
	}
}

func (i *Interpreter) visitBinOpPolish(binary *BinOp) interface{} {
	var leftStr string
	var rightStr string

	vl := i.VisitNodePolish(binary.left)
	switch vvv := vl.(type) {
	case int:
		leftStr = strconv.Itoa(vvv)
	case string:
		leftStr = vvv
	}
	vr := i.VisitNodePolish(binary.right)
	switch vvv := vr.(type) {
	case int:
		rightStr = strconv.Itoa(vvv)
	case string:
		rightStr = vvv
	}
	template := "%s %s %s"
	switch binary.op.typ {
	case Plus:

		return fmt.Sprintf(template, leftStr, rightStr, "+")
	case Minus:
		return fmt.Sprintf(template, leftStr, rightStr, "-")
	case Mul:
		return fmt.Sprintf(template, leftStr, rightStr, "*")
	case Div:
		return fmt.Sprintf(template, leftStr, rightStr, "/")
	default:
		panic("AAA")
	}
}
