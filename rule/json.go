package rule

import (
	"encoding/json"
	"errors"

	"github.com/tidwall/gjson"
)

type operator struct {
	kind     string
	operands []Expr
}

func (o *operator) UnmarshalJSON(data []byte) error {
	var node struct {
		Kind     string
		Operands operands
	}

	err := json.Unmarshal(data, &node)
	if err != nil {
		return err
	}

	o.operands = node.Operands.Exprs
	o.kind = node.Kind

	return nil
}

func (o *operator) MarshalJSON() ([]byte, error) {
	var op struct {
		Kind     string `json:"kind"`
		Operands []Expr `json:"operands"`
	}

	op.Kind = o.kind
	op.Operands = o.operands
	return json.Marshal(&op)
}

func (o *operator) Eval(params Params) (*Value, error) {
	return nil, nil
}

func (o *operator) Operands() []Expr {
	return o.operands
}

type operands struct {
	Ops   []json.RawMessage `json:"operands"`
	Exprs []Expr
}

func (o *operands) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &o.Ops)
	if err != nil {
		return err
	}

	for _, op := range o.Ops {
		r := gjson.Get(string(op), "kind")
		n, err := unmarshalExpr(r.Str, []byte(op))
		if err != nil {
			return err
		}

		o.Exprs = append(o.Exprs, n)
	}

	return nil
}

func unmarshalExpr(kind string, data []byte) (Expr, error) {
	var e Expr
	var err error

	switch kind {
	case "value":
		var v Value
		e = &v
		err = json.Unmarshal(data, &v)
	case "param":
		var p Param
		e = &p
		err = json.Unmarshal(data, &p)
	case "eq":
		var eq exprEq
		e = &eq
		err = eq.UnmarshalJSON(data)
	case "in":
		var in exprIn
		e = &in
		err = in.UnmarshalJSON(data)
	case "not":
		var not exprNot
		e = &not
		err = not.UnmarshalJSON(data)
	case "and":
		var and exprAnd
		e = &and
		err = and.UnmarshalJSON(data)
	case "or":
		var or exprOr
		e = &or
		err = or.UnmarshalJSON(data)
	case "gt":
		var gt exprGT
		e = &gt
		err = gt.UnmarshalJSON(data)
	case "gte":
		var gte exprGTE
		e = &gte
		err = gte.UnmarshalJSON(data)
	case "percentile":
		var percentile exprPercentile
		e = &percentile
		err = percentile.UnmarshalJSON(data)
	case "fnv":
		var fnv exprFNV
		e = &fnv
		err = fnv.UnmarshalJSON(data)
	case "lt":
		var lt exprLT
		e = &lt
		err = lt.UnmarshalJSON(data)
	case "lte":
		var lte exprLTE
		e = &lte
		err = lte.UnmarshalJSON(data)
	default:
		err = errors.New("unknown expression kind " + kind)
	}

	return e, err
}
