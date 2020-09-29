package calc5

type Interpreter struct {
	parser *Parser
}

func (i *Interpreter) interpret() interface{} {
	node := i.parser.parse()
	v, _ := node.Value()
	return v
}
