package calc2

import (
	"fmt"
	"math"
	"strings"
)

type Env map[Var]float64

type Expr interface {
	Eval(env Env) float64
	Check(vars map[Var]bool) error
}

type Var string

func (v Var) Eval(env Env) float64 {
	return env[v]
}

func (v Var) Check(vars map[Var]bool) error {
	vars[v] = true
	return nil
}

type literal float64

func (l literal) Eval(_ Env) float64 {
	return float64(l)
}

func (l literal) Check(_ map[Var]bool) error {
	return nil
}

type unary struct {
	op rune
	x  Expr
}

func (u unary) Eval(env Env) float64 {
	switch u.op {
	case '+':
		return +u.x.Eval(env)
	case '-':
		return -u.x.Eval(env)
	default:
		panic("неподдерживаемый унарный оператор")
	}
	return 0
}

func (u unary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-", u.op) {
		return fmt.Errorf(fmt.Sprintf("некорректный унарный оператор %v", u.op))
	}
	return u.x.Check(vars)
}

type binary struct {
	op   rune
	x, y Expr
}

func (b binary) Eval(env Env) float64 {
	switch b.op {
	case '+':
		return b.x.Eval(env) + b.y.Eval(env)
	case '-':
		return b.x.Eval(env) - b.y.Eval(env)
	case '*':
		return b.x.Eval(env) * b.y.Eval(env)
	case '/':
		return b.x.Eval(env) / b.y.Eval(env)
	default:
		panic("неподдерживаемый бинарный оператор")
	}
}

func (b binary) Check(vars map[Var]bool) error {
	if !strings.ContainsRune("+-*/", b.op) {
		return fmt.Errorf(fmt.Sprintf("некорректный бинарный оператор %v", b.op))
	}
	if err := b.x.Check(vars); err != nil {
		return err
	}
	return b.y.Check(vars)
}

type call struct {
	fn   string
	args []Expr
}

func (c call) Eval(env Env) float64 {
	switch c.fn {
	case "pow":
		return math.Pow(c.args[0].Eval(env), c.args[1].Eval(env))
	case "sqrt":
		return math.Sqrt(c.args[0].Eval(env))
	case "sin":
		return math.Sin(c.args[0].Eval(env))
	default:
		panic("неподдерживаемая функция")
	}
}

var numParams = map[string]int{"pow": 2, "sqrt": 1, "sin": 1}

func (c call) Check(vars map[Var]bool) error {
	arity, ok := numParams[c.fn]
	if !ok {
		return fmt.Errorf("неизвестная функция %v", c.fn)
	}

	if len(c.args) != arity {
		return fmt.Errorf("неверное количество аргументов")
	}

	var err error
	for _, arg := range c.args {
		if err = arg.Check(vars); err != nil {
			return err
		}
	}
	return err
}
