package rule_test

import (
	"fmt"
	"testing"

	"github.com/heetch/regula/rule"
	"github.com/stretchr/testify/require"
)

// IsFulfilledBy indicates whether a given TypedExpression fulfils a Term.
func TestTermIsFulfilledBy(t *testing.T) {
	// We'll express a few fundamental TypedExpressions, to avoid repetition in the test cases.
	boolean := rule.BoolValue(true)
	integer := rule.Int64Value(1)
	float := rule.Float64Value(1.1)
	str := rule.StringValue("foo")
	boolParam := rule.BoolParam("foo")
	stringParam := rule.StringParam("foo")
	intParam := rule.Int64Param("foo")
	floatParam := rule.Float64Param("foo")
	not := rule.Not(boolean).(rule.TypedExpression)
	or := rule.Or(boolean, boolean).(rule.TypedExpression)
	and := rule.And(boolean, boolean).(rule.TypedExpression)
	eq := rule.Eq(boolean, boolean).(rule.TypedExpression)
	in := rule.In(boolean, boolean, boolean).(rule.TypedExpression)

	testCases := []struct {
		// Test cases define a list of positive expressions
		// (those that fulfil the Term), and a list of
		// negative expressions (those that do not fulfil the
		// Term).  We are attempting to be exhaustive here.
		name                string
		positiveExpressions []rule.TypedExpression
		negativeExpressions []rule.TypedExpression
		term                rule.Term
	}{
		{
			name: "Boolean",
			positiveExpressions: []rule.TypedExpression{
				boolean, boolParam, not, or, and, eq, in},
			negativeExpressions: []rule.TypedExpression{
				str, integer, float, stringParam, intParam, floatParam},
			term: rule.Term{Type: rule.BOOLEAN},
		},
		{
			name:                "String",
			positiveExpressions: []rule.TypedExpression{str, stringParam},
			negativeExpressions: []rule.TypedExpression{
				boolean, integer, float, boolParam, intParam, floatParam,
				or, and, eq, in, not},
			term: rule.Term{Type: rule.STRING},
		},
		{
			name:                "Integer",
			positiveExpressions: []rule.TypedExpression{integer, intParam},
			negativeExpressions: []rule.TypedExpression{
				boolean, str, float, boolParam, stringParam, floatParam,
				or, and, eq, in, not},
			term: rule.Term{Type: rule.INTEGER},
		},
		{
			name:                "Float",
			positiveExpressions: []rule.TypedExpression{float, floatParam},
			negativeExpressions: []rule.TypedExpression{
				boolean, str, integer, boolParam, stringParam, intParam,
				or, and, eq, in, not},
			term: rule.Term{Type: rule.FLOAT},
		},
		{
			name: "Number",
			positiveExpressions: []rule.TypedExpression{
				integer, intParam, float, floatParam},
			negativeExpressions: []rule.TypedExpression{
				boolean, str, boolParam, stringParam, or, and, eq, not,
			},
			term: rule.Term{Type: rule.NUMBER},
		},
		{
			name: "Any",
			positiveExpressions: []rule.TypedExpression{
				integer, intParam, float, floatParam,
				boolean, str, boolParam, stringParam,
				or, and, eq, not,
			},
			negativeExpressions: nil,
			term:                rule.Term{Type: rule.ANY},
		},
	}

	// Run "IsFullfilledBy" for each test case with each positive and negative expression.
	for i, tc := range testCases {
		for j, pc := range tc.positiveExpressions {
			t.Run(fmt.Sprintf("%s[%d] positive case %d", tc.name, i, j),
				func(t *testing.T) {
					require.True(t, tc.term.IsFulfilledBy(pc))
				})
		}
		for k, nc := range tc.negativeExpressions {
			t.Run(fmt.Sprintf("%s[%d] negative case %d", tc.name, i, k),
				func(t *testing.T) {
					require.False(t, tc.term.IsFulfilledBy(nc))
				})
		}
	}
}

func TestTermEqual(t *testing.T) {
	t1 := rule.Term{Type: rule.STRING, Cardinality: rule.ONE}
	t2 := rule.Term{Type: rule.BOOLEAN, Cardinality: rule.ONE}
	t3 := rule.Term{Type: rule.STRING, Cardinality: rule.MANY}
	t4 := rule.Term{Type: rule.BOOLEAN, Cardinality: rule.MANY}

	require.True(t, t1.Equal(t1))
	require.False(t, t1.Equal(t2))
	require.False(t, t1.Equal(t3))
	require.False(t, t1.Equal(t4))
}

// GetOperatorExpr returns a TypedExpression by name
func TestGetOperatorExpr(t *testing.T) {
	expected, ok := rule.Eq(rule.BoolValue(true), rule.BoolValue(true)).(rule.TypedExpression)
	require.True(t, ok)
	op, err := rule.GetOperatorExpr("eq")
	require.NoError(t, err)
	te := op.(rule.TypedExpression)
	ec := expected.Contract()
	ac := te.Contract()
	require.True(t, ec.Equal(ac))
}

// Providing a non-existent expression name to GetOperatorExpr results in an error.
func TestGetOperatorExprBadName(t *testing.T) {
	_, err := rule.GetOperatorExpr("dave")
	require.EqualError(t, err, `No operator called "dave" exists`)
}
