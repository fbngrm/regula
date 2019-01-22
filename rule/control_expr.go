package rule

func init() {
	Operators["let"] = func() Operator { return newExprLet() }
}

type exprLet struct {
	operator
}

func newExprLet() *exprLet {
	return &exprLet{
		operator: operator{
			contract: Contract{
				OpCode:     "let",
				ReturnType: ANY,
				Terms: []Term{
					{
						Type:        ANY,
						Cardinality: ONE,
					},
					{
						Type:        ANY,
						Cardinality: ONE,
					},
					{
						Type:        ANY,
						Cardinality: ONE,
						IsBody:      true,
					},
				},
			},
		},
	}
}

// Let creates an expression that defines a new scope with an
// additional parameter.  The Parameter will have the type and value
// of the value expression passed in position 2 and the body form
// (position 3) will be evaluated with the new scope, and therefore
// will be able to reference the new parameter.  The new parameter
// will not be available outside the scope of this Let expression.
// The new parameter may not share a name with a pre-existing parameter.
func Let(parameter Expr, value Expr, body Expr) Expr {
	e := newExprLet()
	e.pushExprOrPanic(parameter)
	e.pushExprOrPanic(value)
	e.pushExprOrPanic(body)
	e.finaliseOrPanic()
	return e
}

// Eval makes exprLet comply with the Expr interface.
func (n *exprLet) Eval(params Params) (*Value, error) {
	// Note we don't evaluate the symbol in position 0.  It will
	// be passed as a Param, but it isn't resolvable outside the
	// scoped Params we create below.
	symb := n.operands[0].(*Param)

	// The Expression in position 1 is the value form we'll bind
	// to the new parameter below.
	val, err := n.operands[1].Eval(params)
	if err != nil {
		return nil, err
	}

	// Create a new scoped Params with the symbol added
	var scopedParams *stack
	switch symb.Type {
	case "string":
		scopedParams = newStack(symb.Name, val.Data, params)
	case "int64":
		i, err := exprToInt64(val, params)
		if err != nil {
			return nil, err
		}
		scopedParams = newStack(symb.Name, i, params)
	case "float64":
		f, err := exprToFloat64(val, params)
		if err != nil {
			return nil, err
		}
		scopedParams = newStack(symb.Name, f, params)
	case "bool":
		b, err := exprToBool(val, params)
		if err != nil {
			return nil, err
		}
		scopedParams = newStack(symb.Name, b, params)
	}

	// Evaluate the body form within the new scope
	return n.operands[2].Eval(scopedParams)
}