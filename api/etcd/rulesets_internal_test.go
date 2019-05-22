package etcd

import (
	"encoding/json"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/heetch/regula"
	pb "github.com/heetch/regula/api/etcd/proto"
	"github.com/heetch/regula/rule"
	"github.com/stretchr/testify/require"
)

// TestPathMethods ensures that the correct path are returned by each method.
func TestPathMethods(t *testing.T) {
	s := &RulesetService{
		Namespace: "test",
	}

	exp := "test/rulesets/rules/path" + versionSeparator + "version"
	require.Equal(t, exp, s.rulesPath("path", "version"))

	exp = "test/rulesets/rules/path"
	require.Equal(t, exp, s.rulesPath("path", ""))

	exp = "test/rulesets/checksums/path"
	require.Equal(t, exp, s.checksumsPath("path"))

	exp = "test/rulesets/signatures/path"
	require.Equal(t, exp, s.signaturesPath("path"))
}

func BenchmarkProtoMarshalling(b *testing.B) {
	rules := []*rule.Rule{
		rule.New(rule.And(rule.Not(rule.BoolValue(false)), rule.BoolParam("param")), rule.BoolValue(true)),
		rule.New(rule.And(rule.BoolParam("1st-param"), rule.BoolParam("2nd-param")), rule.BoolValue(false)),
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := proto.Marshal(rulesToProtobuf(rules))
		require.NoError(b, err)
	}
}

func BenchmarkJSONMarshalling(b *testing.B) {
	rs := regula.NewRuleset(
		rule.New(rule.And(rule.Not(rule.BoolValue(false)), rule.BoolParam("param")), rule.BoolValue(true)),
		rule.New(rule.And(rule.BoolParam("1st-param"), rule.BoolParam("2nd-param")), rule.BoolValue(false)),
	)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, err := json.Marshal(rs)
		require.NoError(b, err)
	}
}

func BenchmarkProtoUnmarshalling(b *testing.B) {
	rules := []*rule.Rule{
		rule.New(rule.And(rule.Not(rule.BoolValue(false)), rule.BoolParam("param")), rule.BoolValue(true)),
		rule.New(rule.And(rule.BoolParam("1st-param"), rule.BoolParam("2nd-param")), rule.BoolValue(false)),
	}

	bb, err := proto.Marshal(rulesToProtobuf(rules))
	require.NoError(b, err)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var pbrs pb.Rules
		err := proto.Unmarshal(bb, &pbrs)
		require.NoError(b, err)
	}
}

func BenchmarkJSONUnmarshalling(b *testing.B) {
	rs := regula.NewRuleset(
		rule.New(rule.And(rule.Not(rule.BoolValue(false)), rule.BoolParam("param")), rule.BoolValue(true)),
		rule.New(rule.And(rule.BoolParam("1st-param"), rule.BoolParam("2nd-param")), rule.BoolValue(false)),
	)

	bb, err := json.Marshal(rs)
	require.NoError(b, err)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var rrs regula.Ruleset
		err := json.Unmarshal(bb, &rrs)
		require.NoError(b, err)
	}
}