package rule

import (
	"errors"
	"fmt"
	"go/token"
	"strconv"

	"hash/fnv"
)

// An Expr is a logical expression that can be evaluated to a value.
type Expr interface {
	Eval(Params) (*Value, error)
}

// A Params is a set of parameters passed on rule evaluation.
// It provides type safe methods to query params.
type Params interface {
	GetString(key string) (string, error)
	GetBool(key string) (bool, error)
	GetInt64(key string) (int64, error)
	GetFloat64(key string) (float64, error)
	Keys() []string
	EncodeValue(key string) (string, error)
}

type exprNot struct {
	operator
}

// Not creates an expression that evaluates the given operand e and returns its opposite.
// e must evaluate to a boolean.
func Not(e Expr) Expr {
	return &exprNot{
		operator: operator{
			kind:     "not",
			operands: []Expr{e},
		},
	}
}

func (n *exprNot) Eval(params Params) (*Value, error) {
	if len(n.operands) < 1 {
		return nil, errors.New("invalid number of operands in Not func")
	}

	op := n.operands[0]
	v, err := op.Eval(params)
	if err != nil {
		return nil, err
	}

	if v.Type != "bool" {
		return nil, errors.New("invalid operand type for Not func")
	}

	if v.Equal(BoolValue(true)) {
		return BoolValue(false), nil
	}

	return BoolValue(true), nil
}

type exprOr struct {
	operator
}

// Or creates an expression that takes at least two operands and evaluates to true if one of the operands evaluates to true.
// All the given operands must evaluate to a boolean.
func Or(v1, v2 Expr, vN ...Expr) Expr {
	return &exprOr{
		operator: operator{
			kind:     "or",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprOr) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in Or func")
	}

	opA := n.operands[0]
	vA, err := opA.Eval(params)
	if err != nil {
		return nil, err
	}
	if vA.Type != "bool" {
		return nil, errors.New("invalid operand type for Or func")
	}

	if vA.Equal(BoolValue(true)) {
		return vA, nil
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}
		if vB.Type != "bool" {
			return nil, errors.New("invalid operand type for Or func")
		}

		if vB.Equal(BoolValue(true)) {
			return vB, nil
		}
	}

	return BoolValue(false), nil
}

type exprAnd struct {
	operator
}

// And creates an expression that takes at least two operands and evaluates to true if all the operands evaluate to true.
// All the given operands must evaluate to a boolean.
func And(v1, v2 Expr, vN ...Expr) Expr {
	return &exprAnd{
		operator: operator{
			kind:     "and",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprAnd) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in And func")
	}

	opA := n.operands[0]
	vA, err := opA.Eval(params)
	if err != nil {
		return nil, err
	}
	if vA.Type != "bool" {
		return nil, errors.New("invalid operand type for And func")
	}

	if vA.Equal(BoolValue(false)) {
		return vA, nil
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}
		if vB.Type != "bool" {
			return nil, errors.New("invalid operand type for And func")
		}

		if vB.Equal(BoolValue(false)) {
			return vB, nil
		}
	}

	return BoolValue(true), nil
}

type exprEq struct {
	operator
}

// Eq creates an expression that takes at least two operands and evaluates to true if all the operands are equal.
func Eq(v1, v2 Expr, vN ...Expr) Expr {
	return &exprEq{
		operator: operator{
			kind:     "eq",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprEq) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in Eq func")
	}

	opA := n.operands[0]
	vA, err := opA.Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		if !vA.Equal(vB) {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

type exprIn struct {
	operator
}

// In creates an expression that takes at least two operands and evaluates to true if the first one is equal to one of the others.
func In(v, e1 Expr, eN ...Expr) Expr {
	return &exprIn{
		operator: operator{
			kind:     "in",
			operands: append([]Expr{v, e1}, eN...),
		},
	}
}

func (n *exprIn) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in In func")
	}

	toFind := n.operands[0]
	vA, err := toFind.Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		if vA.Equal(vB) {
			return BoolValue(true), nil
		}
	}

	return BoolValue(false), nil
}

type exprGT struct {
	operator
}

// GT creates an expression that takes at least two operands and
// evaluates to true if each successive operand has a higher value than
// the next.
func GT(v1, v2 Expr, vN ...Expr) Expr {
	return &exprGT{
		operator: operator{
			kind:     "gt",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprGT) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in GT func")
	}

	vA, err := n.operands[0].Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		res, err := vA.GT(vB)
		if err != nil {
			return nil, err
		}

		if !res {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

type exprGTE struct {
	operator
}

// GTE creates an expression that takes at least two operands and
// evaluates to true if each successive operand has a greater or equal value
// compared to the next.
func GTE(v1, v2 Expr, vN ...Expr) Expr {
	return &exprGTE{
		operator: operator{
			kind:     "gte",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprGTE) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in GTE func")
	}

	vA, err := n.operands[0].Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		res, err := vA.GTE(vB)
		if err != nil {
			return nil, err
		}

		if !res {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

type exprLT struct {
	operator
}

// LT creates an expression that takes at least two operands and
// evaluates to true if each successive operand has a lower value
// compared to the next.
func LT(v1, v2 Expr, vN ...Expr) Expr {
	return &exprLT{
		operator: operator{
			kind:     "lt",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprLT) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in LT func")
	}

	vA, err := n.operands[0].Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		res, err := vA.LT(vB)
		if err != nil {
			return nil, err
		}

		if !res {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

type exprLTE struct {
	operator
}

// LTE creates an expression that takes at least two operands and
// evaluates to true if each successive operand has a lower or equal value
// compared to the next.
func LTE(v1, v2 Expr, vN ...Expr) Expr {
	return &exprLTE{
		operator: operator{
			kind:     "lte",
			operands: append([]Expr{v1, v2}, vN...),
		},
	}
}

func (n *exprLTE) Eval(params Params) (*Value, error) {
	if len(n.operands) < 2 {
		return nil, errors.New("invalid number of operands in LTE func")
	}

	vA, err := n.operands[0].Eval(params)
	if err != nil {
		return nil, err
	}

	for i := 1; i < len(n.operands); i++ {
		vB, err := n.operands[i].Eval(params)
		if err != nil {
			return nil, err
		}

		res, err := vA.LTE(vB)
		if err != nil {
			return nil, err
		}

		if !res {
			return BoolValue(false), nil
		}
	}

	return BoolValue(true), nil
}

type exprFNV struct {
	operator
}

// FNV returns an Integer hash of any value it is provided.  It uses
// the Fowler-Noll-Vo non-cryptographic hash function.
func FNV(v Expr) Expr {
	return &exprFNV{
		operator: operator{
			kind:     "fnv",
			operands: []Expr{v},
		},
	}
}

func (n *exprFNV) Eval(params Params) (*Value, error) {
	if len(n.operands) != 1 {
		return nil, errors.New("invalid number of operands in FNV func")
	}

	h32 := fnv.New32()
	op := n.operands[0]
	v, err := op.Eval(params)
	if err != nil {
		return nil, err
	}
	_, err = h32.Write([]byte(v.Data))
	if err != nil {
		return nil, err
	}
	return Int64Value(int64(h32.Sum32())), nil
}

type exprPercentile struct {
	operator
}

// Percentile indicates whether the provided value is within a given
// percentile of the group of all such values.  It is intended to be
// used to assign values to groups for experimentation.
func Percentile(v, p Expr) Expr {
	return &exprPercentile{
		operator: operator{
			kind:     "percentile",
			operands: []Expr{v, p},
		},
	}
}

func (n *exprPercentile) Eval(params Params) (*Value, error) {
	if len(n.operands) != 2 {
		return nil, errors.New("invalid number of operands in Percentile func")
	}

	hash := FNV(n.operands[0])
	v, err := exprToInt64(hash, params)
	if err != nil {
		return nil, err
	}
	p, err := exprToInt64(n.operands[1], params)
	if err != nil {
		return nil, err
	}
	if (v % 100) <= p {
		return BoolValue(true), nil
	}
	return BoolValue(false), nil
}

// Param is an expression used to select a parameter passed during evaluation and return its corresponding value.
type Param struct {
	Kind string `json:"kind"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// StringParam creates a Param that looks up in the set of params passed during evaluation and returns the value
// of the variable that corresponds to the given name.
// The corresponding value must be a string. If not found it returns an error.
func StringParam(name string) *Param {
	return &Param{
		Kind: "param",
		Type: "string",
		Name: name,
	}
}

// BoolParam creates a Param that looks up in the set of params passed during evaluation and returns the value
// of the variable that corresponds to the given name.
// The corresponding value must be a boolean. If not found it returns an error.
func BoolParam(name string) *Param {
	return &Param{
		Kind: "param",
		Type: "bool",
		Name: name,
	}
}

// Int64Param creates a Param that looks up in the set of params passed during evaluation and returns the value
// of the variable that corresponds to the given name.
// The corresponding value must be an int64. If not found it returns an error.
func Int64Param(name string) *Param {
	return &Param{
		Kind: "param",
		Type: "int64",
		Name: name,
	}
}

// Float64Param creates a Param that looks up in the set of params passed during evaluation and returns the value
// of the variable that corresponds to the given name.
// The corresponding value must be a float64. If not found it returns an error.
func Float64Param(name string) *Param {
	return &Param{
		Kind: "param",
		Type: "float64",
		Name: name,
	}
}

// Eval extracts a value from the given parameters.
func (p *Param) Eval(params Params) (*Value, error) {
	if params == nil {
		return nil, errors.New("params is nil")
	}

	switch p.Type {
	case "string":
		v, err := params.GetString(p.Name)
		if err != nil {
			return nil, err
		}
		return StringValue(v), nil
	case "bool":
		v, err := params.GetBool(p.Name)
		if err != nil {
			return nil, err
		}
		return BoolValue(v), nil
	case "int64":
		v, err := params.GetInt64(p.Name)
		if err != nil {
			return nil, err
		}
		return Int64Value(v), nil
	case "float64":
		v, err := params.GetFloat64(p.Name)
		if err != nil {
			return nil, err
		}
		return Float64Value(v), nil
	}

	return nil, errors.New("unsupported param type")
}

// True creates an expression that always evaluates to true.
func True() Expr {
	return BoolValue(true)
}

// A Value is the result of the evaluation of an expression.
type Value struct {
	Kind string `json:"kind"`
	Type string `json:"type"`
	Data string `json:"data"`
}

func newValue(typ, data string) *Value {
	return &Value{
		Kind: "value",
		Type: typ,
		Data: data,
	}
}

// BoolValue creates a bool type value.
func BoolValue(value bool) *Value {
	return newValue("bool", strconv.FormatBool(value))
}

// StringValue creates a string type value.
func StringValue(value string) *Value {
	return newValue("string", value)
}

// Int64Value creates an int64 type value.
func Int64Value(value int64) *Value {
	return newValue("int64", strconv.FormatInt(value, 10))
}

// Float64Value creates a float64 type value.
func Float64Value(value float64) *Value {
	return newValue("float64", strconv.FormatFloat(value, 'f', 6, 64))
}

// Eval evaluates the value to itself.
func (v *Value) Eval(Params) (*Value, error) {
	return v, nil
}

func (v *Value) compare(op token.Token, other *Value) bool {
	if op != token.EQL {
		return false
	}

	return *v == *other
}

// Equal reports whether v and other represent the same value.
func (v *Value) Equal(other *Value) bool {
	return v.compare(token.EQL, other)
}

// GT reports whether v is greater than other.
func (v *Value) GT(other *Value) (bool, error) {
	switch v.Type {
	case "bool":
		v1, v2, err := parseBoolValues(v, other)
		if err != nil {
			return false, err
		}

		if !v1 {
			// If v1 is False then it's not greater than v2, and we can be done already.
			return false, nil
		}
		if v2 {
			// If v2 is True then v1 can't be greater than it..
			return false, nil
		}
		return true, nil
	case "string":
		if v.Data <= other.Data {
			return false, nil
		}
		return true, nil
	case "int64":
		v1, v2, err := parseInt64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 <= v2 {
			return false, nil
		}
		return true, nil
	case "float64":
		v1, v2, err := parseFloat64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 <= v2 {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown Value type: %s", v.Type)
}

// GTE reports whether v is greater or equal than other.
func (v *Value) GTE(other *Value) (bool, error) {
	switch v.Type {
	case "bool":
		v1, v2, err := parseBoolValues(v, other)
		if err != nil {
			return false, err
		}

		if !v1 && v2 {
			return false, nil
		}
		return true, nil
	case "string":
		if v.Data < other.Data {
			return false, nil
		}
		return true, nil
	case "int64":
		v1, v2, err := parseInt64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 < v2 {
			return false, nil
		}
		return true, nil
	case "float64":
		v1, v2, err := parseFloat64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 < v2 {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown Value type: %s", v.Type)
}

// LT reports whether v is less than other.
func (v *Value) LT(other *Value) (bool, error) {
	switch v.Type {
	case "bool":
		v1, v2, err := parseBoolValues(v, other)
		if err != nil {
			return false, err
		}

		if v1 {
			// If v1 is True then it's not less than v2, and we can be done already.
			return false, nil
		}
		if !v2 {
			// If v2 is False then v1 can't be less than it..
			return false, nil
		}
		return true, nil
	case "string":
		if v.Data >= other.Data {
			return false, nil
		}
		return true, nil
	case "int64":
		v1, v2, err := parseInt64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 >= v2 {
			return false, nil
		}
		return true, nil
	case "float64":
		v1, v2, err := parseFloat64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 >= v2 {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown Value type: %s", v.Type)
}

// LTE reports whether v is less or equal than other.
func (v *Value) LTE(other *Value) (bool, error) {
	switch v.Type {
	case "bool":
		v1, v2, err := parseBoolValues(v, other)
		if err != nil {
			return false, err
		}

		if v1 && !v2 {
			return false, nil
		}
		return true, nil
	case "string":
		if v.Data > other.Data {
			return false, nil
		}
		return true, nil
	case "int64":
		v1, v2, err := parseInt64Values(v, other)
		if err != nil {
			return false, nil
		}

		if v1 > v2 {
			return false, nil
		}
		return true, nil
	case "float64":
		v1, v2, err := parseFloat64Values(v, other)
		if err != nil {
			return false, err
		}

		if v1 > v2 {
			return false, nil
		}
		return true, nil
	}
	return false, fmt.Errorf("unknown Value type: %s", v.Type)
}

func parseBoolValues(v1, v2 *Value) (b1, b2 bool, err error) {
	if b1, err = strconv.ParseBool(v1.Data); err != nil {
		return
	}
	b2, err = strconv.ParseBool(v2.Data)
	return
}

func parseInt64Values(v1, v2 *Value) (i1, i2 int64, err error) {
	if i1, err = strconv.ParseInt(v1.Data, 10, 64); err != nil {
		return
	}
	i2, err = strconv.ParseInt(v2.Data, 10, 64)
	return
}

func parseFloat64Values(v1, v2 *Value) (f1, f2 float64, err error) {
	if f1, err = strconv.ParseFloat(v1.Data, 64); err != nil {
		return
	}
	f2, err = strconv.ParseFloat(v2.Data, 64)
	return
}

type operander interface {
	Operands() []Expr
}

func walk(expr Expr, fn func(Expr) error) error {
	err := fn(expr)
	if err != nil {
		return err
	}

	if o, ok := expr.(operander); ok {
		ops := o.Operands()
		for _, op := range ops {
			err := walk(op, fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// exprToInt64 returns the go-native int64 value of an expression
// evaluated with params.
func exprToInt64(e Expr, params Params) (int64, error) {
	v, err := e.Eval(params)
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(v.Data, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, err
}
